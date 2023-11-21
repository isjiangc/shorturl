package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	v1 "nunu_ginblog/api/v1"
	"nunu_ginblog/internal/model"
	"nunu_ginblog/internal/service"
	"nunu_ginblog/pkg/helper/md5"
	"nunu_ginblog/pkg/helper/useragent"
	"strconv"
	"strings"
)

type ShorturlHandler struct {
	*Handler
	shorturlService service.ShorturlService
}

func NewShorturlHandler(handler *Handler, shorturlService service.ShorturlService) *ShorturlHandler {
	return &ShorturlHandler{
		Handler:         handler,
		shorturlService: shorturlService,
	}
}

// GenShortUrl godoc
// @Summary 生成短链
// @Schemes
// @Description
// @Tags 短链模块
// @Accept json
// @Produce json
// @Param request body v1.GenerateShortUrlRequest true "params"
// @Success 200 {object} v1.GenerateShortUrlResponse
// @Router /shorturl [post]
func (h *ShorturlHandler) GenShortUrl(ctx *gin.Context) {
	req := new(v1.GenerateShortUrlRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	shortUrlData, err := h.shorturlService.GenerateShortUrl(ctx, req)
	if err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.GenerateShortUrlData{
		ShortUrl: shortUrlData.ShortUrl,
	})
}

func (h *ShorturlHandler) ShortUrlDetail(ctx *gin.Context) {
	url := ctx.Param("url")
	if md5.EmptyString(url) {
		ctx.HTML(http.StatusNotFound, "404.html", gin.H{
			"message": "你访问的页面已失效了",
			"code":    http.StatusNotFound,
		})
		return
	}
	memUrl, err := h.shorturlService.Search4ShortUrl(ctx, url)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	if md5.EmptyString(memUrl.DestUrl) {
		ctx.HTML(http.StatusNotFound, "404.html", gin.H{
			"message": "你访问的页面已失效",
			"code":    http.StatusNotFound,
		})
		return
	}
	ua := ctx.Request.UserAgent()
	switch ot := memUrl.OpenType; ot {
	case model.OpenInAndroid:
		if useragent.IsAndroid(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInDingTalk:
		if useragent.IsDingTalk(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInChrome:
		if useragent.IsChrome(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInIPad:
		if useragent.IsIpad(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInIPhone:
		if useragent.IsIPhone(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInSafari:
		if useragent.IsSafari(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInWechat:
		if useragent.IsWaChatUA(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInFirefox:
		if useragent.IsFirefox(ua) {
			redirectSuccess(url, memUrl.DestUrl, ctx)
		} else {
			redirectFail(ctx)
		}
	case model.OpenInAll:
		redirectSuccess(url, memUrl.DestUrl, ctx)
	default:
		redirectFail(ctx)
	}
}

// UpdateUrlState godoc
// @Summary 更新短链状态
// @Schemes
// @Description
// @Tags 短链模块
// @Accept json
// @Produce json
// @Param url path string true "短链"
// @Param request body v1.UpdateShortUrlStateRequest true "params"
// @Success 200 {object} v1.UpdateShortUrlStateResponse
// @Router /shorturl/{url} [put]
func (h *ShorturlHandler) UpdateUrlState(ctx *gin.Context) {
	url := ctx.Param("url")
	enable := v1.UpdateShortUrlStateRequest{}
	ctx.BindJSON(&enable)
	if md5.EmptyString(strings.TrimSpace(url)) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrShortUrlEmpty, nil)
		return
	}
	res, err := h.shorturlService.ChangeState(ctx, url, &enable)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, &res)
}

// DeleteShortUrl godoc
// @Summary 删除短链
// @Schemes
// @Description
// @Tags 短链模块
// @Accept json
// @Produce json
// @Param url path string true "短链"
// @Success 200 {object} v1.DeleteShortUrlStateResponse
// @Router /shorturl/{url} [delete]
func (h *ShorturlHandler) DeleteShortUrl(ctx *gin.Context) {
	url := ctx.Param("url")
	if md5.EmptyString(strings.TrimSpace(url)) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrShortUrlEmpty, nil)
		return
	}
	res, err := h.shorturlService.DeleteShortUrl(ctx, url)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, &res)
}

// GetShortUrlList godoc
// @Summary 分页获取所有短链
// @Schemes
// @Description
// @Tags 短链模块
// @Accept json
// @Produce json
// @Param page query string false "页数"
// @Param size query string false "每页数"
// @Success 200 {object} v1.GetShortUrlListResponse
// @Router /shorturl [get]
func (h *ShorturlHandler) GetShortUrlList(ctx *gin.Context) {
	pageNum := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("size", "10")
	page, _ := strconv.Atoi(pageNum)
	size, _ := strconv.Atoi(pageSize)
	list, err := h.shorturlService.GetShortUrlList(ctx, &v1.GetShortUrlListRequest{PageNum: page, PageSize: size})
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, &list)
}

// GetShortUrlInfo godoc
// @Summary 获取短链信息
// @Schemes
// @Description
// @Tags 短链模块
// @Accept json
// @Produce json
// @Param url path string true "短链"
// @Success 200 {object} model.ShortUrl
// @Router /shorturl/{url} [get]
func (h *ShorturlHandler) GetShortUrlInfo(ctx *gin.Context) {
	url := ctx.Param("url")
	if md5.EmptyString(strings.TrimSpace(url)) {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrShortUrlEmpty, nil)
		return
	}
	stat, err := h.shorturlService.GetShorturlByUrl(ctx, url)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, &stat)
}

func redirectSuccess(shortUrl, descUrl string, ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, descUrl)
}
func redirectFail(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "404.html", gin.H{
		"code":    http.StatusNotFound,
		"message": "不支持的打开方式",
	})
}
