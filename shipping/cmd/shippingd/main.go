package main

import (
	"github.com/excellentprogrammer/goLang_site_work/shipping/internal/platform/broker"
	"github.com/excellentprogrammer/goLang_site_work/shipping/internal/platform/config"
	"github.com/excellentprogrammer/goLang_site_work/shipping/internal/platform/redis"
	"github.com/excellentprogrammer/goLang_site_work/shipping/internal/service"
	"github.com/excellentprogrammer/goLang_site_work/shipping/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	gmbroker "github.com/micro/go-micro/broker"
	"log"
	"time"
)

func main() {
	if err := gmbroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := gmbroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	repo := redis.NewRedisRepository(":6379")
	publisher := broker.NewEventPublisher()
	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	shipping.RegisterShippingHandler(svc.Server(), service.NewShippingService(repo, publisher))

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
