package shared

import (
	"errors"
	"fmt"
	"time"
)

// Init - initialize the cache object
func (gox *GoxModule) Init(id string, name string) {
	gox.Cache = &GoxCache{
		ID:         id,
		Name:       name,
		Storage:    make(map[string]interface{}),
		Expiration: make(map[string]*GoxExpiration),
		SyncKey:    nil,
	}
}

// InitWithSync - initialize the cache object with api key
func (gox *GoxModule) InitWithSync(id string, name string, apiKey string) {
	gox.Cache = &GoxCache{
		ID:         id,
		Name:       name,
		Storage:    make(map[string]interface{}),
		Expiration: make(map[string]*GoxExpiration),
		SyncKey:    &apiKey,
	}

	updateLocal := func() {
		data, err := gox.syncWithBackend()
		if err != nil {
			fmt.Println(fmt.Sprintf("An error occurred while syncing: %s", err.Error()))
		}
		if data != nil {
			if *data.Cache.SyncKey == *gox.Cache.SyncKey {
				if data.Cache.ID == gox.Cache.ID {
					gox.Cache.Storage = data.Cache.Storage
					gox.Cache.Expiration = data.Cache.Expiration
				}
			}
		}
	}

	gox.startProcess(updateLocal, 5*time.Second)
}

// SetItem - add a new item into the cache without overriding
func (gox *GoxModule) SetItem(name string, value interface{}, expires time.Duration) error {
	gox.Cache.Storage[name] = value
	gox.Cache.Expiration[name] = &GoxExpiration{
		CreatedAt:  time.Now(),
		Expiration: expires,
	}

	return nil
}

// GetItem - returns an item and error if the item does not exist or was expired
func (gox *GoxModule) GetItem(name string, purgeIfExpired bool) (interface{}, error) {
	val, ok := gox.Cache.Storage[name]
	if !ok {
		return nil, errors.New("gox item does not exist")
	}

	expVal, ok := gox.Cache.Expiration[name]
	if !ok {
		return nil, errors.New("gox item does not contain an expiration")
	}

	if purgeIfExpired {
		if expVal.CreatedAt.Add(expVal.Expiration).Before(time.Now()) {
			gox.deleteCacheAndExpiration(name)
			return nil, errors.New("gox item has expired")
		}
	}

	return val, nil
}

// RemoveItem - removes an item from the gox cache
func (gox *GoxModule) RemoveItem(name string) error {
	_, ok := gox.Cache.Storage[name]
	if !ok {
		return errors.New("gox item does not exist to delete")
	}

	_, expOk := gox.Cache.Expiration[name]
	if !expOk {
		return errors.New("gox item does not contain an expiration to delete")
	}

	gox.deleteCacheAndExpiration(name)
	return nil
}

// UpdateItem - update an item's value in the gox cache
func (gox *GoxModule) UpdateItem(name string, newVal interface{}, shouldResetExpiration bool) error {
	_, ok := gox.Cache.Storage[name]
	if !ok {
		return errors.New("gox item does not exist to update")
	}

	_, expOk := gox.Cache.Expiration[name]
	if !expOk {
		return errors.New("gox item does not contain an expiration to update")
	}

	gox.Cache.Storage[name] = newVal
	if shouldResetExpiration {
		gox.Cache.Expiration[name].CreatedAt = time.Now()
	}

	return nil
}

// PurgeExpiredItems - returns a count of all expired items
func (gox *GoxModule) PurgeExpiredItems() ([]string, int) {
	var expiredItems []string

	for item, itemExp := range gox.Cache.Expiration {
		if itemExp.CreatedAt.Add(itemExp.Expiration).Before(time.Now()) {
			expiredItems = append(expiredItems, item)
			gox.deleteCacheAndExpiration(item)
		}
	}

	gox.deleteExpiredBackend(expiredItems)

	return expiredItems, len(expiredItems)
}
