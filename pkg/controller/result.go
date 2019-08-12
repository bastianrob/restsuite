package controller

import "github.com/bastianrob/go-httputil/middleware"

//ResultController HTTP API Contract
type ResultController interface {
	Find() middleware.HTTPMiddleware
}
