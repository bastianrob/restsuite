package mongorepo

import (
	"context"
	"time"

	"github.com/bastianrob/arrayutil"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/bastianrob/restsuite/model"
	"github.com/bastianrob/restsuite/pkg/exception"
	"github.com/bastianrob/restsuite/pkg/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

type orgRepo struct {
	client *mongo.Client
}

//NewOrganizationRepo using mongodb
func NewOrganizationRepo(client *mongo.Client) repo.OrganizationRepo {
	return &orgRepo{client: client}
}

func (repo *orgRepo) Get(ctx context.Context, dbname string) (model.Organization, error) {
	filter := bson.M{"alias": dbname}
	res := repo.client.
		Database(dbname).
		Collection("organization").
		FindOne(ctx, filter)

	org := model.Organization{}
	err := res.Decode(&org)
	if err != nil {
		return model.Organization{}, exception.New(500, "Failed to decode organization from database")
	}

	return org, nil
}

func (repo *orgRepo) Add(ctx context.Context, org model.Organization) (model.Organization, error) {
	org.ID = primitive.NewObjectIDFromTimestamp(time.Now()).Hex()

	dbs, err := repo.client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return model.Organization{}, exception.New(500, "Failed to find organizaton list. %s", err.Error())
	}

	if arrayutil.Contains(dbs, org.Alias) {
		return model.Organization{}, exception.New(409, "Organization already exists")
	}

	_, err = repo.client.
		Database(org.Alias).
		Collection("organization").
		InsertOne(ctx, org)
	if err != nil {
		return model.Organization{}, exception.New(500, "Failed to add organization to DB. Err: %s", err.Error())
	}

	return org, nil
}

func (repo *orgRepo) Update(ctx context.Context, org model.Organization) (model.Organization, error) {
	filter := bson.M{"alias": org.Alias}
	update := bson.M{"$set": org}
	_, err := repo.client.
		Database(org.Alias).
		Collection("organization").
		UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Organization{}, exception.New(500, "Failed to update organization to DB. Err: %s", err.Error())
	}

	return org, nil
}
