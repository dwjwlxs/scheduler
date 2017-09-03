package dbsvc

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"scheduler/common/mysql"
	"scheduler/tracker"
)

const RETRY_TIMES = 3

var dbInfo = map[string]interface{}{
	"host":     "your ip",
	"port":     "your port",
	"user":     "your username",
	"password": "your pass",
	"dbname":   "your database name", //table struct refer to init.sql
}
var m *mysql.Mysql

func init() {
	var err error
	if m, err = mysql.NewMysql(dbInfo, false, 2*time.Second); err != nil {
		fmt.Println("error occured when New Mysql: ", err)
		os.Exit(1)
	}
}

func ListEntity(where, order, limit interface{}) (interface{}, error) {
	sql := "SELECT id,type,lastrun,nomore,clock,body,status FROM job"
	if where.(string) != "" {
		sql += " WHERE " + where.(string)
	}
	sql += " " + order.(string) + " " + limit.(string)
	rows, err := m.DB.Query(sql)
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			rows, err = m.DB.Query(sql)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]tracker.JobObject, 0)
	for rows.Next() {
		job := tracker.JobObject{}
		if err := rows.Scan(&job.Jid, &job.Jtype, &job.Lastrun, &job.Nomore, &job.Clock, &job.Body, &job.Status); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}

func ListJob() (interface{}, error) {
	return ListEntity("", "", "")
}

func UpdateEntity(id, set interface{}) error {
	if set == nil {
		return errors.New("set must not be nil")
	}
	sets := set.(map[string]interface{})
	var setstr string
	for key, val := range sets {
		setstr += key + "='" + fmt.Sprintf("%v", val) + "',"
	}
	sql := "UPDATE job SET " + strings.Trim(setstr, ", ") + " WHERE id=" + fmt.Sprintf("%v", id)
	res, err := m.DB.Exec(sql)
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			res, err = m.DB.Exec(sql)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("affected rows: %v", rows)
	}
	return nil
}

func GetEntity(id interface{}) (interface{}, error) {
	job := tracker.JobObject{}
	sql := "SELECT id,type,lastrun,nomore,clock,body,status FROM job WHERE id=" + fmt.Sprintf("%v", id)
	err := m.DB.QueryRow(sql).Scan(&job.Jid, &job.Jtype, &job.Lastrun, &job.Nomore, &job.Clock, &job.Body, &job.Status)
	if err != nil {
		for i := 0; i < RETRY_TIMES; i++ {
			err = m.DB.QueryRow(sql).Scan(&job.Jid, &job.Jtype, &job.Lastrun, &job.Nomore, &job.Clock, &job.Body, &job.Status)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}
