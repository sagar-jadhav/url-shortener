package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sagar-jadhav/url-shortener/model"
	"github.com/sagar-jadhav/url-shortener/pkg/datastore"
	"github.com/sagar-jadhav/url-shortener/pkg/metrics"
	"github.com/sagar-jadhav/url-shortener/pkg/utils"
)

type Shortener struct {
	Datastore            datastore.Datastore
	ShortURLSize         int
	Domain               string
	CollisionRetryCount  int
	GenerateRandomString func(int) string
	Metrics              metrics.DomainMetrics
}

// ShortenURL generates the short URL and store it into memory
func (s *Shortener) ShortenURL(w http.ResponseWriter, req *http.Request) {
	var err error

	// validate that long URL should be present in the request body
	reqBody := &model.Request{}
	err = json.NewDecoder(req.Body).Decode(reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("error in parsing the request body: %v", err), http.StatusInternalServerError)
		return
	}
	if len(reqBody.LongURL) == 0 {
		http.Error(w, "long URL is required", http.StatusInternalServerError)
		return
	} else if !utils.ValidateURL(reqBody.LongURL) {
		http.Error(w, fmt.Sprintf("long URL %s is invalid", reqBody.LongURL), http.StatusBadRequest)
		return
	}

	var shortURL string
	var longURLExist bool
	if longURLExist, err = s.Datastore.DoesLongURLExist(reqBody.LongURL); err != nil {
		http.Error(w, fmt.Sprintf("error in checking whether the long URL %s is exist in the memory or not: %v", reqBody.LongURL, err), http.StatusInternalServerError)
		return
	}

	if longURLExist { // If long URL already exist then return the old short URL
		if shortURL, err = s.Datastore.GetShortURL(reqBody.LongURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else { // else generate the short URL and then insert it into memory
		shortURLExist := false
		// validate that short URL should not be present in the memory
		for i := 0; i < s.CollisionRetryCount; i++ {
			shortURL = s.Domain + s.GenerateRandomString(s.ShortURLSize)
			if shortURLExist, err = s.Datastore.DoesShortURLExist(shortURL); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if !shortURLExist {
				break
			}
		}
		if shortURLExist {
			http.Error(w, fmt.Sprintf("short URL %s already exist in the memory and all the %d retries also exhausted. So please call the api again.", shortURL, s.CollisionRetryCount), http.StatusInternalServerError)
			return
		}

		// adding domain metrics
		domain := utils.GetDomainName(reqBody.LongURL)
		s.Metrics.Add(domain)

		// inserting the data
		s.Datastore.Insert(reqBody.LongURL, shortURL)
	}
	resp := model.Response{
		LongURL:  reqBody.LongURL,
		ShortURL: shortURL,
	}
	var b []byte
	if b, err = json.Marshal(resp); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing the response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	return
}

// Redirect redirects short URL to the long URL
func (s *Shortener) Redirect(w http.ResponseWriter, req *http.Request) {
	shortURL := s.Domain + chi.URLParam(req, "shortUrl")

	var shortURLExist bool
	var err error
	if shortURLExist, err = s.Datastore.DoesShortURLExist(shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if !shortURLExist {
		http.Error(w, fmt.Sprintf("short URL %s not found", shortURL), http.StatusNotFound)
	} else {
		var longURL string
		if longURL, err = s.Datastore.GetLongURL(shortURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, req, longURL, http.StatusSeeOther)
		}
	}
	return
}

// GetMetrics returns the sorted metrics based on the shortened count per domain
func (s *Shortener) GetMetrics(w http.ResponseWriter, req *http.Request) {
	// sort the metrics
	sortedMetrics := s.Metrics.Sort()

	var b []byte
	var err error
	var result []metrics.Metric

	// get the first 3 domain metrics
	if len(sortedMetrics) > 3 {
		result = sortedMetrics[:3]
	} else {
		result = sortedMetrics
	}

	if b, err = json.Marshal(result); err != nil {
		http.Error(w, fmt.Sprintf("error in parsing the response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	return
}
