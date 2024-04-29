package api

import (
	"net/http"
	"time"

	"github.com/0x726f6f6b6965/tiny-url-go/internal/service"
	"github.com/0x726f6f6b6965/tiny-url-go/protos"
	"github.com/0x726f6f6b6965/tiny-url-go/utils"
	"github.com/gin-gonic/gin"
)

type ShortenAPI struct {
	ser service.ShortedURLService
}

func NewShortenAPI(ser service.ShortedURLService) *ShortenAPI {
	return &ShortenAPI{ser: ser}
}

func (s *ShortenAPI) Shorten(ctx *gin.Context) {
	data := new(protos.ShortenedURL)
	if err := ctx.ShouldBindJSON(data); err != nil {
		utils.InvalidParamErr.Message = "Please enter correct data."
		utils.Response(ctx, utils.SuccessCode, utils.InvalidParamErr, nil)
		return
	}
	var expire time.Time
	if data.ExpiresAt != 0 {
		expire = time.Unix(data.ExpiresAt, 0)
	}

	shortUrl, err := s.ser.ShortURL(ctx, data.Owner, data.Original, expire)
	if err != nil {
		utils.InternalServerError.Message = err.Error()
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, shortUrl)
}

func (s *ShortenAPI) RedirectURL(ctx *gin.Context) {
	short := ctx.Param("shorten")
	redirect, err := s.ser.RedirectURL(ctx, short)
	if err != nil {
		utils.InternalServerError.Message = err.Error()
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}
	ctx.Redirect(http.StatusPermanentRedirect, redirect)
}

func (s *ShortenAPI) DeleteURL(ctx *gin.Context) {
	data := new(protos.ShortenedURL)
	if err := ctx.ShouldBindJSON(data); err != nil {
		utils.InvalidParamErr.Message = "Please enter correct data."
		utils.Response(ctx, utils.SuccessCode, utils.InvalidParamErr, nil)
		return
	}
	err := s.ser.DeleteURL(ctx, data.Owner, data.Original)
	if err != nil {
		utils.InternalServerError.Message = err.Error()
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, nil)
}

func (s *ShortenAPI) UpdateURL(ctx *gin.Context) {
	data := new(protos.ShortenedURL)
	if err := ctx.ShouldBindJSON(data); err != nil {
		utils.InvalidParamErr.Message = "Please enter correct data."
		utils.Response(ctx, utils.SuccessCode, utils.InvalidParamErr, nil)
		return
	}
	var expire time.Time
	if data.ExpiresAt != 0 {
		expire = time.Unix(data.ExpiresAt, 0)
	}
	err := s.ser.UpdateURL(ctx, data.Owner, data.Shorten, data.Original, expire)
	if err != nil {
		utils.InternalServerError.Message = err.Error()
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, nil)
}
