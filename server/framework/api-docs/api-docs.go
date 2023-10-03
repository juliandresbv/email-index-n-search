package apidocs

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	httpswagger "github.com/swaggo/http-swagger"

	customlogger "pkg/custom-logger"
)

var logger = customlogger.NewLogger()

func NewApiDocs(server *chi.Mux, hostPort string) {
	server.Get("/docs/*", httpswagger.Handler(
		httpswagger.URL(fmt.Sprintf("http://%v/docs/doc.json", hostPort)),
	))

	logger.Printf("API docs running on: http://%v/docs/index.html", hostPort)
}
