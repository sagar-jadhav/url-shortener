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
	"github.com/sagar-jadhav/url-shortener/pkg/utils"
)

const (
	DEFAULT_SHORT_URL_SIZE = 5
	COLLISION_RETRY_COUNT  = 5
)

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
	collisionRetryCountStr := os.Getenv("COLLISION_RETRY_COUNT")

	var shortURLSize int
	if shortURLSize, err = strconv.Atoi(shortURLSizeStr); err != nil {
		log.Printf("error converting short URL size %s. So setting it to %d", shortURLSizeStr, DEFAULT_SHORT_URL_SIZE)
	}

	var collisionRetryCount int
	if collisionRetryCount, err = strconv.Atoi(collisionRetryCountStr); err != nil {
		log.Printf("error converting collision retry count %s. So setting it to %d", collisionRetryCountStr, COLLISION_RETRY_COUNT)
	}

	// Initialising the shortener service
	s := v1.Shortener{
		Datastore:            &datastore.MemoryDatastore{},
		ShortURLSize:         shortURLSize,
		CollisionRetryCount:  collisionRetryCount,
		Domain:               scheme + "://" + appDomain + ":" + port + "/",
		GenerateRandomString: utils.GenerateRandomString,
	}

	// set up the server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", s.ShortenURL)
	r.Get("/{shortUrl}", s.Redirect)

	// starting the server
	log.Printf("server is starting at %s:%s", appDomain, port)
	log.Fatal(http.ListenAndServe(appDomain+":"+port, r))
}
