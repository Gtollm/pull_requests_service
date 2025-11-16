package router

import "net/http"

type Router interface {
	Handle(method, path string, handler http.Handler)
	GET(path string, handler http.Handler)
	POST(path string, handler http.Handler)
	PUT(path string, handler http.Handler)
	DELETE(path string, handler http.Handler)

	Group(prefix string) Router
	Use(middleware ...func(http.Handler) http.Handler)
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}