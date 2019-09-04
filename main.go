package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
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

		name, err := repo.GetBookName(isbn)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %+v\n", err)
			return
		}

		fmt.Fprintf(w, "Your book is: %s\n", name)
	}
}
