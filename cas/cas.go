package cas

import (
	goCas "github.com/go-cas/cas"
)

type casMiddleware struct {
	casClient *goCas.Client
}

func (casMiddleware casMiddleware) casMiddleware(c *gin.Context) {
	casMiddleware.casClient.Handler(c.Writer, c.Request)
	c.Set("CASUsername", casMiddleware.casClient.Username(c.Request))
	c.Set("CASAttributes", casMiddleware.casClient.Attributes(c.Request))
}

func MiddlewareFunc(casClient *goCas.Client) gin.HandlerFunc {
	return casMiddleware{casClient: casClient}.middlewareFunc
}
