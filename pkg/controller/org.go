package controller

import "github.com/bastianrob/go-httputil/middleware"

//OrganizationController HTTP API Contract
type OrganizationController interface {
	Add() middleware.HTTPMiddleware
	Get() middleware.HTTPMiddleware
	Update() middleware.HTTPMiddleware
}
