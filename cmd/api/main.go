package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"

	"github.com/racoon-proger/wb-l0/internal/cache"
	"github.com/racoon-proger/wb-l0/internal/domain"
	server "github.com/racoon-proger/wb-l0/internal/server/http"
	"github.com/racoon-proger/wb-l0/internal/service"
	"github.com/racoon-proger/wb-l0/storage"
)

type configuration struct {
	DatabasHost      string `envconfig:"DB_HOST" required:"true"`
	DatabasePort     int    `envconfig:"DB_PORT" required:"true"`
	DatabaseUser     string `envconfig:"DB_USER" required:"true"`
	DatabasePassword string `envconfig:"DB_PASSWORD" required:"true"`
	DatabaseName     string `envconfig:"DB_NAME" required:"true"`
	ServerPort       int    `envconfig:"SERVER_PORT" required:"true"`
	NatsAddr         string `envconfig:"NATS_ADDR" default:"nats://127.0.0.1:4222"`
	NatsQueue        string `envconfig:"NATS_QUEUE" default:"order"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabasHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	nc, err := nats.Connect(cfg.NatsAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	storage := storage.NewStorage(db)
	orders, err := storage.GetOrders(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cache := cache.NewCache()
	svc := service.NewService(cache, storage)
	nc.Subscribe(cfg.NatsQueue, func(msg *nats.Msg) {
		var order domain.Order
		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("failed to parse data: %s", err)
			return
		}
		err = svc.CreateOrder(ctx, &order)
		if err != nil {
			log.Printf("failed to create order: %s", err)
			return
		}
		msg.Ack()
	})

	for i := range orders {
		cache.SetOrder(&orders[i])
	}
	server := server.NewServer(svc, cfg.ServerPort)
	http.HandleFunc("/", server.ServeHTML)
	http.HandleFunc("/get-order", server.GetOrder)
	go func() {
		err := server.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
}
