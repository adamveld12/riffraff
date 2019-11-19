package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/adamveld12/riffraff/internal"
	"github.com/gobuffalo/packr"
)

func main() {
	port := flag.Int("port", 80, "port to listen on")
	bindAddr := flag.String("bind", "127.0.0.1", "interface to bind to")
	dbPath := flag.String("data", "./data.json", "path to save shortcut database")
	enableAccessLogging := flag.Bool("accesslog", true, "Enable access logging")
	flag.Parse()

	box := packr.NewBox("./internal/templates")
	tp := internal.TemplateRenderer{FS: box}

	ss := &internal.ShortcutStore{Path: *dbPath}

	if err := ss.Init(); err != nil {
		log.Fatalf("could not access database file: %v", err)
	}

	server := internal.NewServer(tp, ss, *enableAccessLogging)

	log.SetPrefix("[INFO] ")

	addr := fmt.Sprintf("%s:%d", *bindAddr, *port)
	log.Printf("Listening @ %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
