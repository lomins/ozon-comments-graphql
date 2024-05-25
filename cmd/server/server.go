package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lomins/ozon-comments-graphql/database"
	"github.com/lomins/ozon-comments-graphql/internal/graph"
	"github.com/lomins/ozon-comments-graphql/internal/graph/model"
	"github.com/lomins/ozon-comments-graphql/pkg/storage"
)

const defaultPort = "8080"

func main() {
	db := database.InitDB()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	storageType := os.Getenv("STORAGE_TYPE")
	var store storage.Storage

	if storageType == "postgres" {
		db := database.InitDB()
		defer db.Close()
		store = storage.NewPostgresStorage(db)
	} else {
		store = storage.NewInMemoryStorage()
	}

	resolver := &graph.Resolver{Storage: store, Comments: make(map[string]chan *model.Comment)}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
