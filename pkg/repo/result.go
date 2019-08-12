package repo

import (
	"context"

	restify "github.com/bastianrob/go-restify"
)

//ResultRepo contract
type ResultRepo interface {
	//Find lists of results
	//params consists of: where, fields, sorts, paging
	//returns: Array of test resuts, Total Datacount based on filter, and/or error
	Find(ctx context.Context, dbname string, params ...map[string]interface{}) ([]restify.TestResult, int64, error)
}
