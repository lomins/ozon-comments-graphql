package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lomins/ozon-comments-graphql/config"
	"github.com/lomins/ozon-comments-graphql/internal/graph"
	"github.com/lomins/ozon-comments-graphql/internal/graph/model"

	"github.com/lomins/ozon-comments-graphql/pkg/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	cfg.ParseFlags()

	var store storage.Storage

	switch cfg.App.Storage {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name, cfg.DB.Password)

		log.Printf("Storage type: %s", cfg.App.Storage)
		store, err = storage.NewPostgresStorage(dsn)
		if err != nil {
			log.Fatalf("Не удалось подключиться к Postgres: %s", err)
		}
		defer store.Close()
	case "inmemory":
		log.Printf("Storage type: %s", cfg.App.Storage)
		store = storage.NewInMemoryStorage()
	default:
		log.Printf("unknown storage type: %s, defaulting to inmemory", cfg.App.Storage)
		store = storage.NewInMemoryStorage()
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Storage: store, Comments: make(map[string]chan *model.Comment),
	}}))

	srv.AddTransport(&transport.Websocket{})

	startServer(cfg.App.Port, srv)
}

func startServer(port int, srv *handler.Server) {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	address := fmt.Sprintf(":%d", port)
	log.Printf("connect to http://localhost%s/ for GraphQL playground", address)
	log.Fatal(http.ListenAndServe(address, nil))
	address := fmt.Sprintf(":%d", port)
	log.Printf("connect to http://localhost%s/ for GraphQL playground", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
