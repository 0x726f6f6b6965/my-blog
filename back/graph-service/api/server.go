package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/0x726f6f6b6965/my-blog/graph-service/graph"
	"github.com/0x726f6f6b6965/my-blog/graph-service/internal/client"
	"github.com/0x726f6f6b6965/my-blog/graph-service/internal/config"
	"github.com/0x726f6f6b6965/my-blog/graph-service/internal/middleware"
	"github.com/0x726f6f6b6965/my-blog/lib/logger"
	pbBlog "github.com/0x726f6f6b6965/my-blog/protos/blog/v1"
	pbSearch "github.com/0x726f6f6b6965/my-blog/protos/search/v1"
	pbUser "github.com/0x726f6f6b6965/my-blog/protos/user/v1"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const defaultPort = 8080

func main() {
	godotenv.Load()
	path := os.Getenv("CONFIG")
	var cfg config.Config
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read yaml error", err)
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("unmarshal yaml error", err)
		return
	}
	zaplog, cleanup, err := logger.NewLogger(&cfg.Log)
	if err != nil {
		log.Fatal("create log error", err)
		return
	}
	defer cleanup()
	handler, err := NewGraphQLHandler(&cfg, zaplog)
	if err != nil {
		return
	}
	port := cfg.Rest.Port
	if port == 0 {
		port = defaultPort
	}

	log.Printf("connect to http://%s:%d/ for GraphQL playground", cfg.Rest.Host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Rest.Host, port), handler))
}

// NewGraphQLHandler returns handler for GraphQL application
func NewGraphQLHandler(cfg *config.Config, zaplog *zap.Logger) (*chi.Mux, error) {
	// create a new router
	var router *chi.Mux = chi.NewRouter()

	rds := client.InitRedisClient(&cfg.Redis)
	defer rds.Close()

	connUser, err := client.NewGrpcConn(
		context.Background(),
		fmt.Sprintf("%s:%d", cfg.Clients.User.Grpc.Host, cfg.Clients.User.Grpc.Port),
		cfg.Clients.User.ConnectionTimeout)
	if err != nil {
		zaplog.Error("connect user service error", zap.Error(err))
		return nil, errors.New("connect user service error")
	}

	connBlog, err := client.NewGrpcConn(
		context.Background(),
		fmt.Sprintf("%s:%d", cfg.Clients.Blog.Grpc.Host, cfg.Clients.Blog.Grpc.Port),
		cfg.Clients.Blog.ConnectionTimeout)
	if err != nil {
		zaplog.Error("connect blog service error", zap.Error(err))
		return nil, errors.New("connect blog service error")
	}

	connSearch, err := client.NewGrpcConn(
		context.Background(),
		fmt.Sprintf("%s:%d", cfg.Clients.Search.Grpc.Host, cfg.Clients.Search.Grpc.Port),
		cfg.Clients.Search.ConnectionTimeout)
	if err != nil {
		zaplog.Error("connect search service error", zap.Error(err))
		return nil, errors.New("connect search service error")
	}

	resolver := &graph.Resolver{
		UserService:   pbUser.NewUserServiceClient(connUser.GetConn()),
		BlogService:   pbBlog.NewBlogServiceClient(connBlog.GetConn()),
		SearchService: pbSearch.NewSearchServiceClient(connSearch.GetConn()),
		Log:           zaplog,
	}

	// use the middleware component
	router.Use(middleware.NewMiddleware(rds))

	// create a GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// assign some handlers for the GraphQL server
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	// return the handler
	return router, nil
}
