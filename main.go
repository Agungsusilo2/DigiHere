package main

import (
	"DigiHero_Web/controllers"
	"log"
	"net/http"
)

func main() {
	postController := &controllers.PostImp{} // Creating an instance of the controller

	http.HandleFunc("/", postController.Index)
	http.HandleFunc("/posts", postController.Index)
	http.HandleFunc("/post/create", postController.Create)
	http.HandleFunc("/post/store", postController.Store)
	http.HandleFunc("/post/delete", postController.Delete)

	log.Print("Server started on: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
