package v1

import (
	"context"
	"regexp"
	"strings"

	thennable "github.com/bastianrob/go-thennable"
	"github.com/bastianrob/restsuite/model"
	"github.com/bastianrob/restsuite/pkg/repo"
	"github.com/bastianrob/restsuite/pkg/service"
)

type orgService struct {
	repo repo.OrganizationRepo
}

//NewOrganizationService v1
func NewOrganizationService(repo repo.OrganizationRepo) service.OrganizationService {
	return &orgService{repo: repo}
}

func (svc *orgService) Get(ctx context.Context) (model.Organization, error) {
	res, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(svc.repo.Get).
		End()

	if err != nil {
		return model.Organization{}, err
	}

	return res[0].(model.Organization), nil
}

func (svc *orgService) Add(ctx context.Context, name string) (model.Organization, error) {
	org := model.Organization{Name: name}

	//Create alias
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	name = strings.ToLower(name)
	name = reg.ReplaceAllString(name, "")
	org.Alias = name

	//Set subscription as free tier
	org.Subscription.Type = "Free"
	org.Subscription.Limit.Endpoint = 10
	org.Subscription.Limit.Scenario = 5
	org.Subscription.Limit.Cases = 25
	org.Subscription.Limit.Data = 1e+8 //100MB

	return svc.repo.Add(ctx, org)
}

func (svc *orgService) Update(ctx context.Context, update model.Organization) (model.Organization, error) {
	res, err := thennable.Start(ctx).
		Then(service.GetOrganizationName).
		Then(svc.repo.Get).
		Then(func(existing model.Organization) (context.Context, model.Organization, error) {
			existing.Endpoints = update.Endpoints
			existing.Environments = update.Environments

			return ctx, existing, nil
		}).
		Then(svc.repo.Update).
		End()

	if err != nil {
		return model.Organization{}, err
	}

	return res[0].(model.Organization), nil
}
