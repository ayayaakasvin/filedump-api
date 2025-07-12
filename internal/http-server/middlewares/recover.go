package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (mw *Middlewares) RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				msg := fmt.Sprintf("panic recovered: %v", rec)
				mw.logger.Errorf("%s%s\n%s", ColorRed, msg, debug.Stack())

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}
