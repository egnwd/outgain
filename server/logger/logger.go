package logger

import (
	"io"
	"log"
	"net/http"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// ServerLogger is wrapper for a http.Handler that logs output
// to the specified writer
func ServerLogger(l io.Writer, next http.Handler) http.Handler {
	log.SetOutput(l)

	fn := func(w http.ResponseWriter, r *http.Request) {
		var logger logResponseWriter = &responseLogger{w: w}
		next.ServeHTTP(logger, r)

		writeLog(r.URL.Path, logger.Status(), r.Method)
	}

	return http.HandlerFunc(fn)
}

func writeLog(uri string, status int, method string) {
	statusColour := colourForStatus(status)
	methodColour := colourForMethod(method)
	log.Printf("|%s %d %s| %s %s %s - %s",
		statusColour, status, reset,
		methodColour, method, reset,
		uri,
	)
}

func colourForStatus(status int) string {
	switch {
	case status < 300:
		return green
	case 300 <= status && status < 400:
		return white
	case 400 <= status && status < 500:
		return yellow
	default:
		return red
	}
}

func colourForMethod(method string) string {
	switch method {
	case http.MethodGet:
		return cyan
	case http.MethodPost:
		return yellow
	case http.MethodPut:
		return blue
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

type logResponseWriter interface {
	http.ResponseWriter
	Status() int
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
}

func (l *responseLogger) Status() int {
	if l.status == 0 {
		return http.StatusOK
	}

	return l.status
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(bs []byte) (int, error) {
	if l.status == 0 {
		l.status = http.StatusOK
	}
	return l.w.Write(bs)
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}
