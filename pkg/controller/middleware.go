package controller

import (
	"context"
	"net/http"

	"github.com/bastianrob/restsuite/pkg/ctxkey"

	"github.com/bastianrob/go-httputil/middleware"
	oauth "github.com/bastianrob/go-oauth/handler"
	"github.com/bastianrob/go-oauth/model"
)

//UnwrapJWT JWT custom claims and store organization name to context
func UnwrapJWT() middleware.HTTPMiddleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authClaims, ok := ctx.Value(oauth.ContextKeyClaims).(model.AuthClaims)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}

			organizationName, ok := authClaims.CustomClaims["org"]
			if !ok {
				h.ServeHTTP(w, r)
				return
			}

			ctx = context.WithValue(ctx, ctxkey.OrganizationName, organizationName)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		}
	}
}
