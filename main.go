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
	bindAddr := flag.String("bind", "0.0.0.0", "interface to bind to")
	enableAccessLogging := flag.Bool("accesslog", true, "Enable access logging")
	flag.Parse()

	box := packr.NewBox("./internal/templates")
	tp := internal.TemplateRenderer{FS: box}

	config, err := internal.LoadConfig(*cfgPath)
	if err != nil  {
		log.Fatalf("could not load config from file: %v", err)
	}

	server := internal.NewServer(tp, config.Shortcuts, *enableAccessLogging)

	log.SetPrefix("[INFO] ")

	addr := fmt.Sprintf("%s:%d", *bindAddr, *port)
	log.Printf("Listening @ %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
