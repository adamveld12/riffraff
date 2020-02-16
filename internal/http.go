package internal

import (
	"net/http"
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

	return AccessLoggerMiddleware(mux)
}
