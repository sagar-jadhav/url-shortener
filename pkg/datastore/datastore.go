package datastore

type Datastore interface {
	Get(string) (string, error)
	Insert(string, string) error
	DoesLongURLExist(string) (bool, error)
	DoesShortURLExist(string) (bool, error)
}
