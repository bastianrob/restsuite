package main

import (
	"log"
	"net/http"
	"os"

	handlerV1 "github.com/bastianrob/restsuite/pkg/controller/v1"
	serviceV1 "github.com/bastianrob/restsuite/pkg/service/v1"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bastianrob/restsuite/pkg/repo/mongorepo"
)

func main() {
	mongoConn := os.Getenv("MONGO_CONN")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoConn))
	if err != nil {
		log.Fatal(err)
	}

	scenarioRepo := mongorepo.NewScenarioRepo(mongoClient)
	scenarioService := serviceV1.NewScenarioService(scenarioRepo)
	scenarioHandler := handlerV1.NewScenarioHandler(scenarioService)

	routes := httprouter.New()
	healthcheck(routes)
	routesV1(routes, scenarioHandler)

	httphandler := handlers.CombinedLoggingHandler(os.Stdout, routes)
	httphandler = handlers.RecoveryHandler()(httphandler)
	log.Fatal(http.ListenAndServe(":3000", httphandler))

}
