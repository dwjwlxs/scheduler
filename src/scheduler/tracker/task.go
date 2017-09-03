package tracker

type Task struct {
	TID    uint32                 //task id
	PreIDs []uint32               //pre tasks
	Worker string                 //the worker's name
	Fields map[string]interface{} //like ingredients
}
