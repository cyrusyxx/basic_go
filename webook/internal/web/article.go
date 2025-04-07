package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/errs"
	"webook/webook/internal/service"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/pkg/ginx"
	"webook/webook/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type ArticleHandler struct {
	svc        service.ArticleService
	interSvc   service.InteractiveService
	rankSvc    service.RankingService
	commentSvc service.CommentService

	l   logger.Logger
	biz string
}

type Page struct {
	Limit  int64
	Offset int64
}

type CreateCommentReq struct {
	Content string `json:"content"`
}

func NewArticleHandler(l logger.Logger,
	svc service.ArticleService,
	intersvc service.InteractiveService,
	ranksvc service.RankingService,
	commentSvc service.CommentService) *ArticleHandler {
	return &ArticleHandler{
		svc:        svc,
		interSvc:   intersvc,
		rankSvc:    ranksvc,
		commentSvc: commentSvc,
		l:          l,
		biz:        "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/article")

	g.POST("/edit", ginx.WrapBodyAndClaims(h.Edit))
	g.POST("/publish", ginx.WrapBodyAndClaims(h.Publish))
	g.POST("/withdraw", h.Withdraw)

	g.GET("/detail/:id", h.Detail)
	g.POST("/list", h.List)

	g.GET("/pub/:id", h.PubDetail)
	g.POST("/pub/like", h.Like)
	g.POST("/pub/collect", h.Collect)
	g.POST("/pub/list", h.PubList)

	g.GET("/pub/top", h.Top)

	// 评论相关路由
	g.POST("/:id/comment", h.CreateComment)
	g.GET("/:id/comments", h.ListComments)
	g.POST("/:id/comment/:commentId/delete", h.DeleteComment)
}

func (h *ArticleHandler) Edit(ctx *gin.Context, req EditReq, uc ijwt.UserClaims) (ginx.Result, error) {

	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return ginx.Result{
			Code: errs.ArticleInternalServerError,
			Msg:  "System Error",
			Data: nil,
		}, fmt.Errorf("failed to save article: %w", err)
	}

	return ginx.Result{
		Data: id,
	}, nil
}

func (h *ArticleHandler) Publish(ctx *gin.Context, req PublishReq, uc ijwt.UserClaims) (ginx.Result, error) {
	id, err := h.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return ginx.Result{
			Code: errs.ArticleInternalServerError,
			Msg:  "System Error",
			Data: nil,
		}, fmt.Errorf("failed to publish article: %w", err)
	}
	return ginx.Result{
		Data: id,
	}, nil
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req

	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)

	err := h.svc.Withdraw(ctx, uc.Uid, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 0,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to withdraw article", logger.Error(err))
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}
	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)

	artis, err := h.svc.GetByAuthor(ctx, uc.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article list", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Data: toAbstractVos(artis, make(map[int64]domain.InteractiveCount)),
	})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	// Get id from path
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "Invalid parameter: id",
		})
		return
	}

	// Get article by id
	arti, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article", logger.Error(err))
		return
	}

	// Check if the id is right
	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	if arti.Author.Id != uc.Uid {
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "System Error",
			Code: 5,
		})
		h.l.Error("invalid article id",
			logger.Int64("id", id),
			logger.Int64("uid", uc.Uid))
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Data: toContentVo(arti, domain.InteractiveCount{}),
	})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	// Get id from path
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "Invalid parameter: id",
		})
		return
	}

	var (
		eg    errgroup.Group
		arti  domain.Article
		inter domain.InteractiveCount
	)

	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)

	eg.Go(func() error {
		var er error
		arti, er = h.svc.GetPubById(ctx, uc.Uid, id)
		return er
	})

	eg.Go(func() error {
		var er error
		inter, er = h.interSvc.Get(ctx, h.biz, id, uc.Uid)
		return er
	})

	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article", logger.Error(err))
		return
	}

	// Increase view count
	err = h.interSvc.IncreaseViewCount(ctx, h.biz, id)

	// Return article
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: toContentVo(arti, inter),
		Msg:  "OK",
	})
}

func (h *ArticleHandler) PubList(ctx *gin.Context) {
	// Get offset and limit from query
	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}

	var (
		artis    []domain.Article
		interMap map[int64]domain.InteractiveCount
	)

	artis, err := h.svc.ListPub(ctx, time.Now(), page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article", logger.Error(err))
		return
	}

	interMap, err = h.interSvc.GetByIds(ctx, h.biz, domain.ArticleList(artis).Ids())

	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article", logger.Error(err))
		return
	}

	// Return article
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: toAbstractVos(artis, interMap),
		Msg:  "OK",
	})
}

func (h *ArticleHandler) Like(ctx *gin.Context) {
	type Req struct {
		Id   int64 `json:"id"`
		Like bool  `json:"like"`
	}
	var req Req
	var err error
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)

	if req.Like {
		err = h.interSvc.Like(ctx, h.biz, req.Id, uc.Uid)
	} else {
		err = h.interSvc.CancelLike(ctx, h.biz, req.Id, uc.Uid)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
		})
		h.l.Error("Failed to like article", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) Collect(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
		// Collections ID
		Cid int64 `json:"cid"`
		// true: 收藏, false: 取消收藏
		Collect bool `json:"collect"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)

	var err error
	if req.Collect {
		err = h.interSvc.Collect(ctx, h.biz, req.Id, req.Cid, uc.Uid)
	} else {
		err = h.interSvc.CancelCollect(ctx, h.biz, req.Id, uc.Uid)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
		})
		h.l.Error("Failed to collect article", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) Top(ctx *gin.Context) {
	// 获取排行榜文章
	articles, err := h.rankSvc.GetTopN(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get top articles", logger.Error(err))
		return
	}

	if len(articles) == 0 {
		ctx.JSON(http.StatusOK, ginx.Result{
			Data: []domain.Article{},
			Msg:  "OK",
		})
		return
	}

	// 获取文章ID列表
	var articleIds []int64
	for _, article := range articles {
		articleIds = append(articleIds, article.Id)
	}

	// 获取互动数据
	interMap, err := h.interSvc.GetByIds(ctx, h.biz, articleIds)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get interactive data", logger.Error(err))
		return
	}

	// 返回文章列表
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: toAbstractVos(articles, interMap),
		Msg:  "OK",
	})
}

func toAbstractVo(articles domain.Article) ArticleVo {
	return _tovo(articles, domain.InteractiveCount{}, true)
}

func toContentVo(articles domain.Article, inter domain.InteractiveCount) ArticleVo {
	return _tovo(articles, inter, false)
}

func toAbstractVos(articles []domain.Article, interMap map[int64]domain.InteractiveCount) []ArticleVo {
	var vos []ArticleVo

	for _, article := range articles {
		inter, ok := interMap[article.Id]
		if !ok {
			inter = domain.InteractiveCount{}
		}
		vos = append(vos, _tovo(article, inter, true))
	}
	return vos
}

func toContentVos(articles []domain.Article) []ArticleVo {
	var vos []ArticleVo
	for _, article := range articles {
		vos = append(vos, _tovo(article, domain.InteractiveCount{}, false))
	}
	return vos
}

func _tovo(article domain.Article, inter domain.InteractiveCount, isAbstract bool) ArticleVo {
	vo := ArticleVo{
		Id:         article.Id,
		Title:      article.Title,
		Abstract:   "",
		Content:    "",
		AuthorId:   article.Author.Id,
		AuthorName: article.Author.Name,
		Status:     uint8(article.Status),

		// interactive field
		ViewCnt:    inter.ViewCnt,
		LikeCnt:    inter.LikeCnt,
		CollectCnt: inter.CollectCnt,
		Liked:      inter.Liked,
		Collected:  inter.Collected,

		// Format time to string
		Ctime: article.Ctime.Format(time.DateTime),
		Utime: article.Utime.Format(time.DateTime),
	}
	if isAbstract {
		vo.Abstract = article.Abstract()
	} else {
		vo.Content = article.Content
	}
	return vo
}

func (h *ArticleHandler) CreateComment(ctx *gin.Context) {
	articleIdStr := ctx.Param("id")
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		})
		return
	}

	var req CreateCommentReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	id, err := h.commentSvc.Create(ctx, domain.Comment{
		Content:   req.Content,
		ArticleId: articleId,
		User: domain.User{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("创建评论失败", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Data: id,
	})
}

func (h *ArticleHandler) ListComments(ctx *gin.Context) {
	articleIdStr := ctx.Param("id")
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		})
		return
	}

	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}

	comments, err := h.commentSvc.GetByArticleId(ctx, articleId, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获取评论列表失败", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Data: comments,
	})
}

func (h *ArticleHandler) DeleteComment(ctx *gin.Context) {
	commentIdStr := ctx.Param("commentId")
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		})
		return
	}

	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	err = h.commentSvc.DeleteById(ctx, commentId, uc.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("删除评论失败", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}
