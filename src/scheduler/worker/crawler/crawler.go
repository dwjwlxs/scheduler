package crawler

type Crawler struct {
	Fields map[string]interface{}
}

func NewCrawler(fields map[string]interface{}) (*Crawler, error) {
	return &Crawler{
		Fields: fields,
	}, nil
}

func (this *Crawler) Init(fields map[string]interface{}) error {
	this.Fields = fields
	return nil
}

func (this *Crawler) Execute() (interface{}, error) {
	return nil, nil
}
