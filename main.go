package main

import (

	//"fmt"

	//"fmt"
	//"log"

	"net/http"
	"os"

	"web-server/controller"
	"web-server/logger"
	"web-server/service"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.uber.org/zap"
)

var (
	logg *logger.Logger
	db   *service.Database
)

func init() {
	logg = logger.Init()

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logg.Error("Error loading .env file", zap.Error(err))
		return
	}

	var err error
	db, err = service.InitDatabase(logg)
	if err != nil {
		logg.Error("Error connecting to database", zap.Error(err))
		os.Exit(1)
	}
}

func main() {

	defer db.Disconnect()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", controller.HandleGetUsers(db))
		r.Post("/", controller.HandleCreateUser(db))
		r.Put("/", controller.HandleUpdateUser(db))
		r.Delete("/", controller.HandleDeleteUser(db))
	})

	port := os.Getenv("PORT")
	if port == "" {
		logg.Error("Set the 'PORT' environment variable.")
		//return nil, errors.New("no PORT number provided")
	}

	logg.Info("Server started on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		logg.Error("Failed to start server", zap.Error(err))

	}
}
