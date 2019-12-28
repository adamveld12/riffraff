package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
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

	return accessLoggerMiddleware(mux)
}

func accessLoggerMiddleware(middleware http.Handler) http.Handler {
	accessLogger := log.New(os.Stdout, "[access] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		requestId := uuid.NewV4()
		setupHeaders(res.Header(), requestId.String())
		accessLogger.Printf("[%s] %s %s '%s'", requestId, getRemoteIP(req), req.Method, req.URL.String())
		middleware.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), "id", requestId.String())))
	})
}

func searchHandler(scs *CommandHandler, shortcutStore *ShortcutStore, logAccess bool) http.HandlerFunc {
	searchLogger := log.New(os.Stdout, "[search] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	return func(res http.ResponseWriter, req *http.Request) {
		requestId := req.Context().Value("id")
		v := req.URL.Query()
		commandString := v.Get("q")

		action, err := scs.Handle(commandString)
		if err != nil {
			if logAccess {
				searchLogger.Printf("[%s] '%s' -> got error: '%s'", requestId, commandString, err.Error())
			}

			http.Redirect(res, req, fmt.Sprintf("https://duckduckgo.com?q=%s", commandString), http.StatusFound)
			return
		}

		searchLogger.Printf("[%s] '%s' %s -> 302 %s", requestId, commandString, action.Action, action.Location)
		if action.Action != "lookup" {
			if err := shortcutStore.SaveShortcuts(scs.Shortcuts, nil); err != nil {
				if logAccess {
					searchLogger.Printf("[%s] '%s' %s -> could not save shortcuts database file: %v", requestId, commandString, action.Action, err)
				}

				http.Error(res, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		http.Redirect(res, req, action.Location, http.StatusFound)
	}
}

func setupHeaders(headers http.Header, requestId string) {
	headers.Set("X-RiffRaff-Request-Id", requestId)
	headers.Set("Server", "Riff Raff")
	headers.Set("Cache-Control", fmt.Sprintf("public, max-age=%s", time.Hour/time.Second))
}

func getRemoteIP(req *http.Request) string {
	ip := req.RemoteAddr
	h := req.Header

	if realIpHeader := h.Get("X-Real-Ip"); realIpHeader != "" {
		ip = realIpHeader
	}

	if forwardIpHeader := h.Get("X-Forwarded-For"); forwardIpHeader != "" {
		ip = forwardIpHeader
	}

	return ip
}
