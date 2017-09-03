package jobtracker

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"scheduler/common/dbsvc"
	"scheduler/common/utils"
	"scheduler/tracker"
)

type Tracker struct {
	name    string
	jid     uint32
	jobhold tracker.JobObject
}

const (
	JTK_PREFIX = "cruiser|jbtrk|jobid|"
)

/**
in go style, always return error when New object.
*/
func NewTracker(jid interface{}) (*Tracker, error) {
	host, _ := os.Hostname()
	id := jid.(uint32)
	return &Tracker{
		name: host,
		jid:  id,
	}, nil
}

/**
a goroutine to publish job to work queue
always running as a goroutine.
scan job list, to check if a job need to publish,
if it is delayed enough?
if it is published by other jobtracker?
other conditions?
then push it in to job queue!
*/
func (this *Tracker) Publish() {
	firstrun := true
	for {
		now := time.Now()
		fmt.Printf("%v: i am ready to put a job into queue\n", time.Now().String())
		//get job Entity each time to be aware of job updating.
		job, gerr := dbsvc.GetEntity(this.jid)
		if gerr != nil {
			break
		}
		this.jobhold = job.(tracker.JobObject)

		if this.jobhold.Nomore == 1 {
			break
		}
		if !firstrun {
			go this.put()
		}
		firstrun = false

		elapse := time.Since(now)
		now = time.Now()
		delay, err := this.calcDelaySeconds(elapse)
		// fmt.Printf("calc delay: %v\n", delay)
		if err != nil {
			fmt.Printf("calc delay err: %v\n", err)
			break
		}
		//tolerate time difference between two machines
		if rerr := this.refreshAliveFlag(delay.(int64) + HBTimeout); rerr != nil {
			//treat as this job tracker is down
			break
		}
		// fmt.Printf("report alive: %v\n", delay.(int64)+HBTimeout)
		elapse = time.Since(now)
		delayDuration, _ := time.ParseDuration(fmt.Sprintf("%vs", delay.(int64)))
		// fmt.Printf("sleep duration: %v %T\n", delayDuration-elapse, delayDuration-elapse)
		time.Sleep(delayDuration - elapse)
	}
}

func (this *Tracker) put() {
	var tasks []interface{}
	json.Unmarshal([]byte(this.jobhold.Body), &tasks)
	var taskSlice = make([]tracker.Task, 0)
	for _, task := range tasks {
		t := task.(map[string]interface{})
		//memory leaks or not?
		T := tracker.Task{
			Worker: t["Worker"].(string),
			Fields: t["Fields"].(map[string]interface{}),
		}
		taskSlice = append(taskSlice, T)
	}
	job := tracker.Job{
		OrderReady: true, //let's set true temporarily
		Tasks:      taskSlice,
	}

	for {
		if err := Put(job); err != nil {
			continue
		}
		break
	}
	fmt.Printf("%v: i put a job into queue\n", time.Now().String())

	if this.jobhold.Jtype == 0 {
		set := map[string]interface{}{"nomore": 1}
		_ = dbsvc.UpdateEntity(this.jid, set)
	}
}

func (this *Tracker) calcDelaySeconds(d time.Duration) (interface{}, error) {
	switch this.jobhold.Jtype {
	case tracker.DELAY_TYPE:
		fallthrough
	case tracker.TICK_TYPE:
		if strings.Index(strings.Trim(this.jobhold.Clock, " "), " ") > 0 {
			return nil, fmt.Errorf("Clock is wrong: %v", this.jobhold.Clock)
		}
		delay, err := strconv.ParseInt(this.jobhold.Clock, 10, 64)
		if err != nil {
			return 0, err
		}
		return delay - int64(d.Seconds()), nil
	case tracker.CLOCK_TYPE:
		return utils.NearestFuture(this.jobhold.Clock)
	default:
		return nil, fmt.Errorf("unknown job type: %v", this.jobhold.Jtype)
	}
}

func (this *Tracker) refreshAliveFlag(life int64) error {
	key := JobTrackerKey(this.jid)
	if _, err := dbsvc.Setex(key, life, this.name); err != nil {
		return err
	}
	return nil
}

/**
request a job from queue,
then return a Job struct,
each time it called by a task tracker.
*/
func Reserve() (interface{}, error) {
	job, err := queue.Pop()
	if err != nil {
		fmt.Println("Reserve error: %v", err)
		return nil, err
	}
	return job, nil
}

func Put(job interface{}) error {
	err := queue.Push(job)
	if err != nil {
		fmt.Println("Put error: %v", err)
		return err
	}
	return nil
}

func JobTrackerKey(jid interface{}) string {
	return JTK_PREFIX + fmt.Sprintf("%v", jid)
}
