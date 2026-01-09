package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	err := templates.ExecuteTemplate(writer, "welcome_page.html", nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/footer", welcomeHandler)

	// starts the server on the port 8080 and nil means that the predefined router has to be used
	log.Fatal(http.ListenAndServe(":8080", nil))
}
