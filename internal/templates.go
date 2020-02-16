package internal

import "net/http"

type TemplateRenderer struct {
	http.FileSystem
}

func (tp TemplateRenderer) Handle(templateFile string, model interface{}) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	})
}
