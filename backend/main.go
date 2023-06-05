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

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("BACKEND_REDIS_INTERNAL_ADDRESS"),
	})

	// clusterURLs = []string{"http://elasticsearch:9200", "http://elasticsearch:9300"}
	clusterURLs := []string{"http://localhost:9201", "http://localhost:9301"}
	username := "elastic"
	password := "kopinikmatnyamandilambung"
	cfg := elasticsearch.Config{
		Addresses: clusterURLs,
		Username:  username,
		Password:  password,
	}
	embeddingsServiceBaseURL := "http://localhost:8000"
	embeddings := embeddings_helper.NewEmbeddings(embeddingsServiceBaseURL)
	profileIndex := elasticsarch_helper.NewElasticsearchProfileIndex(cfg, embeddings)
	profileIndex.CreateIndexIfNotExists(ctx)

	db := gorm_helper.ConnectDatabase()
	usecases := usecase.NewUsecases(db, rdb, profileIndex)
	controllers := controller.NewControllers(usecases)
	r := gin_adapter.SetupRouter(controllers)

	r.Run(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT")))
}
