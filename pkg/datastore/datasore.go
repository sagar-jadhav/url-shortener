package datastore

type Datastore interface {
	Get(string) (string, error)
	Insert(string, string) error
	Exist(string) (bool, error)
}
