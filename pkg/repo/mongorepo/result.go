package mongorepo

import (
	"context"

	restify "github.com/bastianrob/go-restify"
	"github.com/bastianrob/restsuite/pkg/exception"
	"github.com/bastianrob/restsuite/pkg/repo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type resultRepo struct {
	client *mongo.Client
}

//NewResultRepo using mongodb
func NewResultRepo(client *mongo.Client) repo.ResultRepo {
	return &resultRepo{client: client}
}

//deconstructParam defaults passed parameters into 4 category:
//1. where filter,
//2. fields projection
//3. sorts
//4. and paging (skip & limit) with default skip = 0, and limit to 25
func deconstructParam(params ...map[string]interface{}) []map[string]interface{} {
	result := []map[string]interface{}{nil, nil, nil, {"skip": 0, "limit": 25}}
	for i := 0; i <= 3; i++ {
		if len(params) <= i {
			break
		}

		result[i] = params[i]
	}

	return result
}

func (repo *resultRepo) Find(ctx context.Context, dbname string, params ...map[string]interface{}) ([]restify.TestResult, int64, error) {
	params = deconstructParam(params...)
	where, fields, sorts, paging := params[0], params[1], params[2], params[3]

	collection := repo.client.Database(dbname).Collection("results")
	opts := options.Find().
		SetProjection(fields).
		SetSort(sorts).
		SetSkip(paging["skip"].(int64)).
		SetLimit(paging["limit"].(int64))

	total, err := collection.CountDocuments(ctx, where)
	if err != nil {
		return nil, 0, exception.New(500, "Failed to find total test result count from database. %s", err.Error())
	}

	res, err := collection.Find(ctx, where, opts)
	if err != nil {
		return nil, 0, exception.New(500, "Failed to find test result from database. %s", err.Error())
	}

	results := []restify.TestResult{}
	defer res.Close(ctx)
	for res.Next(ctx) {
		result := restify.TestResult{}
		err = res.Decode(&result)
		if err != nil {
			return nil, 0, exception.New(500, "Failed to decode test result from database")
		}

		results = append(results, result)
	}

	return results, total, nil
}
