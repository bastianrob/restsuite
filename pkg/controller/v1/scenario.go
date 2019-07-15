package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/bastianrob/go-httputil/middleware"
	"github.com/bastianrob/go-restify/scenario"
	"github.com/bastianrob/restsuite/pkg/controller"
	"github.com/bastianrob/restsuite/pkg/service"
)

type scenarioHandler struct {
	svc service.ScenarioService
}

//NewScenarioHandler v1
func NewScenarioHandler(svc service.ScenarioService) controller.ScenarioController {
	return &scenarioHandler{svc}
}

func (hndl *scenarioHandler) Find() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			name := r.FormValue("name")
			scenarios, err := hndl.svc.Find(ctx, name)
			if err != nil {
				log.Println("ERR Failed to find test scenario.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
			}

			resp := controller.Response{Data: scenarios}
			body, _ := json.Marshal(resp)

			w.Write(body)
			w.WriteHeader(http.StatusOK)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *scenarioHandler) Get() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ps := httprouter.ParamsFromContext(ctx)
			id := ps.ByName("id")

			scenario, err := hndl.svc.Get(ctx, id)
			if err != nil {
				log.Println("ERR Failed to get test scenario with ID:", id, ".", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
			}

			resp := controller.Response{Data: scenario}
			body, _ := json.Marshal(resp)

			w.Write(body)
			w.WriteHeader(http.StatusOK)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *scenarioHandler) Add() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("ERR Failed to parse request body.", err.Error())

				w.WriteHeader(http.StatusBadRequest)
				h.ServeHTTP(w, r)
			}

			scenario := scenario.New()
			json.Unmarshal(reqBody, scenario)

			scenario, err = hndl.svc.Add(ctx, scenario)
			if err != nil {
				log.Println("ERR Failed to add new test scenario.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
			}

			resp := controller.Response{Data: scenario}
			body, _ := json.Marshal(resp)

			w.Write(body)
			w.WriteHeader(http.StatusCreated)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *scenarioHandler) Update() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ps := httprouter.ParamsFromContext(ctx)
			id := ps.ByName("id")

			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("ERR Failed to parse request body.", err.Error())

				w.WriteHeader(http.StatusBadRequest)
				h.ServeHTTP(w, r)
			}

			scenario := scenario.New()
			json.Unmarshal(reqBody, scenario)

			scenario, err = hndl.svc.Update(ctx, id, scenario)
			if err != nil {
				log.Println("ERR Failed to update test scenario.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
			}

			resp := controller.Response{Data: scenario}
			body, _ := json.Marshal(resp)

			w.Write(body)
			w.WriteHeader(http.StatusCreated)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *scenarioHandler) Delete() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ps := httprouter.ParamsFromContext(ctx)
			id := ps.ByName("id")

			err := hndl.svc.Delete(ctx, id)
			if err != nil {
				log.Println("ERR Failed to delete test scenario with ID:", id, ".", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
			}

			w.WriteHeader(http.StatusNoContent)
			h.ServeHTTP(w, r)
		}
	}
}
