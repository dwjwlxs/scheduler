package mailer

import (
	"fmt"
)

/**
params:map[string]interface{}
address,
content,
*/
type Mailer struct {
	Fields map[string]interface{}
}

func NewMailer(fields map[string]interface{}) (*Mailer, error) {
	return &Mailer{
		Fields: fields,
	}, nil
}

func (this *Mailer) Init(fields map[string]interface{}) error {
	this.Fields = fields
	return nil
}

func (this *Mailer) Execute() (interface{}, error) {
	ret := fmt.Sprintf("Send message: %#v to people: %#v", this.Fields["content"], this.Fields["address"])
	fmt.Println(ret)
	return ret, nil
}
