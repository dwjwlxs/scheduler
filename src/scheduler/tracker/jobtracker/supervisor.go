package jobtracker

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"scheduler/common/dbsvc"
	"scheduler/tracker"
)

type Supervisor struct {
	name string
}

const (
	SCAN_INTERVAL      = 25 //interval between each table scanning, increase it when num of nodes increased.
	MAX_JOBTRACKER     = 256
	SUPER_LOCK         = "cruiser|supvsr_lock"
	SUPER_LOCK_TIMEOUT = 30
	HBTimeout          = 20 //timeout for heart beat
	MaxNewTrackers     = 1
)

var (
	wg          sync.WaitGroup
	JobTrackers = make(map[string]interface{})
	muJT        sync.Mutex
	LimitG      = make(chan int, MAX_JOBTRACKER) //limit number of job trackers in a node.
)

func NewSupervisor() (*Supervisor, error) {
	host, _ := os.Hostname()
	return &Supervisor{
		name: host,
	}, nil
}

func (this *Supervisor) Run() {
	wg.Add(2)
	go this.watchLocalG()
	go this.discoverNewJob()
	wg.Wait()
}

/**
a goroutine try to restart a job tracker instantly when a job tracker is down.
the job may be tracked in other node also.
*/
func (this *Supervisor) watchLocalG() {
	defer wg.Done()
	for {
		time.Sleep(time.Hour)
		//process LimitG
		var count int
		JT := make(map[string]interface{})
		for jobid, _ := range JobTrackers {
			//
			jtkey := JobTrackerKey(jobid)
			if val, err := dbsvc.Get(jtkey); err == nil && val == this.name {
				count += 1
				JT[jobid] = JobTrackers[jobid]
				fmt.Printf("%v: Supervisor@%v has check a running job tracker(%v) spawned by itself\n", time.Now().String(), this.name, jobid)
			}
		}
		muJT.Lock()
		JobTrackers = JT
		muJT.Unlock()
		spare := MAX_JOBTRACKER - count
		for i := 0; i < spare; i++ {
			<-LimitG
		}
		fmt.Printf("%v: Supervisor@%v has updated the count(%v) of running job trackers\n", time.Now().String(), this.name, count)
	}
}

/**
a goroutine to start a new job tracker for a job if needed.
when a job is not tracked, a new job tracker for it needs to be built.
*/
func (this *Supervisor) discoverNewJob() {
	defer wg.Done()
	for {
		fmt.Printf("%v: Supervisor@%v is trying to discover new jobs\n", time.Now().String(), this.name)
		l, _ := queue.Length()
		fmt.Printf("%v: Supervisor@%v checked the size of work queue: %v\n", time.Now().String(), this.name, l)
		LimitG <- 1 //if reach the MAX_JOBTRACKER blocked here.

		timestamp := fmt.Sprintf("%v", (time.Now()).Unix())
		locked, sleep, err := lock(timestamp)
		if !locked {
			if err != nil {
				fmt.Printf("%v: Supervisor@%v is failed to lock, ERROR: %v\n", time.Now().String(), this.name, err)
			}
			if sleep {
				time.Sleep(SCAN_INTERVAL * time.Second)
			}
			continue
		} else {
			fmt.Printf("%v: Supervisor@%v is locked, key: %v\n", time.Now().String(), this.name, timestamp)
		}

		//in the code below, we should delete the lock before the loop ends
		//scan job list
		jobs, err := dbsvc.ListJob()
		if err != nil {
			fmt.Printf("%v: Supervisor@%v encountered an error when query jobs, will unlock and sleep for some time. Error: %v\n", time.Now().String(), this.name, err)
			if err = unlock(timestamp); err != nil {
				fmt.Printf("%v: Supervisor@%v is failed to unlock, ERROR: %v\n", time.Now().String(), this.name, err)
			} else {
				fmt.Printf("%v: Supervisor@%v is unlocked, key: %v\n", time.Now().String(), this.name, timestamp)
			}
			time.Sleep(SCAN_INTERVAL * time.Second)
			continue
		}
		jobobjects := jobs.([]tracker.JobObject)
		if len(jobobjects) <= 0 {
			fmt.Printf("%v: Supervisor@%v found no job\n", time.Now().String(), this.name)
		}
		count := 0
		for _, job := range jobobjects {
			//check if there is a job tracker for each job.
			//by check a expired string in redis
			//yes? continue annother loop
			if count >= MaxNewTrackers {
				break
			}
			jtkey := JobTrackerKey(job.Jid)
			if ok, err := dbsvc.Exists(jtkey); err != nil || ok == 1 {
				node, _ := dbsvc.Get(jtkey)
				fmt.Printf("%v: Supervisor@%v has found a existing job tracker: %v on node(%v), will continue\n", time.Now().String(), this.name, job.Jid, node)
				continue
			}
			time.Sleep(time.Second * HBTimeout)
			//to make sure the job tracker does not exists
			if ok, err := dbsvc.Exists(jtkey); err != nil || ok == 1 {
				node, _ := dbsvc.Get(jtkey)
				fmt.Printf("%v: Supervisor@%v has found a existing job tracker: %v on node(%v), will continue\n", time.Now().String(), this.name, job.Jid, node)
				continue
			}
			//not exists? New a job tracker for the job
			tracker, _ := NewTracker(job.Jid)
			count++
			muJT.Lock()
			JobTrackers[fmt.Sprintf("%v", job.Jid)] = job
			muJT.Unlock()
			go tracker.Publish()
			fmt.Printf("%v: Supervisor@%v has built a job tracker: %v\n", time.Now().String(), this.name, job.Jid)
		}
		if err = unlock(timestamp); err != nil {
			fmt.Printf("%v: Supervisor@%v is failed to unlock, ERROR: %v\n", time.Now().String(), this.name, err)
		} else {
			fmt.Printf("%v: Supervisor@%v is unlocked: %v\n", time.Now().String(), this.name, timestamp)
		}
		time.Sleep(SCAN_INTERVAL * time.Second)
	}
}

func unlock(key string) error {
	lock_val, err := dbsvc.Get(SUPER_LOCK)
	if err != nil {
		return err
	}
	//avoid to delete other's lock, caused by timeout
	if lock_val != key {
		err = fmt.Errorf("not my lock")
		return err
	}
	if _, err = dbsvc.Del(SUPER_LOCK); err != nil {
		return err
	}
	return nil
}

func lock(key string) (locked, sleep bool, err error) {
	//distributed Lock
	ok, err := dbsvc.Setnx(SUPER_LOCK, key)
	if err != nil {
		// time.Sleep(SCAN_INTERVAL * time.Second)
		// continue
		return false, true, err
	}
	if ok != 1 {
		lock_val, err := dbsvc.Get(SUPER_LOCK)
		if err != nil {
			// time.Sleep(SCAN_INTERVAL * time.Second)
			// continue
			return false, true, err
		}
		oldtime, _ := strconv.ParseInt(lock_val, 10, 64)
		newtime := (time.Now()).Unix()
		if newtime-oldtime < SUPER_LOCK_TIMEOUT {
			// time.Sleep(SCAN_INTERVAL * time.Second)
			// continue
			err = fmt.Errorf("not timeouted lock")
			return false, true, err
		}

		oldtime_2_str, err := dbsvc.GetSet(SUPER_LOCK, fmt.Sprintf("%v", newtime))
		if err != nil {
			// time.Sleep(SCAN_INTERVAL * time.Second)
			// continue
			return false, true, err
		}
		oldtime_2, _ := strconv.ParseInt(oldtime_2_str, 10, 64)
		if oldtime_2 == oldtime {
			//then this instance has the right to delete the supervisor lock
			_, err = dbsvc.Del(SUPER_LOCK)
			if err == nil {
				fmt.Printf("@Supervisor has deleted the timeouted lock: %v\n", oldtime)
			}
			// continue
			return false, false, err
		} else {
			//no right, wash wash and go to sleep
			// time.Sleep(SCAN_INTERVAL * time.Second)
			// continue
			return false, true, err
		}
	}
	return true, false, nil
}
