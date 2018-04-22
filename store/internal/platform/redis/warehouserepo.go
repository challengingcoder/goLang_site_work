package redis

import (
	"fmt"
	"github.com/excellentprogrammer/goLang_site_work/store/proto"
	"github.com/garyburd/redigo/redis"
)

// StoreRepository represents a redis repository over store data
type StoreRepository struct {
	redisDialString string
}

// NewStoreRepository creates a new store repo
func NewStoreRepository(redisDialString string) *StoreRepository {
	return &StoreRepository{redisDialString: redisDialString}
}

// GetStoreDetails queries the information about physical inventory in the store for a given SKU
func (r *StoreRepository) GetStoreDetails(sku string) (details *store.StoreDetails, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	storeKey := fmt.Sprintf("store:%s", sku)
	res, err := redis.Values(c.Do("HGETALL", storeKey))
	var itemDetails redisStoreDetails
	err = redis.ScanStruct(res, &itemDetails)
	if err != nil {
		return nil, err
	}
	stockKey := fmt.Sprintf("store:%s:stock", sku)
	stockCount, err := redis.Int(c.Do("GET", stockKey))
	if err != nil {
		return nil, err
	}
	details = &store.StoreDetails{
		Sku:            itemDetails.SKU,
		Manufacturer:   itemDetails.Manufacturer,
		ModelNumber:    itemDetails.ModelNumber,
		StockRemaining: uint32(stockCount),
	}
	return details, nil
}

// SkuExists indicates whether the SKU exists in the store inventory (regardless of in-stock quantity)
func (r *StoreRepository) SkuExists(sku string) (exists bool, err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return false, err
	}
	defer c.Close()
	storeKey := fmt.Sprintf("store:%s", sku)
	exists, err = redis.Bool(c.Do("EXISTS", storeKey))
	return exists, err
}

// DecrementStock will reduce the on-hand quantity of a SKU by 1
func (r *StoreRepository) DecrementStock(sku string) (err error) {
	c, err := redis.Dial("tcp", r.redisDialString)
	if err != nil {
		return err
	}
	defer c.Close()
	storeKey := fmt.Sprintf("store:%s:stock", sku)
	_, err = c.Do("INCRBY", storeKey, "-1")
	if err != nil {
		return err
	}

	return nil
}

type redisStoreDetails struct {
	SKU          string `redis:"sku"`
	Manufacturer string `redis:"mfr"`
	ModelNumber  string `redis:"model"`
}
