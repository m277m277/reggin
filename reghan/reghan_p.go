package reghan

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yyle88/erero"
)

// Handle0pFunc 适用于没有参数且无ctx的处理函数的场景，认为不带ctx的属于非正式的场景，没法拿到上下文的信息，比如监控或者超时等信息，但比较简单
type Handle0pFunc[RES any] func() (RES, error)

// Handle1pFunc 适用于一个参数且无ctx的处理函数的场景
type Handle1pFunc[ARG, RES any] func(arg *ARG) (RES, error)

func Handle0p[RES any, RESPONSE any](run Handle0pFunc[RES], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, erx := run() //区别只在这里，这个不传ctx信息，因此处理逻辑里拿不到ctx信息，适用于简单场景
		ctx.SecureJSON(http.StatusOK, respFunc(ctx, res, erx))
	}
}

func Handle1p[ARG, RES any, RESPONSE any](run Handle1pFunc[ARG, RES], parseReq ParseReqFunc[ARG], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		arg, erx := parseReq(ctx)
		if erx != nil {
			var res RES // zero
			ctx.SecureJSON(http.StatusOK, respFunc(ctx, res, erero.WithMessage(erx, "PARAM IS WRONG")))
			return
		}
		res, erx := run(arg) //区别只在这里，这个不传ctx信息，因此处理逻辑里拿不到ctx信息，适用于简单场景
		ctx.SecureJSON(http.StatusOK, respFunc(ctx, res, erx))
	}
}

func P0[RES any, RESPONSE any](run Handle0pFunc[RES], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return Handle0p(run, respFunc)
}

func P1[ARG, RES any, RESPONSE any](run Handle1pFunc[ARG, RES], parseReq ParseReqFunc[ARG], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return Handle1p(run, parseReq, respFunc)
}

func PX[ARG, RES any, RESPONSE any](run Handle1pFunc[ARG, RES], respFunc MakeRespFunc[RES, RESPONSE]) gin.HandlerFunc {
	return Handle1p[ARG, RES, RESPONSE](run, BIND[ARG], respFunc)
}
