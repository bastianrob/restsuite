package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bastianrob/go-httputil/middleware"
	"github.com/bastianrob/go-httputil/queryables"
	"github.com/bastianrob/go-httputil/queryables/qtype"
	"github.com/bastianrob/restsuite/pkg/controller"
	"github.com/bastianrob/restsuite/pkg/exception"
	"github.com/bastianrob/restsuite/pkg/service"
)

var (
	resultQueryables = queryables.Collection{
		{QueryKey: "time", DBKey: "timestamp", Type: qtype.Integer, Transform: timestampTransform},
		{QueryKey: "scenario_name", DBKey: "scenario_name", Type: qtype.String, Transform: stringRegexTransform},
		{QueryKey: "case_name", DBKey: "test_case_name", Type: qtype.String, Transform: stringRegexTransform},
		{QueryKey: "response_code", DBKey: "response_code", Type: qtype.Integer, Transform: statusCodeTransform},
		{QueryKey: "expected_code", DBKey: "expected_code", Type: qtype.Integer, Transform: statusCodeTransform},
	}

	resultPaging = queryables.Collection{
		{QueryKey: "skip", DBKey: "skip", Type: qtype.Integer, Transform: pagingTransform},
		{QueryKey: "limit", DBKey: "limit", Type: qtype.Integer, Transform: pagingTransform},
	}
)

type resultHandler struct {
	svc service.ResultService
}

//NewResultHandler v1
func NewResultHandler(svc service.ResultService) controller.ResultController {
	return &resultHandler{svc}
}

func pagingTransform(k string, v []interface{}) (rk string, rv interface{}, err error) {
	return k, v[0], nil
}

//timestamp received from javascript so expecting millisecond
func timestampTransform(k string, v []interface{}) (rk string, rv interface{}, err error) {
	length := len(v)
	if length == 0 {
		return "", nil, errors.New("Query value is empty")
	}

	//exact timestamp, which is very unlikely
	if length == 1 {
		return k, v[0].(int64) * 1000000, nil
	}

	//timestamp is between
	if length >= 2 {
		start, end := v[0].(int64), v[1].(int64)
		return k, map[string]interface{}{
			"$gte": start * 1000000,
			"$lte": end * 1000000,
		}, nil
	}

	return "", nil, errors.New(("Unhandled query value length"))
}

func stringRegexTransform(k string, v []interface{}) (rk string, rv interface{}, err error) {
	length := len(v)
	if length == 0 {
		return "", nil, errors.New("Query value is empty")
	}

	if length == 1 {
		return k, map[string]interface{}{
			"$regex":   v[0],
			"$options": "i",
		}, nil
	}

	//else $or: [{key: {$regex: str[0], $options: i}}, {key: {$regex: str[2], $options: i}}, ...]
	//VERY SLOW!
	col := make([]map[string]interface{}, length)
	for i, str := range v {
		col[i] = map[string]interface{}{
			k: map[string]interface{}{
				"$regex":   str,
				"$options": "i",
			},
		}
	}
	return "$or", col, nil
}

func statusCodeTransform(k string, v []interface{}) (rk string, rv interface{}, err error) {
	length := len(v)
	if length == 0 {
		return "", nil, errors.New("Query value is empty")
	}

	if length == 1 {
		return k, v[0], nil
	}

	//if length > 1,
	return k, map[string]interface{}{"$in": v}, nil
}

func (hndl *resultHandler) Find() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			resp := controller.Response{}

			var fields map[string]interface{} //empty field projection
			query := resultQueryables.ToQuery(r)
			paging := resultPaging.ToQuery(r)
			sorts := map[string]interface{}{"timestamp": -1} //sort by timestamp, descending
			log.Println("WHERE:", query, "\nFIELDS:", fields, "\nSORTS:", sorts, "\nPAGING:", paging)

			results, _, err := hndl.svc.Find(ctx, query, fields, sorts, paging)
			if err != nil {
				status, message := http.StatusInternalServerError, fmt.Sprintf("ERR Failed to find test result. %s", err.Error())
				exc, ok := exception.IsException(err)
				if ok {
					status, message = exc.Code(), exc.Message()
				}

				resp.Message = message
				body, _ := json.Marshal(resp)
				w.WriteHeader(status)
				w.Write(body)
				h.ServeHTTP(w, r)
				return
			}

			resp.Data = results
			body, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			h.ServeHTTP(w, r)
		}
	}
}
