package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	customlogger "pkg/custom-logger"
	searchdbclient "pkg/zinc-search-db/db-client/search-db-client"
	dbclienttypes "pkg/zinc-search-db/db-client/types"

	emailscontroller "server/adapter/controllers/emails"
	usecases "server/application/use-cases"
	apidocs "server/framework/api-docs"
	"server/framework/config"

	_ "server/docs"
)

var logger = customlogger.NewLogger()

func init() {
	config.LoadEnvVars()
}

// @title					Emails API
// @version				1.0.0
// @description		This is the API doc for Emails API.
// @host					0.0.0.0:7070
// @BasePath			/
// @contact.name	Julian Bermudez Valderrama
// @contact.email	julian.berval@gmail.com
func main() {
	server := chi.NewRouter()

	server.Use(middleware.Logger)
	server.Use(render.SetContentType(render.ContentTypeJSON))
	server.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*", "http://localhost:5173", "http://0.0.0.0:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	searchDbClient := searchdbclient.NewSearchDbClient(dbclienttypes.DbClientConfig{})
	emailsUseCase := usecases.NewEmailsUseCase(searchDbClient)
	emailscontroller.NewEmailsController(server, emailsUseCase)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "7070"
	}

	apidocs.NewApiDocs(server)

	logger.Println("Server is running on port " + port)

	err := http.ListenAndServe(":"+port, server)

	if err != nil {
		logger.Println(err)

		panic(err)
	}
}
