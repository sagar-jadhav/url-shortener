package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sagar-jadhav/url-shortener/model"
	"github.com/sagar-jadhav/url-shortener/pkg/datastore"
	"github.com/sagar-jadhav/url-shortener/pkg/utils"
)

var (
	longURL  = "http://github.com/sagar-jadhav"
	shortURL = "http://localhost:3000/abcde"
)

func Test_ShortenURL(t *testing.T) {
	tests := []struct {
		name      string
		reqBody   *model.Request
		shortener Shortener
		status    int
	}{
		{
			name: "long URL in request body is empty",
			reqBody: &model.Request{
				LongURL: "",
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{},
				},
				ShortURLSize:         5,
				Domain:               "http://localhost:3000/",
				CollisionRetryCount:  5,
				GenerateRandomString: utils.GenerateRandomString,
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "long URL exist in the memory",
			reqBody: &model.Request{
				LongURL: longURL,
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{
						longURL: shortURL,
					},
				},
				ShortURLSize:         5,
				Domain:               "http://localhost:3000/",
				CollisionRetryCount:  5,
				GenerateRandomString: utils.GenerateRandomString,
			},
			status: http.StatusOK,
		},
		{
			name: "long URL & short URL not exist in the memory",
			reqBody: &model.Request{
				LongURL: longURL,
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{},
				},
				ShortURLSize:         5,
				Domain:               "http://localhost:3000/",
				CollisionRetryCount:  5,
				GenerateRandomString: utils.GenerateRandomString,
			},
			status: http.StatusOK,
		},
		{
			name: "short URL exist in the memory",
			reqBody: &model.Request{
				LongURL: longURL,
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{
						longURL + "/url-shortener": shortURL,
					},
				},
				ShortURLSize:        5,
				Domain:              "http://localhost:3000/",
				CollisionRetryCount: 5,
				GenerateRandomString: func(i int) string {
					return "abcde"
				},
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// starting the server
			r := chi.NewRouter()
			r.Use(middleware.Logger)
			r.Post("/", test.shortener.ShortenURL)

			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatalf("error not expecting but got: %v", err)
			}
			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.status {
				t.Fatalf("ShortenURL() => expected %d but got %d", test.status, rr.Code)
			}
		})
	}
}

func Test_Redirect(t *testing.T) {
	tests := []struct {
		name      string
		shortener Shortener
		status    int
		shortURL  string
	}{
		{
			name: "short URL not found",
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{},
				},
				ShortURLSize:         5,
				Domain:               "http://localhost:3000/",
				CollisionRetryCount:  5,
				GenerateRandomString: utils.GenerateRandomString,
			},
			shortURL: "abcde",
			status:   http.StatusNotFound,
		},
		{
			name: "short URL present",
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{
						longURL: shortURL,
					},
				},
				ShortURLSize:         5,
				Domain:               "http://localhost:3000/",
				CollisionRetryCount:  5,
				GenerateRandomString: utils.GenerateRandomString,
			},
			shortURL: "abcde",
			status:   http.StatusSeeOther,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// starting the server
			r := chi.NewRouter()
			r.Use(middleware.Logger)
			r.Get("/{shortUrl}", test.shortener.Redirect)

			req, _ := http.NewRequest(http.MethodGet, "/"+test.shortURL, nil)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != test.status {
				t.Fatalf("Redirect() => expected %d but got %d", test.status, rr.Code)
			}
		})
	}
}
