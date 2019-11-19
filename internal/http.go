package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Shortcuts map[string]string

func NewServer(tp TemplateRenderer, shortcutStore *ShortcutStore, accessLogging bool) http.Handler {
	mux := http.NewServeMux()

	shorts, _ := shortcutStore.LoadShortcuts(nil)
	ss := &CommandHandler{
		Mutex:     &sync.Mutex{},
		Shortcuts: shorts,
	}

	mux.HandleFunc("/", tp.RenderHandler("index.html", ss, nil))
	mux.HandleFunc("/index.html", tp.RenderHandler("index.html", ss, nil))
	mux.HandleFunc("/search_plugin.xml", tp.RenderHandler("search_plugin.xml", ss, http.Header{
		"Content-Type": []string{"application/opensearchdescription+xml"},
	}))

	handlerFunc := searchHandler(ss, shortcutStore, accessLogging)
	mux.HandleFunc("/search", handlerFunc)
	mux.HandleFunc("/search_to_home", handlerFunc)

	return mux
}

func searchHandler(scs *CommandHandler, shortcutStore *ShortcutStore, logAccess bool) http.HandlerFunc {
	accessLogger := log.New(os.Stdout, "[access] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	return func(res http.ResponseWriter, req *http.Request) {
		v := req.URL.Query()
		commandString := v.Get("q")

		action, err := scs.Handle(commandString)
		if err != nil {
			if logAccess {
				accessLogger.Printf("'%s' -> got error: '%s'", commandString, err.Error())
			}
			http.Redirect(res, req, fmt.Sprintf("https://duckduckgo.com?q=%s", commandString), http.StatusFound)

			return
		}

		accessLogger.Printf("'%s' %s -> 302 %s", commandString, action.Action, action.Location)

		if action.Action != "lookup" {
			if err := shortcutStore.SaveShortcuts(scs.Shortcuts, nil); err != nil {
				if logAccess {
					accessLogger.Printf("'%s' %s -> could not save shortcuts database file: %v", commandString, action.Action, err)
				}
				http.Error(res, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		http.Redirect(res, req, action.Location, http.StatusFound)
	}
}
