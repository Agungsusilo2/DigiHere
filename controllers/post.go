package controllers

import "net/http"

type Post interface {
	Store(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Index(w http.ResponseWriter, r *http.Request)
}
