package shared

import "time"

// GoxCache - actual cache struct
type GoxCache struct {
	ID         string                    `json:"id"`
	Name       string                    `json:"name"`
	Storage    map[string]interface{}    `json:"data"`
	Expiration map[string]*GoxExpiration `json:"expiration"`
	SyncKey    *string                   `json:"key"`
}

// GoxSyncCache - holds the content of a gox sync cache
type GoxSyncCache struct {
	Cache *GoxCache `json:"content"`
}

// GoxExpiration - expiration struct container
type GoxExpiration struct {
	CreatedAt  time.Time
	Expiration time.Duration
}

// GoxModule - module with cache functions
type GoxModule struct {
	Cache *GoxCache
}
