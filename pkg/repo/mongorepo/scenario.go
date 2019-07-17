package mongorepo

import (
	"context"
	"time"

	"github.com/bastianrob/go-restify/scenario"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	restify "github.com/bastianrob/go-restify"
	"github.com/bastianrob/restsuite/pkg/exception"
	"github.com/bastianrob/restsuite/pkg/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

type scenarioRepo struct {
	client *mongo.Client
}

//NewScenarioRepo using mongodb
func NewScenarioRepo(client *mongo.Client) repo.ScenarioRepo {
	return &scenarioRepo{client: client}
}

func (repo *scenarioRepo) Find(ctx context.Context, dbname, name string) ([]restify.Scenario, error) {
	var filter bson.M
	if name != "" {
		filter = bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
	}
	res, err := repo.client.
		Database(dbname).
		Collection("scenario").
		Find(ctx, filter)
	if err != nil {
		return nil, exception.New(500, "Failed to find scenario from database. %s", err.Error())
	}

	scenarios := []restify.Scenario{}
	defer res.Close(ctx)
	for res.Next(ctx) {
		scenario := scenario.New()
		err = res.Decode(scenario)
		if err != nil {
			return nil, exception.New(500, "Failed to decode scenario from database")
		}

		scenarios = append(scenarios, scenario)
	}

	return scenarios, nil
}

func (repo *scenarioRepo) Get(ctx context.Context, dbname, id string) (restify.Scenario, error) {
	filter := bson.M{"_id": id}
	res := repo.client.Database(dbname).Collection("scenario").FindOne(ctx, filter)

	scenario := scenario.New()
	err := res.Decode(scenario)
	if err != nil {
		return nil, exception.New(500, "Failed to decode scenario from database")
	}

	return scenario, nil
}

func (repo *scenarioRepo) Add(ctx context.Context, dbname string, scenario restify.Scenario) (restify.Scenario, error) {
	id := primitive.NewObjectIDFromTimestamp(time.Now()).Hex()
	scenario.Set().ID(id)

	_, err := repo.client.
		Database(dbname).
		Collection("scenario").
		InsertOne(ctx, scenario)
	if err != nil {
		return nil, exception.New(500, "Failed to add scenario to DB: %s, Scenario: %s")
	}

	return scenario, nil
}

func (repo *scenarioRepo) Update(ctx context.Context, dbname, id string, scenario restify.Scenario) (restify.Scenario, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": scenario}
	_, err := repo.client.
		Database(dbname).
		Collection("scenario").
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, exception.New(500, "Failed to add scenario to DB: %s, Scenario: %s")
	}

	return scenario, nil
}

func (repo *scenarioRepo) Delete(ctx context.Context, dbname, id string) error {
	filter := bson.M{"_id": id}
	res, err := repo.client.
		Database(dbname).
		Collection("scenario").
		DeleteOne(ctx, filter)

	if res.DeletedCount <= 0 {
		return exception.New(404, "Failed, scenario with id: %s is not found", id)
	}

	if err != nil {
		return exception.New(500, "Failed to delete scenario. %s", err.Error())
	}

	return nil
}
