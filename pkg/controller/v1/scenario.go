package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
				return
			}

			resp := controller.Response{Data: scenarios}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
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
				return
			}

			resp := controller.Response{Data: scenario}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *scenarioHandler) Run() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ps := httprouter.ParamsFromContext(ctx)
			id := ps.ByName("id")

			output, err := hndl.svc.Run(ctx, id)
			if err != nil {
				log.Println("ERR Failed to run test scenario with ID:", id, err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
				return
			}

			result := strings.Split(output, "\n")
			resp := controller.Response{Data: result}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
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
				return
			}

			scenario := scenario.New()
			json.Unmarshal(reqBody, scenario)

			scenario, err = hndl.svc.Add(ctx, scenario)
			if err != nil {
				log.Println("ERR Failed to add new test scenario.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
				return
			}

			resp := controller.Response{Data: scenario}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusCreated)
			w.Write(body)
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
				return
			}

			scenario := scenario.New()
			json.Unmarshal(reqBody, scenario)

			scenario, err = hndl.svc.Update(ctx, id, scenario)
			if err != nil {
				log.Println("ERR Failed to update test scenario.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
				return
			}

			resp := controller.Response{Data: scenario}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusCreated)
			w.Write(body)
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
				return
			}

			w.WriteHeader(http.StatusNoContent)
			h.ServeHTTP(w, r)
		}
	}
}
