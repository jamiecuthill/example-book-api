package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	//
	// {"type": "foo", "data": {"name": "foo"}}
	//

	r.Methods(http.MethodPost).
		Path("/req").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Type string
				Data json.RawMessage
			}
			_ = json.NewDecoder(r.Body).Decode(&body)

			// var subtype struct{}
			// json.Unmarshal(body.Data, &subtype)

			// if body["hello"] == "" {
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	return
			// }

			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, "msg is %+v", body)
		})

	server := &http.Server{
		Addr:    ":8080",
		Handler: withAuth(r),
	}

	go func() {
		fmt.Println("Running http server")
		err := server.ListenAndServe()
		fmt.Println(err)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	inter := <-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	fmt.Printf("Finished %s", inter)
}

func withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
