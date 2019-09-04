package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type bookrepo struct {
	server *url.URL
}

func (r *bookrepo) GetBookName(isbn string) (string, error) {
	// https://openlibrary.org/api/books?bibkeys=ISBN:9780134190440&format=json&jscmd=data

	bookUrl := *r.server
	bookUrl.Path = path.Join(r.server.Path, "/api/books")
	q := bookUrl.Query()
	q.Set("format", "json")
	q.Set("bibkeys", "ISBN:"+isbn)
	q.Set("jscmd", "data")
	bookUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, bookUrl.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code from openlibrary: %d", resp.StatusCode)
	}

	var respData map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&respData)

	name := respData["ISBN:"+isbn].(map[string]interface{})["title"].(string)

	return name, nil
}
