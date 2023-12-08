package datastore

type Datastore interface {
	GetShortURL(string) (string, error)
	GetLongURL(string) (string, error)
	Insert(string, string) error
	DoesLongURLExist(string) (bool, error)
	DoesShortURLExist(string) (bool, error)
}
