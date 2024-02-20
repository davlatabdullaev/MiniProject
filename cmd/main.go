package main

import (
	"context"
	"test/config"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	// newRedis := redis.New(cfg)

	pgStore, err := postgres.New(context.Background(), cfg, log, newRedis)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return
	}
	defer pgStore.Close()

	services := service.New(pgStore, log, newRedis)

	server := api.New(services, log)

	log.Info("Service is running on", logger.Int("port", 8080))
	if err = server.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
