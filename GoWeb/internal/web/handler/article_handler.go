package handler

import (
	ginx "github.com/LEILEI0628/GinPro/GinX"
	jwtx "github.com/LEILEI0628/GinPro/middleware/jwt"
	loggerx "github.com/LEILEI0628/GinPro/middleware/logger"
	service2 "github.com/LEILEI0628/GoWeb/interactive/service"
	"github.com/LEILEI0628/GoWeb/internal/domain"
	"github.com/LEILEI0628/GoWeb/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArticleHandler struct {
	svc    service.ArticleServiceInterface
	l      loggerx.Logger
	itrSvc service2.InteractiveServiceInterface
}

func NewArticleHandler(svc service.ArticleServiceInterface, itrSvc service2.InteractiveServiceInterface, logger loggerx.Logger) *ArticleHandler {
	return &ArticleHandler{svc: svc, itrSvc: itrSvc, l: logger}
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//TODO 检测输入

	c := ctx.MustGet("claims")
	claims, ok := c.(*jwtx.UserClaims)
	if !ok {
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, ginx.Result{Code: 5, Msg: "系统错误"})
		h.l.Error("未发现用户的session信息")
		return
	}
	UID := claims.UID

	// 调用svc
	id, err := h.svc.Save(ctx.Request.Context(), domain.Article{Id: req.Id, Title: req.Title, Content: req.Content,
		Author: domain.Author{Id: UID}})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: 5, Msg: "系统错误"})
		h.l.Error("article保存失败", loggerx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "OK", Data: id})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*jwtx.UserClaims)
	if !ok {
		// 可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}

	id, err := h.svc.Publish(ctx, req.toDomain(claims.UID))
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		h.l.Error("发表帖子失败", loggerx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "OK",
		Data: id,
	})
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
