package handlers

import "net/http"

// Simple ping handler to check health of service and test if server is available externally
func (h *Handlers) PingHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("pong"))
    }
}