package main

import (
	"SongLibrary/config"
	"SongLibrary/pkg/song"
	"SongLibrary/pkg/storage"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "SongLibrary/docs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title SongLibrary Swagger API
// @version 1.0
// @description API for song library
func main() {

	cfg, err := config.ReadConfig("app", "./config")
	if err != nil {
		log.Println("reading config file failed")
		return
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("open sql connect failed")
		return
	}
	defer db.Close()

	migration, err := migrate.New("file://migration", dsn)
	if err != nil {
		log.Println("error1 here:", err)
		return
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("error2 here:", err)
		return
	}

	err = db.Ping()
	if err != nil {
		log.Println("ping failed")
		return
	}

	log.Println("db successfully connected")

	songHandler := &song.SongHandler{
		Storage:     storage.NewStorage(db),
		ExternalAPI: cfg.ExternalAPI,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/", swaggerHandler)
	mux.HandleFunc("GET /api/songs", songHandler.GetAll)
	mux.HandleFunc("GET /api/songs/{id}", songHandler.Get)
	mux.HandleFunc("PUT /api/songs", songHandler.New)
	mux.HandleFunc("POST /api/songs", songHandler.Update)
	mux.HandleFunc("DELETE /api/songs", songHandler.Delete)

	log.Println("start server at", cfg.HTTPPort)
	http.ListenAndServe(":"+cfg.HTTPPort, mux)

}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	httpSwagger.WrapHandler(w, r)
}
