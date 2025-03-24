package data

import "github.com/go-kratos/kratos/v2/log"

type exampleRepo struct {
	data *Data
	log  *log.Helper
}

func NewExampleRepo(data *Data, logger log.Logger) biz.ExampleRepo {
	return &ExampleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
