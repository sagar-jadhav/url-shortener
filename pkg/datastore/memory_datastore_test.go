package datastore

import "testing"

func Test_Insert(t *testing.T) {
	mds := MemoryDatastore{}
	longURL := "https://localhost:3000/longurl"
	shortURL := "https://localhost:3000/shorturl"
	mds.Insert(longURL, shortURL)

	val, ok := mds.Data[longURL]
	if !ok {
		t.Fatal("Insert() => expected data to be inserted into the map but data not present")
	}
	if val != shortURL {
		t.Fatalf("Insert() => expected %s got %s", shortURL, val)
	}
}

func Test_DoesLongURLExist(t *testing.T) {
	mds := MemoryDatastore{}
	longURL := "https://localhost:3000/longurl"
	shortURL := "https://localhost:3000/shorturl"
	mds.Insert(longURL, shortURL)

	exist, err := mds.DoesLongURLExist(longURL)
	if !exist {
		t.Fatal("DoesLongURLExist() => expected long URL to be present in the map but its not present")
	}
	if err != nil {
		t.Fatalf("DoesLongURLExist() => error not expected but got %v", err)
	}
}

func Test_GetShortURL(t *testing.T) {
	mds := MemoryDatastore{}
	longURL1 := "https://localhost:3000/longurl1"
	shortURL1 := "https://localhost:3000/shorturl1"
	longURL2 := "https://localhost:3000/longurl2"
	mds.Insert(longURL1, shortURL1)

	_, err := mds.GetShortURL(longURL2)
	if err == nil {
		t.Fatal("Get() => expected error to be not nil")
	}

	var val string
	val, err = mds.GetShortURL(longURL1)
	if val != shortURL1 {
		t.Fatalf("Get() => expected %s got %s", shortURL1, val)
	}
}

func Test_DoesShortURLExist(t *testing.T) {
	mds := MemoryDatastore{}
	longURL := "https://localhost:3000/longurl"
	shortURL := "https://localhost:3000/shorturl"
	mds.Insert(longURL, shortURL)

	exist, err := mds.DoesShortURLExist(shortURL)
	if !exist {
		t.Fatal("DoesShortURLExist() => expected short URL to be present in the map but its not present")
	}
	if err != nil {
		t.Fatalf("DoesShortURLExist() => error not expected but got %v", err)
	}
}

func Test_GetLongURL(t *testing.T) {
	mds := MemoryDatastore{}
	longURL1 := "https://localhost:3000/longurl1"
	shortURL1 := "https://localhost:3000/shorturl1"
	shortURL2 := "https://localhost:3000/shorturl2"
	mds.Insert(longURL1, shortURL1)

	_, err := mds.GetLongURL(shortURL2)
	if err == nil {
		t.Fatal("Get() => expected error to be not nil")
	}

	var val string
	val, err = mds.GetLongURL(shortURL1)
	if val != longURL1 {
		t.Fatalf("Get() => expected %s got %s", longURL1, val)
	}
}
