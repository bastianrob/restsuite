package service

import (
	"context"

	"github.com/bastianrob/restsuite/model"
)

//OrganizationService contract
type OrganizationService interface {
	Add(ctx context.Context, name string) (model.Organization, error)
	Get(ctx context.Context) (model.Organization, error)
	Update(ctx context.Context, org model.Organization) (model.Organization, error)
}
