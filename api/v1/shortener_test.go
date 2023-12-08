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
)

var (
	longURL  = "https://github.com/sagar-jadhav"
	shortURL = "https://localhost:3000/abcde"
)

func Test_ShortenURL(t *testing.T) {
	tests := []struct {
		name      string
		reqBody   *model.Request
		shortener Shortener
		status    int
	}{
		{
			name: "longURL in request body is empty",
			reqBody: &model.Request{
				LongURL: "",
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{},
				},
				ShortURLSize: 5,
				Domain:       "http://localhost:3000/",
			},
			status: http.StatusInternalServerError,
		},
		{
			name: "longURL exist in the memory",
			reqBody: &model.Request{
				LongURL: longURL,
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{
						longURL: shortURL,
					},
				},
				ShortURLSize: 5,
				Domain:       "http://localhost:3000/",
			},
			status: http.StatusOK,
		},
		{
			name: "longURL not exist in the memory",
			reqBody: &model.Request{
				LongURL: longURL,
			},
			shortener: Shortener{
				Datastore: &datastore.MemoryDatastore{
					Data: map[string]string{},
				},
				ShortURLSize: 5,
				Domain:       "http://localhost:3000/",
			},
			status: http.StatusOK,
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
