package cas

import (
	"github.com/gin-gonic/gin"
	goCas "github.com/go-cas/cas"
	"net/http"
)

type casMiddleware struct {
	casClient *goCas.Client
	handler http.Handler
}

func (casMiddleware casMiddleware) authed(c *gin.Context) {
	c.Set("CASUsername", goCas.Username(c.Request))
	c.Set("CASAttributes", goCas.Attributes(c.Request))
	c.Next()
	return
}

func (casMiddleware casMiddleware) middlewareFunc(c *gin.Context) {
	casMiddleware.handler.ServeHTTP(c.Writer, c.Request)
	if goCas.IsAuthenticated(c.Request) {
		casMiddleware.authed(c)
		return
	}
	c.Abort()
}

type Options = goCas.Options
func MiddlewareFunc(options *Options) gin.HandlerFunc {
	casClient := goCas.NewClient((*goCas.Options)(options))
	rawHandler := func(res http.ResponseWriter, req *http.Request) {
		if goCas.IsAuthenticated(req) {
			return
		}
		if goCas.IsNewLogin(req) {
			return
		}
		casClient.RedirectToLogin(res, req)
	}
	return casMiddleware{
		casClient: casClient,
		handler: casClient.HandleFunc(rawHandler),
	}.middlewareFunc
}
