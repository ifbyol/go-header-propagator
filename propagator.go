package goheaderpropagator

import (
	"context"
	"net/http"
)

// Propagator struct to be used to get a header into the go context, and also as transport for http client
type Propagator struct {
	Header string
	Base   http.RoundTripper
}

type ctxKey struct{}

// Middleware it injects the header configured in the Propagator fo the golang context in the HTTP request
func (p *Propagator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		customHeader := r.Header.Get(p.Header)

		ctx = context.WithValue(ctx, ctxKey{}, customHeader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoundTrip it adds the header to the request if it is present in the context
func (p *Propagator) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())

	if val := req.Context().Value(ctxKey{}); val != nil {
		req.Header.Set(p.Header, val.(string))
	}

	base := p.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}
