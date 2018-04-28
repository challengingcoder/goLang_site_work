package service_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"

	stderrors "errors"
	"github.com/excellentprogrammer/goLang_site_work/shipping/proto"
	"github.com/excellentprogrammer/goLang_site_work/store/internal/service"
	"github.com/excellentprogrammer/goLang_site_work/store/proto"
	"github.com/micro/go-micro/errors"
	"net/http"
)

func TestStoreService_GetStoreDetails(t *testing.T) {
	Convey("Given a store service", t, func() {
		ctx := context.Background()
		stockChan := make(chan string)
		repo := &fakeRepo{stockChan: stockChan}
		shippedChannel := make(chan *shipping.ItemShippedEvent)
		svc := service.NewStoreService(repo, shippedChannel)

		Convey("requesting store details should invoke the repository", func() {
			repo.shouldFail = false
			var resp store.DetailsResponse
			err := svc.GetStoreDetails(ctx, &store.DetailsRequest{Sku: "111111"}, &resp)
			So(err, ShouldBeNil)
			So(resp.Details, ShouldNotBeNil)
			So(resp.Details.Manufacturer, ShouldEqual, "TOSHIBA")
			So(resp.Details.StockRemaining, ShouldEqual, 42)
			So(resp.Details.ModelNumber, ShouldEqual, "T-1000")
		})

		Convey("requesting store details should fail when the repo fails", func() {
			repo.shouldFail = true
			var resp store.DetailsResponse
			err := svc.GetStoreDetails(ctx, &store.DetailsRequest{Sku: "111111"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("requesting store details for a non-existent sku should fail with a 404", func() {
			repo.shouldFail = false
			var resp store.DetailsResponse
			err := svc.GetStoreDetails(ctx, &store.DetailsRequest{Sku: "nevergonnahappen"}, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("requesting store details with a nil request should fail", func() {
			repo.shouldFail = false
			var resp store.DetailsResponse
			err := svc.GetStoreDetails(ctx, nil, &resp)
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("requesting store details with a bad SKU should fail", func() {
			repo.shouldFail = false
			var resp store.DetailsResponse
			err := svc.GetStoreDetails(ctx, &store.DetailsRequest{Sku: "1111"}, &resp) // SKU is too short
			So(err, ShouldNotBeNil)
			realError := errors.Parse(err.Error())
			So(realError, ShouldNotBeNil)
			So(realError.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("when an item shipped event is received, it should ask the repo to decrement stock", func() {
			repo.shouldFail = false
			evt := &shipping.ItemShippedEvent{
				Sku:            "111111",
				ShippingMethod: shipping.ShippingMethod_SM_FEDEX,
				TrackingNumber: "abc1233",
			}
			shippedChannel <- evt
			s := <-stockChan
			So(s, ShouldEqual, "111111")
		})
	})
}

type fakeRepo struct {
	shouldFail bool
	stockChan  chan string
}

func (r *fakeRepo) GetStoreDetails(sku string) (details *store.StoreDetails, err error) {
	if r.shouldFail {
		return nil, stderrors.New("Faily Fail")
	}
	return &store.StoreDetails{
		ModelNumber:    "T-1000",
		StockRemaining: 42,
		Manufacturer:   "TOSHIBA",
		Sku:            "111111",
	}, nil
}

func (r *fakeRepo) DecrementStock(sku string) (err error) {
	r.stockChan <- sku
	return nil
}

func (r *fakeRepo) SkuExists(sku string) (exists bool, err error) {
	return sku == "111111", nil
}
