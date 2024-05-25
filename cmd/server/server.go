package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jinzhu/gorm"
	"github.com/lomins/ozon-comments-graphql/config"
	"github.com/lomins/ozon-comments-graphql/internal/graph"

	"github.com/lomins/ozon-comments-graphql/pkg/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	var store storage.Storage

	switch cfg.App.Storage {
	case "postgres":
		db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name, cfg.DB.Password))
		if err != nil {
			log.Fatalf("could not connect to the database: %v", err)
		}
		defer db.Close()

		store = storage.NewPostgresStorage(db)
	case "inmemory":
		store = storage.NewInMemoryStorage()
	default:
		log.Fatalf("unknown storage type: %v", cfg.App.Storage)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Storage: store,
	}}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	port := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("connect to http://localhost%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
