package main

import (
	"github.com/excellentprogrammer/goLang_site_work/shipping/proto"
	"github.com/excellentprogrammer/goLang_site_work/store/internal/platform/broker"
	"github.com/excellentprogrammer/goLang_site_work/store/internal/platform/config"
	"github.com/excellentprogrammer/goLang_site_work/store/internal/platform/redis"
	"github.com/excellentprogrammer/goLang_site_work/store/internal/service"
	"github.com/excellentprogrammer/goLang_site_work/store/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	gmbroker "github.com/micro/go-micro/broker"
	"time"
)

func main() {

	if err := gmbroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := gmbroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	itemShippedChannel := make(chan *shipping.ItemShippedEvent)
	broker.CreateEventConsumer(itemShippedChannel)

	repo := redis.NewStoreRepository(":6379")

	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	store.RegisterStoreHandler(svc.Server(), service.NewStoreService(repo, itemShippedChannel))

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
