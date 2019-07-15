package repo

import (
	"context"

	restify "github.com/bastianrob/go-restify"
)

//ScenarioRepo contract
type ScenarioRepo interface {
	Find(ctx context.Context, dbname, name string) ([]restify.Scenario, error)
	Get(ctx context.Context, dbname, id string) (restify.Scenario, error)

	Add(ctx context.Context, dbname string, scenario restify.Scenario) (restify.Scenario, error)
	Update(ctx context.Context, dbname, id string, scenario restify.Scenario) (restify.Scenario, error)
	Delete(ctx context.Context, dbname, id string) error
}
