package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"html/template"
)

type TemplateRenderer struct {
	FS http.FileSystem
}

type TemplateVariables struct {
	Host    string
	Entries map[string]string
}

func ToTemplateVars(req *http.Request, shortcuts map[string]string) TemplateVariables {
	scheme := req.URL.Scheme
	host := req.URL.Host

	if scheme == "" {
		scheme = "http"
	}

	if host == "" {
		if host = req.Host; host == "" {
			host = req.Header.Get("X-Forwarded-For")
		}
	}

	fqdn := fmt.Sprintf("%s://%s", scheme, host)

	return TemplateVariables{
		Host:    fqdn,
		Entries: shortcuts,
	}
}

func (tr TemplateRenderer) Render(name string, templVars TemplateVariables, w io.Writer) error {
	templateFileName := fmt.Sprintf("%s.tpl", name)

	f, err := tr.FS.Open(templateFileName)
	if err != nil {
		return err
	}

	defer f.Close()

	templateBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	templateStr := string(templateBytes)

	templateObj, parseTemplateErr := template.New(templateFileName).Parse(templateStr)
	if parseTemplateErr != nil {
		return fmt.Errorf("could not parse template for '%s': %w", templateFileName, parseTemplateErr)
	}

	log.Printf("%+v", templVars)

	return templateObj.Execute(w, templVars)
}

func (tr TemplateRenderer) RenderHandler(filename string, ss *CommandHandler, headers http.Header) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		templVar := ToTemplateVars(req, ss.Shortcuts)
		h := res.Header()

		for k, v := range headers {
			h[k] = v
		}

		h.Set("Cache-Control", "max-age 0; no-cache; private")

		if err := tr.Render(filename, templVar, res); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
