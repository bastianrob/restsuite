package controller

import "github.com/bastianrob/go-httputil/middleware"

//ScenarioController HTTP API Contract
type ScenarioController interface {
	Find() middleware.HTTPMiddleware
	Get() middleware.HTTPMiddleware

	Add() middleware.HTTPMiddleware
	Update() middleware.HTTPMiddleware
	Delete() middleware.HTTPMiddleware
}
