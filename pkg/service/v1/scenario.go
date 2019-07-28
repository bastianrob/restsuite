package v1

import (
	"bytes"
	"context"

	thennable "github.com/bastianrob/go-thennable"

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
		Then(func(ctx context.Context, dbname string) (context.Context, string, string, error) {
			return ctx, dbname, name, nil
		}).
		Then(svc.repo.Find).
		End()

	if err != nil {
		return nil, err
	}

	return res[0].([]restify.Scenario), nil
}

func (svc *scenarioService) Get(ctx context.Context, id string) (restify.Scenario, error) {
	provideData := func(ctx context.Context, dbname string) (context.Context, string, string, error) {
		return ctx, dbname, id, nil
	}

	res, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(provideData).
		Then(svc.repo.Get).
		End()

	if err != nil {
		return nil, err
	}

	return res[0].(restify.Scenario), nil
}

func (svc *scenarioService) Run(ctx context.Context, id string) (string, error) {
	buffer := bytes.Buffer{}
	_, err := thennable.Start(ctx, id).
		Then(svc.Get).
		Then(func(scenario restify.Scenario) error {
			scenario.Run(&buffer)
			return nil
		}).
		End()

	return buffer.String(), err
}

func (svc *scenarioService) Add(ctx context.Context, scenario restify.Scenario) (restify.Scenario, error) {
	_, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(svc.repo.Add).
		End()

	return scenario, err
}

func (svc *scenarioService) Update(ctx context.Context, id string, scenario restify.Scenario) (restify.Scenario, error) {
	provideUpdateData := func(ctx context.Context, dbname string) (context.Context, string, string, restify.Scenario, error) {
		return ctx, dbname, id, scenario, nil
	}

	_, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(provideUpdateData).
		Then(svc.repo.Update).
		End()

	return scenario, err
}

func (svc *scenarioService) Delete(ctx context.Context, id string) error {
	provideData := func(ctx context.Context, dbname string) (context.Context, string, string, error) {
		return ctx, dbname, id, nil
	}

	_, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(provideData).
		Then(svc.repo.Delete).
		End()

	return err
}
