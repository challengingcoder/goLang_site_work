package main

import (
	"github.com/excellentprogrammer/goLang_site_work/api/internal/platform/config"
	"github.com/excellentprogrammer/goLang_site_work/api/internal/service"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	"github.com/micro/go-log"
	_ "github.com/micro/go-plugins/client/grpc"
)

const (
	serviceName = "go.shopping.api.v1.commerce"
)

func main() {
	webService := web.NewService(
		web.Name(serviceName),
		web.Version(config.Version),
	)

	webService.Init()
	handler := service.NewCommerceService(client.DefaultClient)

	wc := restful.NewContainer()
	ws := new(restful.WebService)

	ws.
		Path("/v1/commerce").
		Doc("Aggregated Commerce API").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/products/{sku}").To(handler.GetProductDetails)).
		Doc("Query product details")

	wc.Add(ws)
	webService.Handle("/", wc)
	if err := webService.Run(); err != nil {
		log.Fatal(err)
	}

}
