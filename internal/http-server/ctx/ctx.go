package ctx

import (
	"context"
	"net/http"
)

// used as key for context.WithValue
type contextKey string

// most used keys in ctx
const (
	CtxSessionIDKey contextKey 	= "session_id"
	CtxUserIDKey 	contextKey	= "user_id"
	CtxFinishedKey 	contextKey 	= "finished"
)

func WrapValueIntoRequest(r *http.Request, key contextKey, value any) *http.Request {
	newCtx := context.WithValue(r.Context(), key, value)
	return r.WithContext(newCtx)
}