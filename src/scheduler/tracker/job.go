package tracker

type Job struct {
	OrderReady bool   //if tasks in right order
	Tasks      []Task //tasks to be done in some right order.
}

type JobObject struct {
	Jid     uint32
	Jtype   uint8  //0 1 2
	Lastrun int64  //timestamp of last run
	Nomore  uint8  //one more time?
	Clock   string //depend on jtype
	Body    string //json
	Status  uint8  //complete? status
}

//job type
const (
	DELAY_TYPE = iota
	TICK_TYPE  = iota
	CLOCK_TYPE = iota
)
