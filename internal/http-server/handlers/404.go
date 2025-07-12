package handlers

import "net/http"

const HtmlFilePath = "./static/404.html"

// Handler for 404 case, no route matches request path
func (h *Handlers) NotFound404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, HtmlFilePath)
	}
}
