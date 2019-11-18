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
	cfgPath := flag.String("config", "./config.json", "path to config file")
	enableAccessLogging := flag.Bool("accesslog", true, "Enable access logging")
	flag.Parse()

	box := packr.NewBox("./internal/templates")
	tp := internal.TemplateRenderer{FS: box}

	config, err := internal.LoadShortcutsFromConfig(*cfgPath)
	if err != nil  {
		log.Fatalf("could not load config from file: %v", err)
	}

	server := internal.NewServer(tp, config.Shortcuts, *enableAccessLogging)

	log.SetPrefix("[INFO] ")

	addr := fmt.Sprintf("%s:%d", config.ListenAddress, *port)
	log.Printf("Listening @ %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
