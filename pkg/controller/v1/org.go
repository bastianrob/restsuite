package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bastianrob/go-httputil/middleware"
	"github.com/bastianrob/restsuite/model"
	"github.com/bastianrob/restsuite/pkg/controller"
	"github.com/bastianrob/restsuite/pkg/exception"
	"github.com/bastianrob/restsuite/pkg/service"
)

type orgHandler struct {
	svc service.OrganizationService
}

//NewOrganizationHandler v1
func NewOrganizationHandler(svc service.OrganizationService) controller.OrganizationController {
	return &orgHandler{svc}
}

func (hndl *orgHandler) Get() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			org, err := hndl.svc.Get(ctx)
			if err != nil {
				log.Println("ERR Failed to get test organization.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
				return
			}

			resp := controller.Response{Data: org}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusOK)
			w.Write(body)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *orgHandler) Add() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			orgname := r.FormValue("org")
			resp := controller.Response{}

			//Fail case 1
			if len(orgname) < 4 {
				resp.Message = "Organization name must be at least 4 characters"
				body, _ := json.Marshal(resp)

				w.WriteHeader(http.StatusBadRequest)
				w.Write(body)
				h.ServeHTTP(w, r)
				return
			}

			org, err := hndl.svc.Add(ctx, orgname)

			//Fail case 2
			if err != nil {
				status, message := http.StatusInternalServerError, err.Error()
				exc, ok := exception.IsException(err)
				if ok {
					status, message = exc.Code(), exc.Message()
				}

				log.Println("ERR Failed to add new organization.", err.Error())

				resp.Message = message
				body, _ := json.Marshal(resp)

				w.WriteHeader(status)
				w.Write(body)
				h.ServeHTTP(w, r)
				return
			}

			//Success
			resp.Data = org
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusCreated)
			w.Write(body)
			h.ServeHTTP(w, r)
		}
	}
}

func (hndl *orgHandler) Update() middleware.HTTPMiddleware {
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

			org := model.Organization{}
			json.Unmarshal(reqBody, &org)

			org, err = hndl.svc.Update(ctx, org)
			if err != nil {
				log.Println("ERR Failed to update organization settings.", err.Error())

				w.WriteHeader(http.StatusInternalServerError)
				h.ServeHTTP(w, r)
				return
			}

			resp := controller.Response{Data: org}
			body, _ := json.Marshal(resp)

			w.WriteHeader(http.StatusCreated)
			w.Write(body)
			h.ServeHTTP(w, r)
		}
	}
}
