package main

import (
	"DigiHero_Web/controllers"
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed uploads
var uploads embed.FS

func main() {
	sub, err := fs.Sub(uploads, "uploads")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(sub))
	postController := &controllers.PostImp{}

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", postController.Index)
	mux.HandleFunc("/posts", postController.Index)
	mux.HandleFunc("/post/create", postController.Create)
	mux.HandleFunc("/post/store", postController.Store)
	mux.HandleFunc("/post/delete", postController.Delete)
	mux.HandleFunc("/post/download", postController.DownloadFile)

	server := http.Server{Addr: "localhost:8000", Handler: mux}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	log.Print("Server started on: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
