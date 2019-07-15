package service

import (
	"context"

	restify "github.com/bastianrob/go-restify"
	"github.com/bastianrob/restsuite/pkg/ctxkey"
	"github.com/bastianrob/restsuite/pkg/exception"
)

//ScenarioService contract
type ScenarioService interface {
	Find(ctx context.Context, name string) ([]restify.Scenario, error)
	Get(ctx context.Context, id string) (restify.Scenario, error)

	Add(ctx context.Context, scenario restify.Scenario) (restify.Scenario, error)
	Update(ctx context.Context, id string, scenario restify.Scenario) (restify.Scenario, error)
	Delete(ctx context.Context, id string) error
}

//GetOrganizationName from context
func GetOrganizationName(ctx context.Context) (context.Context, string, error) {
	dbname, ok := ctx.Value(ctxkey.OrganizationName).(string)
	if !ok || dbname == "" {
		return ctx, "", exception.New(400, "Organization name is not found in request context")
	}

	return ctx, dbname, nil
}
