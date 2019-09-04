package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	openlibrary, err := url.Parse("https://openlibrary.org")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.Methods(http.MethodGet).
		Path("/books/{isbn}").
		HandlerFunc(handleBookReq(&bookrepo{
			server: openlibrary,
		}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Running http server")
		err := server.ListenAndServe()
		fmt.Printf("server error: %+v\n", err)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	inter := <-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	fmt.Printf("Finished %s\n", inter)
}

func handleBookReq(repo *bookrepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isbn := mux.Vars(r)["isbn"]

		name, err := repo.getBookName(isbn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %+v\n", err)
			return
		}

		fmt.Fprintf(w, "Your book is: %s\n", name)
	}
}

type bookrepo struct {
	server *url.URL
}

func (r *bookrepo) getBookName(isbn string) (string, error) {
	// https://openlibrary.org/api/books?bibkeys=ISBN:9780134190440&format=json

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
