package controllers

import "github.com/mpanelo/gocookit/views"

type Static struct {
	Home *views.View
	// About   *views.View
	// Contact *views.View
}

func NewStatic() *Static {
	return &Static{
		Home: views.NewView("static/home"),
	}
}
