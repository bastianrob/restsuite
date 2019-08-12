package service

import (
	"context"

	restify "github.com/bastianrob/go-restify"
)

//ResultQueryable field map of what is queryable in result resource
type ResultQueryable struct {
	QueryKey  string
	DBKey     string
	Transform interface{}
}

//ResultService Contract
type ResultService interface {
	//Find list of results, can contains http query
	Find(ctx context.Context, params ...map[string]interface{}) ([]restify.TestResult, int64, error)
}
