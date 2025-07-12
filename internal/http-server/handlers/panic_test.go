package handlers

import "net/http"

// test handler for panic and recover from it
func (h *Handlers) PanicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}
}