package datastore

import "fmt"

type MemoryDatastore struct {
	Data map[string]string
}

// Inserts the Long & Short URL in the map
func (mds *MemoryDatastore) Insert(longURL string, shortURL string) error {
	if mds.Data == nil {
		mds.Data = make(map[string]string)
	}
	mds.Data[longURL] = shortURL
	return nil
}

// Exist checks returns true if Long URL exists in the Map else return false
func (mds *MemoryDatastore) Exist(longURL string) (bool, error) {
	if mds.Data == nil {
		return false, nil
	}
	_, ok := mds.Data[longURL]
	return ok, nil
}

// Get returns the Short URL for the given Long URL from the Map
func (mds *MemoryDatastore) Get(longURL string) (string, error) {
	if exist, _ := mds.Exist(longURL); !exist {
		return "", fmt.Errorf("%s URL not exist in the map", longURL)
	}
	return mds.Data[longURL], nil
}
