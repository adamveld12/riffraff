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
	enableAccessLogging := flag.Bool("accesslog", true, "Enable access logging")
	flag.Parse()

	box := packr.NewBox("./internal/templates")
	tp := internal.TemplateRenderer{FS: box}

	server := internal.NewServer(tp, *enableAccessLogging)

	log.SetPrefix("[INFO] ")

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Printf("Listening @ %s", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
