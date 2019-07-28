package repo

import (
	"context"

	"github.com/bastianrob/restsuite/model"
)

//OrganizationRepo contract
type OrganizationRepo interface {
	Get(ctx context.Context, dbname string) (model.Organization, error)
	Add(ctx context.Context, org model.Organization) (model.Organization, error)
	Update(ctx context.Context, org model.Organization) (model.Organization, error)
}
