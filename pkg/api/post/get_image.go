package post

import (
	"io"
	"net/http"
)

func (f *File) Get(client *http.Client) (*[]byte, error) {
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil { return nil, err }
	resp, err := client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil { return nil, err }
	return &body, nil
}

