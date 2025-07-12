package middlewares

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var randomizer *rand.Rand

func init() {
	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
}

const (
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorReset   = "\033[0m"
)

func (m *Middlewares) LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		coloredReqID := fmt.Sprintf("%s%s\033[0m", randomColorOfReqID(), reqID)
		m.logger.Info(requestInfo(r, coloredReqID))
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			result := "failed"
			color := "\033[31m"
			duration := time.Since(start)

			if rw.finished {
				result = "successfully"
				color = "\033[32m"
			}

			m.logger.Infof("[END]\n\tReqID=%s\n\tStatus=%d\n\tDuration=%s\n\tResult=%s%s\033[0m\n",
				coloredReqID,
				rw.statusCode,
				duration,
				color,
				result,
			)
		}()

		next.ServeHTTP(rw, r)
	}
}

func requestInfo(r *http.Request, coloredReqID string) string {
	return fmt.Sprintf(
		"[START]\n\tReqID=%s\n\tMethod=%s\n\tURL=%s\n\tRemoteAddr=%s\n\tUserAgent=%s\n\tHeaders=%v\n",
		coloredReqID,
		r.Method,
		r.URL.String(),
		r.RemoteAddr,
		r.UserAgent(),
		r.Header,
	)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	finished   bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.finished = true
	rw.ResponseWriter.WriteHeader(code)
}

func randomColorOfReqID() string {
	colors := []string{
		ColorRed,
		ColorGreen,
		ColorYellow,
		ColorBlue,
		ColorMagenta,
		ColorCyan,
	}
	return colors[randomizer.Intn(len(colors))]
}