package reghan

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yyle88/erero"
)

type Handle0pFunc[RES any] func() (*RES, error)
type Handle1pFunc[ARG, RES any] func(arg *ARG) (*RES, error)

type MakeRespFunc[RES any, RESPONSE any] func(ctx *gin.Context, res *RES, erx error) *RESPONSE //使用指针类型拼返回值

func Handle0p[RES any, RESPONSE any](run Handle0pFunc[RES], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, erx := run()
		ctx.SecureJSON(http.StatusOK, respFunc(ctx, res, erx))
	}
}

func Handle1p[ARG, RES any, RESPONSE any](run Handle1pFunc[ARG, RES], parseReq ParseReqFunc[ARG], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		arg, erx := parseReq(ctx)
		if erx != nil {
			//出错时就没有返回值啦
			ctx.SecureJSON(http.StatusOK, respFunc(ctx, nil, erero.WithMessage(erx, "PARAM IS WRONG")))
			return
		}
		res, erx := run(arg)
		ctx.SecureJSON(http.StatusOK, respFunc(ctx, res, erx))
	}
}

func Handle1x[ARG, RES any, RESPONSE any](run Handle1pFunc[ARG, RES], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return Handle1p[ARG, RES, RESPONSE](run, BindJson[ARG], respFunc)
}
