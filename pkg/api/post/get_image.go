package post

import (
	"io"
	"net/http"
	"qtbooru/pkg/fixedmaps"
	"strings"
)

const cacheLength = 100
var imageCache fixedmaps.FixedSizeMap[URL, *[]byte] = *fixedmaps.NewFixedSizeMap[URL, *[]byte](cacheLength)

func (f *File) Get(client *http.Client) (*[]byte, error) {
	if val, has := imageCache.Get(f.URL); has {
		return val, nil
	} else {
		req, err := http.NewRequest("GET", strings.Replace(string(f.URL), "localhost", "loki2", 1), nil)
		if err != nil { return nil, err }
		resp, err := client.Do(req)
		if err != nil { return nil, err }
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil { return nil, err }
		imageCache.Set(f.URL, &body)
		return &body, nil
	}
}

