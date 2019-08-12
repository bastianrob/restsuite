package v1

import (
	"context"

	restify "github.com/bastianrob/go-restify"
	"github.com/bastianrob/restsuite/pkg/repo"
	"github.com/bastianrob/restsuite/pkg/service"
)

type resultService struct {
	repo repo.ResultRepo
}

//NewResultService v1
func NewResultService(repo repo.ResultRepo) service.ResultService {
	return &resultService{repo: repo}
}

func (svc *resultService) Find(ctx context.Context, params ...map[string]interface{}) ([]restify.TestResult, int64, error) {
	ctx, dbname, err := service.GetOrganizationName(ctx)
	if err != nil {
		return nil, 0, err
	}

	return svc.repo.Find(ctx, dbname, params...)
}
