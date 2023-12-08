package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	v1 "github.com/sagar-jadhav/url-shortener/api/v1"
	"github.com/sagar-jadhav/url-shortener/pkg/datastore"
)

const DEFAULT_SHORT_URL_SIZE = 5

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	appDomain := os.Getenv("APP_DOMAIN")
	port := os.Getenv("PORT")
	shortURLSizeStr := os.Getenv("SHORT_URL_SIZE")
	scheme := os.Getenv("SCHEME")
	var shortURLSize int
	if shortURLSize, err = strconv.Atoi(shortURLSizeStr); err != nil {
		log.Printf("error converting short URL size %s. So setting it to %d", shortURLSizeStr, DEFAULT_SHORT_URL_SIZE)
	}

	// Initialising the shortener service
	s := v1.Shortener{
		Datastore:    &datastore.MemoryDatastore{},
		ShortURLSize: shortURLSize,
		Domain:       scheme + "://" + appDomain + ":" + port + "/",
	}

	// set up the server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", s.ShortenURL)

	// starting the server
	log.Printf("server is starting at %s:%s", appDomain, port)
	log.Fatal(http.ListenAndServe(appDomain+":"+port, r))
}
