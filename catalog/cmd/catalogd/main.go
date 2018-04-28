package main

import (
	"github.com/excellentprogrammer/goLang_site_work/catalog/internal/platform/config"
	"github.com/excellentprogrammer/goLang_site_work/catalog/internal/platform/redis"
	"github.com/excellentprogrammer/goLang_site_work/catalog/internal/service"
	"github.com/excellentprogrammer/goLang_site_work/catalog/proto"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	svc := grpc.NewService(
		micro.Name(config.ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(config.Version),
	)
	svc.Init()

	redisCatalogRepository := redis.NewRedisRepository(":6379")
	catalog.RegisterCatalogHandler(svc.Server(), service.NewCatalogService(redisCatalogRepository))

	if err := svc.Run(); err != nil {
		panic(err)
	}
}
