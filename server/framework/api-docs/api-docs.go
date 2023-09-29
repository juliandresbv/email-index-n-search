package apidocs

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	httpswagger "github.com/swaggo/http-swagger"
)

func NewApiDocs(server *chi.Mux) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "7070"
	}

	server.Get("/docs/*", httpswagger.Handler(
		httpswagger.URL(fmt.Sprintf("http://%v:%v/docs/doc.json", host, port)),
	))
}
