package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sagar-jadhav/url-shortener/model"
	"github.com/sagar-jadhav/url-shortener/pkg/datastore"
	"github.com/sagar-jadhav/url-shortener/pkg/utils"
)

type Shortener struct {
	Datastore    datastore.Datastore
	ShortURLSize int
	Domain       string
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
		http.Error(w, "longURL is required", http.StatusInternalServerError)
		return
	}

	var shortURL string
	var exist bool
	if exist, err = s.Datastore.Exist(reqBody.LongURL); err != nil {
		http.Error(w, fmt.Sprintf("error in checking whether the long URL %s is exist in the memory or not: %v", reqBody.LongURL, err), http.StatusInternalServerError)
		return
	}

	if exist { // If long URL already exist then return the old short URL
		if shortURL, err = s.Datastore.Get(reqBody.LongURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else { // else generate the short URL and then insert it into memory
		shortURL = s.Domain + utils.GenerateRandomString(s.ShortURLSize)
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
