package main

import (
	"fmt"
	"nafisah/delivery/gin_adapter"
	"nafisah/delivery/gin_adapter/controller"
	"nafisah/pkg"
	"nafisah/usecase"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_INTERNAL_ADDRESS"),
	})

	db := pkg.ConnectDatabase()
	usecases := usecase.NewUsecases(db, rdb)
	controllers := controller.NewControllers(usecases)
	r := gin_adapter.SetupRouter(controllers)

	r.Run(fmt.Sprintf(":%s", os.Getenv("API_PORT")))
}
