package main

import (
	"io/ioutil"

	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-lambda-graphql/config"
	"go-lambda-graphql/gql"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/volatiletech/sqlboiler/boil"
	"xi2.org/x/httpgzip"
)

var schema *graphql.Schema

func checkPanicError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func init() {
	// opentracing.SetGlobalTracer(graphql.Tracer)
	rawSchema, err := ioutil.ReadFile("gql/schema.gql")
	checkPanicError(err)
	schema = graphql.MustParseSchema(string(rawSchema), &gql.Resolver{})
	boil.DebugMode = !config.IsProduction
	db, err := sql.Open("postgres", config.ConnectionString)
	checkPanicError(err)
	boil.SetDB(db)
}

func main() {
	router := httprouter.New()

	// routes
	router.Handler("POST", "/query", httpgzip.NewHandler(&relay.Handler{Schema: schema}, nil))
	router.NotFound = httpgzip.NewHandler(http.FileServer(http.Dir(config.Directory)), nil).ServeHTTP

	s := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("go server listening on port " + config.Port)
	log.Fatal(s.ListenAndServe())
}
