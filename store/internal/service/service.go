package service

import (
	"github.com/excellentprogrammer/goLang_site_work/shipping/proto"
	"github.com/excellentprogrammer/goLang_site_work/store/proto"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type storeService struct {
	repo     storeRepository
	shipChan chan *shipping.ItemShippedEvent
}

type storeRepository interface {
	GetStoreDetails(sku string) (details *store.StoreDetails, err error)
	SkuExists(sku string) (exists bool, err error)
	DecrementStock(sku string) (err error)
}

// NewStoreService returns an instance of a store handler
func NewStoreService(repo storeRepository, itemShippedChannel chan *shipping.ItemShippedEvent) store.StoreHandler {
	svc := &storeService{repo: repo, shipChan: itemShippedChannel}
	go svc.awaitItemShippedEvents()
	return svc
}

func (w *storeService) GetStoreDetails(ctx context.Context, request *store.DetailsRequest,
	response *store.DetailsResponse) error {

	if request == nil {
		return errors.BadRequest("", "Missing details request")
	}
	if len(request.Sku) < 6 {
		return errors.BadRequest("", "Invalid SKU")
	}
	exists, err := w.repo.SkuExists(request.Sku)
	if err != nil {
		return errors.InternalServerError(request.Sku, "Failed to check for SKU existence: %s", err)
	}
	if !exists {
		return errors.NotFound(request.Sku, "No such SKU")
	}

	details, err := w.repo.GetStoreDetails(request.Sku)
	if err != nil {
		return errors.InternalServerError(request.Sku, "Failed to query store details: %s", err)
	}

	response.Details = details

	return nil
}

func (w *storeService) awaitItemShippedEvents() {
	for shippedEvent := range w.shipChan {
		log.Logf("Received an item shipped event! %+v\n", shippedEvent)
		w.repo.DecrementStock(shippedEvent.Sku)
	}
}
