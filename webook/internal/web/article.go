package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/errs"
	"webook/webook/internal/service"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/pkg/ginx"
	"webook/webook/pkg/logger"
)

type ArticleHandler struct {
	svc      service.ArticleService
	interSvc service.InteractiveService

	l   logger.Logger
	biz string
}

type Page struct {
	Limit  int64
	Offset int64
}

func NewArticleHandler(l logger.Logger,
	svc service.ArticleService,
	intersvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		interSvc: intersvc,
		l:        l,
		biz:      "article",
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
		Data: toAbstractVos(artis),
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

func (h *ArticleHandler) Like(ctx *gin.Context) {
	type Req struct {
		Id   int64 `json:"id"`
		Like bool  `json:"like"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	var err error
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
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	err := h.interSvc.Collect(ctx, h.biz, req.Id, req.Cid, uc.Uid)
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

func toAbstractVo(articles domain.Article) ArticleVo {
	return _tovo(articles, domain.InteractiveCount{}, true)
}

func toContentVo(articles domain.Article, inter domain.InteractiveCount) ArticleVo {
	return _tovo(articles, inter, false)
}

func toAbstractVos(articles []domain.Article) []ArticleVo {
	var vos []ArticleVo
	for _, article := range articles {
		vos = append(vos, _tovo(article, domain.InteractiveCount{}, true))
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
