package orchestrator

type DataStore interface {
	CheckIfAPExists(AP) (bool, error)
	InsertAP(AP) (int, string)
}
