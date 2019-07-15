package v1

import (
	"context"
	"time"

	thennable "github.com/bastianrob/go-thennable"
	"go.mongodb.org/mongo-driver/bson/primitive"

	restify "github.com/bastianrob/go-restify"
	"github.com/bastianrob/restsuite/pkg/repo"
	"github.com/bastianrob/restsuite/pkg/service"
)

type scenarioService struct {
	repo repo.ScenarioRepo
}

//NewScenarioService v1
func NewScenarioService(repo repo.ScenarioRepo) service.ScenarioService {
	return &scenarioService{repo}
}

func (svc *scenarioService) Find(ctx context.Context, name string) ([]restify.Scenario, error) {
	res, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(func(dbname string) (context.Context, string, string, error) {
			return ctx, dbname, name, nil
		}).
		Then(svc.repo.Get).
		End()

	if err != nil {
		return nil, err
	}

	return res[0].([]restify.Scenario), nil
}

func (svc *scenarioService) Get(ctx context.Context, id string) (restify.Scenario, error) {
	res, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(func(dbname string) (context.Context, string, string, error) {
			return ctx, dbname, id, nil
		}).
		Then(svc.repo.Get).
		End()

	if err != nil {
		return nil, err
	}

	return res[0].(restify.Scenario), nil
}

func (svc *scenarioService) Add(ctx context.Context, scenario restify.Scenario) (restify.Scenario, error) {
	_, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(func(ctx context.Context, dbname string) (context.Context, string, restify.Scenario, error) {
			id := primitive.NewObjectIDFromTimestamp(time.Now()).Hex()
			scenario.Set().ID(id)

			return ctx, dbname, scenario, nil
		}).
		Then(svc.repo.Add).
		End()

	return scenario, err
}

func (svc *scenarioService) Update(ctx context.Context, id string, scenario restify.Scenario) (restify.Scenario, error) {
	_, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(func(ctx context.Context, dbname string) (context.Context, string, string, restify.Scenario, error) {
			return ctx, dbname, id, scenario, nil
		}).
		Then(svc.repo.Update).
		End()

	return scenario, err
}

func (svc *scenarioService) Delete(ctx context.Context, id string) error {
	_, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(func(dbname string) (context.Context, string, string, error) {
			return ctx, dbname, id, nil
		}).
		Then(svc.repo.Delete).
		End()

	return err
}
