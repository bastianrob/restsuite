package main

import (
	"net/http"

	"github.com/bastianrob/go-httputil/middleware"
	oauth "github.com/bastianrob/go-oauth/handler"
	"github.com/bastianrob/restsuite/pkg/controller"
	"github.com/julienschmidt/httprouter"
)

func pipe(function middleware.HTTPMiddleware) http.HandlerFunc {
	return middleware.NewPipeline().
		Do(oauth.Authenticate()).
		Do(function).
		For(func(w http.ResponseWriter, r *http.Request) {})
}

func healthcheck(router *httprouter.Router) {
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})
}

func routesV1(router *httprouter.Router, controller controller.ScenarioController) {
	router.HandlerFunc("GET", "/v1/scenarios", pipe(controller.Find()))
	router.HandlerFunc("POST", "/v1/scenarios", pipe(controller.Add()))
	router.HandlerFunc("GET", "/v1/scenarios/:id", pipe(controller.Get()))
	router.HandlerFunc("PATCH", "/v1/scenarios/:id", pipe(controller.Update()))
}
