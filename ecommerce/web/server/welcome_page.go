package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates = loadTemplates()

func checkerr(writer http.ResponseWriter, err error) {
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func loadTemplates() *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	basePath := filepath.Dir(wd)
	templatesPath := filepath.Join(basePath, "templates", "*.html")

	return template.Must(template.ParseGlob(templatesPath))
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, templates.ExecuteTemplate(writer, "welcome_page.html", nil))
}

func productCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, templates.ExecuteTemplate(writer, "product_catalog.html", nil))
}

func shoppingCartHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, templates.ExecuteTemplate(writer, "shopping_cart.html", nil))
}

func accountHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, templates.ExecuteTemplate(writer, "account.html", nil))
}

func main() {
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/product_catalog", productCatalogHandler)
	http.HandleFunc("/shopping_cart", shoppingCartHandler)
	http.HandleFunc("/account", accountHandler)

	// starts the server on the port 8080 and nil means that the predefined router has to be used
	log.Fatal(http.ListenAndServe(":8080", nil))
}
