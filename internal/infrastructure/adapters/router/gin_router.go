package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type GinRouter struct {
	engine *gin.Engine
	group  *gin.RouterGroup
}

func NewGinRouter() *GinRouter {
	g := gin.New()
	g.Use(gin.Recovery(), gin.Logger())
	return &GinRouter{g, nil}
}

func (r *GinRouter) activeGroup() *gin.RouterGroup {
	if r.group == nil {
		return r.group
	}
	return &r.engine.RouterGroup
}

func httpHandlerAdapter(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func middlewareAdapter(middleware func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		terminal := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				c.Next()
			},
		)
		wrapped := middleware(terminal)
		wrapped.ServeHTTP(c.Writer, c.Request)
	}
}

func (r *GinRouter) Handle(method, path string, handler http.Handler) {
	g := r.activeGroup()
	g.Handle(method, path, httpHandlerAdapter(handler))
}

func (r *GinRouter) GET(path string, handler http.Handler) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *GinRouter) POST(path string, handler http.Handler) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *GinRouter) PUT(path string, handler http.Handler) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *GinRouter) DELETE(path string, handler http.Handler) {
	r.Handle(http.MethodDelete, path, handler)
}

func (r *GinRouter) Group(prefix string) Router {
	group := r.activeGroup().Group(prefix)
	return &GinRouter{r.engine, group}
}

func (r *GinRouter) Use(middleware ...func(http.Handler) http.Handler) {
	for _, m := range middleware {
		r.activeGroup().Use(middlewareAdapter(m))
	}
}

func (r *GinRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.engine.ServeHTTP(writer, request)
}