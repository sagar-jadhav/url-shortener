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

// DoesLongURLExist returns true if Long URL exists in the Map else return false
func (mds *MemoryDatastore) DoesLongURLExist(longURL string) (bool, error) {
	if mds.Data == nil {
		return false, nil
	}
	_, ok := mds.Data[longURL]
	return ok, nil
}

// GetShortURL returns the Short URL for the given Long URL from the Map
func (mds *MemoryDatastore) GetShortURL(longURL string) (string, error) {
	if exist, _ := mds.DoesLongURLExist(longURL); !exist {
		return "", fmt.Errorf("%s URL not exist in the map", longURL)
	}
	return mds.Data[longURL], nil
}

// GetLongURL returns the Long URL for the given Short URL from the Map
func (mds *MemoryDatastore) GetLongURL(shortURL string) (string, error) {
	if exist, _ := mds.DoesShortURLExist(shortURL); !exist {
		return "", fmt.Errorf("%s URL not exist in the map", shortURL)
	}
	for longURL, url := range mds.Data {
		if url == shortURL {
			return longURL, nil
		}
	}
	return "", nil
}

// DoesShortURLExist returns true if Short URL exists in the Map else return false
func (mds *MemoryDatastore) DoesShortURLExist(shortURL string) (bool, error) {
	if mds.Data == nil {
		return false, nil
	}
	for _, url := range mds.Data {
		if url == shortURL {
			return true, nil
		}
	}
	return false, nil
}
