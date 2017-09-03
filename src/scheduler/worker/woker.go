package worker

import (
	"scheduler/worker/crawler"
	"scheduler/worker/mailer"
)

type Worker interface {
	Init(fields map[string]interface{}) error
	Execute() (interface{}, error)
}

func Instance(name string, fields map[string]interface{}) Worker {
	var w Worker
	switch name {
	case "mailer":
		w, _ = mailer.NewMailer(fields)
	case "crawler":
		w, _ = crawler.NewCrawler(fields)
	default:
		return w
	}
	return w
}
