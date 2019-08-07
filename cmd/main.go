package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	handlerV1 "github.com/bastianrob/restsuite/pkg/controller/v1"
	serviceV1 "github.com/bastianrob/restsuite/pkg/service/v1"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

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

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	scenarioRepo := mongorepo.NewScenarioRepo(mongoClient)
	scenarioService := serviceV1.NewScenarioService(scenarioRepo)
	scenarioHandler := handlerV1.NewScenarioHandler(scenarioService)

	orgRepo := mongorepo.NewOrganizationRepo(mongoClient)
	orgService := serviceV1.NewOrganizationService(orgRepo)
	orgHandler := handlerV1.NewOrganizationHandler(orgService)

	routes := httprouter.New()
	healthcheck(routes)
	scenarioV1(routes, scenarioHandler)
	orgV1(routes, orgHandler)

	httphandler := handlers.CombinedLoggingHandler(os.Stdout, routes)
	httphandler = handlers.RecoveryHandler()(httphandler)
	httphandler = cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:*",
			"http://lapelio.com",
			"https://lapelio.com",
			"http://*.lapelio.com",
			"https://*.lapelio.com",
		},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			"RUN", //custom method
		},
		AllowedHeaders: []string{"*"},
	}).Handler(httphandler)
	log.Println("Serving :7001")
	log.Fatal(http.ListenAndServe(":7001", httphandler))
}
