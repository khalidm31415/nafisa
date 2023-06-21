package main

import (
	"backend/delivery/gin_adapter"
	"backend/delivery/gin_adapter/controller"
	elasticsarch_helper "backend/package_helper/elasticsearch_helper"
	"backend/package_helper/embeddings_helper"
	"backend/package_helper/gorm_helper"
	"backend/usecase"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("BACKEND_REDIS_INTERNAL_ADDRESS"),
	})

	clusterURLs := strings.Split(os.Getenv("BACKEND_ELASTICSEARCH_CLUSTER_URLS"), ",")
	username := "elastic"
	password := os.Getenv("BACKEND_ELASTICSEARCH_PASSWORD")
	cfg := elasticsearch.Config{
		Addresses: clusterURLs,
		Username:  username,
		Password:  password,
	}
	embeddingsServiceBaseURL := os.Getenv("BACKEND_EMBEDDINGS_API_BASE_URL")
	embeddings := embeddings_helper.NewEmbeddings(embeddingsServiceBaseURL)
	profileIndex := elasticsarch_helper.NewElasticsearchProfileIndex(cfg, embeddings)
	profileIndex.CreateIndexIfNotExists(ctx)

	db := gorm_helper.ConnectDatabase()
	usecases := usecase.NewUsecases(db, rdb, profileIndex)
	googleOauthConfig := oauth2.Config{
		ClientID:     os.Getenv("BACKEND_GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("BACKEND_GOOGLE_OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BACKEND_GOOGLE_OAUTH_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	controllers := controller.NewControllers(googleOauthConfig, usecases)
	r := gin_adapter.SetupRouter(controllers)

	r.Run(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")))
}
