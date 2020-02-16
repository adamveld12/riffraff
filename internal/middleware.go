package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

func AccessLoggerMiddleware(middleware http.Handler) http.Handler {
	accessLogger := log.New(os.Stdout, "[access] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		requestId := uuid.Must(uuid.NewV4())
		setupHeaders(res.Header(), requestId.String())
		accessLogger.Printf("[%s] %s %s '%s'", requestId, getRemoteIP(req), req.Method, req.URL.String())
		middleware.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), "id", requestId.String())))
	})
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
