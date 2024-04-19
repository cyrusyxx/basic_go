package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/service"
	ijwt "webook/webook/internal/web/jwt"
	"webook/webook/pkg/logger"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.Logger
}

type Page struct {
	Limit  int64
	Offset int64
}

func NewArticleHandler(l logger.Logger,
	svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/article")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)

	g.GET("/detail/:id", h.Detail)
	g.POST("/list", h.List)

	g.GET("/pub/:id", h.PubDetail)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to save article", logger.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	uc := ctx.MustGet("userclaim").(ijwt.UserClaims)
	id, err := h.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to save article", logger.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
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
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to withdraw article", logger.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
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
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article list", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: toAbstractVos(artis),
	})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	// Get id from path
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Invalid parameter: id",
		})
		return
	}

	// Get article by id
	arti, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
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
		ctx.JSON(http.StatusOK, Result{
			Msg:  "System Error",
			Code: 5,
		})
		h.l.Error("invalid article id",
			logger.Int64(id),
			logger.Int64(uc.Uid))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: toContentVo(arti),
	})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	// Get id from path
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Invalid parameter: id",
		})
		return
	}

	// Get article by id
	arti, err := h.svc.GetPubById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "System Error",
			Data: nil,
		})
		h.l.Error("Failed to get article", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: toContentVo(arti),
		Msg:  "OK",
	})
}

func toAbstractVo(articles domain.Article) ArticleVo {
	return _tovo(articles, true)
}

func toContentVo(articles domain.Article) ArticleVo {
	return _tovo(articles, false)
}

func toAbstractVos(articles []domain.Article) []ArticleVo {
	var vos []ArticleVo
	for _, article := range articles {
		vos = append(vos, _tovo(article, true))
	}
	return vos
}

func toContentVos(articles []domain.Article) []ArticleVo {
	var vos []ArticleVo
	for _, article := range articles {
		vos = append(vos, _tovo(article, false))
	}
	return vos
}

func _tovo(article domain.Article, isAbstract bool) ArticleVo {
	vo := ArticleVo{
		Id:         article.Id,
		Title:      article.Title,
		Abstract:   "",
		Content:    "",
		AuthorId:   article.Author.Id,
		AuthorName: article.Author.Name,
		Status:     uint8(article.Status),

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
