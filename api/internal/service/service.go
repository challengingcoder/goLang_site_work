package service

import (
	stderrors "errors"
	"github.com/excellentprogrammer/goLang_site_work/catalog/proto"
	"github.com/excellentprogrammer/goLang_site_work/shipping/proto"
	"github.com/excellentprogrammer/goLang_site_work/store/proto"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"net/http"
)

const (
	catalogService   = "go.shopping.srv.catalog"
	shippingService  = "go.shopping.srv.shipping"
	storeService = "go.shopping.srv.store"
)

type CommerceService struct {
	storeClient store.StoreClient
	shippingClient  shipping.ShippingClient
	catalogClient   catalog.CatalogClient
}

type catalogResults struct {
	catalogResponse *catalog.DetailResponse
	err             error
}

type storeResults struct {
	storeResponse *store.DetailsResponse
	err               error
}

func NewCommerceService(c client.Client) *CommerceService {
	return &CommerceService{
		storeClient: store.NewStoreClient(storeService, c),
		shippingClient:  shipping.NewShippingClient(shippingService, c),
		catalogClient:   catalog.NewCatalogClient(catalogService, c),
	}
}

func (cs *CommerceService) GetProductDetails(request *restful.Request, response *restful.Response) {

	sku := request.PathParameter("sku")
	log.Logf("Received request for product details: %s", sku)
	ctx := context.Background()
	catalogCh := cs.getCatalogDetails(ctx, sku)
	storeCh := cs.getStoreDetails(ctx, sku)

	catalogReply := <-catalogCh
	if catalogReply.err != nil {
		writeError(response, catalogReply.err)
		return
	}

	storeReply := <-storeCh
	if storeReply.err != nil {
		writeError(response, storeReply.err)
		return
	}
	product := catalogReply.catalogResponse.Product

	details := productDetails{
		SKU:            product.Sku,
		StockRemaining: storeReply.storeResponse.Details.StockRemaining,
		Manufacturer:   product.Manufacturer,
		Price:          product.Price,
		Model:          product.Model,
		Name:           product.Name,
		Description:    product.Description,
	}
	response.WriteEntity(details)
}

func (cs *CommerceService) getCatalogDetails(ctx context.Context, sku string) chan catalogResults {
	ch := make(chan catalogResults, 1)

	go func() {
		res, err := cs.catalogClient.GetProductDetails(ctx, &catalog.DetailRequest{Sku: sku})
		ch <- catalogResults{catalogResponse: res, err: err}
	}()

	return ch
}

func (cs *CommerceService) getStoreDetails(ctx context.Context, sku string) chan storeResults {
	ch := make(chan storeResults, 1)

	go func() {
		res, err := cs.storeClient.GetStoreDetails(ctx, &store.DetailsRequest{Sku: sku})
		ch <- storeResults{storeResponse: res, err: err}
	}()

	return ch
}

func writeError(response *restful.Response, err error) {
	realError := errors.Parse(err.Error())
	if realError != nil {
		response.WriteError(int(realError.Code), stderrors.New(realError.Detail))
		return
	}
	response.WriteError(http.StatusInternalServerError, err)

}
