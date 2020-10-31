package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (gox *GoxModule) deleteCacheAndExpiration(name string) {
	gox.Cache.Storage[name] = nil
	gox.Cache.Expiration[name] = nil
	delete(gox.Cache.Storage, name)
	delete(gox.Cache.Expiration, name)
}

func (gox *GoxModule) sendPost(url string, data []byte) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", *gox.Cache.SyncKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	returned, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return returned, nil
}

func (gox *GoxModule) syncWithBackend() (*GoxSyncCache, error) {
	objectData := make(map[string]interface{})
	if gox.Cache.Storage != nil && gox.Cache.Expiration != nil {
		objectData["cache_data"] = gox.Cache.Storage
		objectData["cache_expiration"] = gox.Cache.Expiration
		objectData["cache_name"] = gox.Cache.Name
		objectData["cache_id"] = gox.Cache.ID
	}

	data, err := json.Marshal(objectData)
	if err != nil {
		return nil, err
	}

	resp, err := gox.sendPost(fmt.Sprintf("%s/sync", SyncAPI), data)
	if err != nil {
		return nil, err
	}

	formatted := &GoxSyncCache{}
	json.Unmarshal(resp, &formatted)

	return formatted, nil
}

func (gox *GoxModule) deleteExpiredBackend(expires []string) (*GoxSyncCache, error) {
	objectData := map[string]interface{}{
		"items":      expires,
		"cache_name": gox.Cache.Name,
		"cache_id":   gox.Cache.ID,
	}

	data, err := json.Marshal(objectData)
	if err != nil {
		return nil, err
	}

	resp, err := gox.sendPost(fmt.Sprintf("%s/delete", SyncAPI), data)
	if err != nil {
		return nil, err
	}

	formatted := &GoxSyncCache{}
	json.Unmarshal(resp, &formatted)

	return formatted, nil
}

func (gox *GoxModule) startProcess(action func(), length time.Duration) chan bool {
	stop := make(chan bool)
	action()
	go func() {
		for {
			action()
			select {
			case <-time.After(length):
			case <-stop:
				return
			}
		}
	}()
	return stop
}
