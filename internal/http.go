package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

func NewServer(tp TemplateRenderer, accessLogging bool) http.Handler {
	mux := http.NewServeMux()

	ss := &CommandHandler{
		Mutex: &sync.Mutex{},
		Shortcuts: map[string]string{
			"*":    DefaultSearchProvider,
			"help": "/",
		},
	}

	mux.HandleFunc("/", tp.RenderHandler("index.html", ss, nil))
	mux.HandleFunc("/index.html", tp.RenderHandler("index.html", ss, nil))
	mux.HandleFunc("/search_plugin.xml", tp.RenderHandler("search_plugin.xml", ss, http.Header{
		"Content-Type": []string{"application/opensearchdescription+xml"},
	}))

	handlerFunc := searchHandler(ss, accessLogging)
	mux.HandleFunc("/search", handlerFunc)
	mux.HandleFunc("/search_to_home", handlerFunc)

	return mux
}

func searchHandler(scs *CommandHandler, logAccess bool) http.HandlerFunc {
	accessLogger := log.New(os.Stdout, "[access] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	return func(res http.ResponseWriter, req *http.Request) {
		v := req.URL.Query()
		commandString := v.Get("q")
		action, err := scs.Handle(commandString)
		if err != nil {
			accessLogger.Printf("'%s' -> got error: '%s'", commandString, err.Error())
			http.Redirect(res, req, fmt.Sprintf("https://duckduckgo.com?q=%s", commandString), http.StatusFound)
			return
		}

		accessLogger.Printf("'%s' %s -> 302 %s", commandString, action.Action, action.Location)
		http.Redirect(res, req, action.Location, http.StatusFound)
	}
}
