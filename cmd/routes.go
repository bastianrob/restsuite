package main

import (
	"net/http"

	"github.com/bastianrob/go-httputil/middleware"
	oauth "github.com/bastianrob/go-oauth/handler"
	"github.com/bastianrob/restsuite/pkg/controller"
	"github.com/julienschmidt/httprouter"
)

func final(w http.ResponseWriter, r *http.Request) {}

func pipe(endpoint func() middleware.HTTPMiddleware) http.HandlerFunc {
	return middleware.
		NewPipeline().              //HTTP pipeline
		Do(oauth.Authenticate()).   //Authenticate user
		Do(controller.UnwrapJWT()). //Unwrap JWT token to get oranization name
		Do(endpoint()).             //Call the endpoint
		For(final)                  //Final call, do nothing? track metric?
}

func healthcheck(router *httprouter.Router) {
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})
}

func scenarioV1(router *httprouter.Router, controller controller.ScenarioController) {
	router.HandlerFunc("GET", "/v1/scenarios", pipe(controller.Find))
	router.HandlerFunc("POST", "/v1/scenarios", pipe(controller.Add))
	router.HandlerFunc("GET", "/v1/scenarios/:id", pipe(controller.Get))
	router.HandlerFunc("PATCH", "/v1/scenarios/:id", pipe(controller.Update))
	router.HandlerFunc("RUN", "/v1/scenarios/:id", pipe(controller.Run))
}

func orgV1(router *httprouter.Router, controller controller.OrganizationController) {
	router.HandlerFunc("GET", "/v1/org", pipe(controller.Get))
	router.HandlerFunc("POST", "/v1/org", pipe(controller.Add))
	router.HandlerFunc("PATCH", "/v1/org", pipe(controller.Update))
}
