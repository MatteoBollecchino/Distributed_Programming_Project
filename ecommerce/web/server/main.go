package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	pbCatalog "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/clients"
)

var port = ":8080"

// 1. Creiamo un'ambiente (Server) per fare Dependency Injection
type WebServer struct {
	templates *template.Template
	clients   *clients.ServiceClients
}

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

// 2. Trasformiamo gli handler in METODI di WebServer. Ora hanno accesso a s.clients!

func (s *WebServer) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "welcome_page.html", nil))
}

func (s *WebServer) productCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// Request products from the catalog microservice via gRPC using s.clients
	catalogRes, err := s.clients.Catalog.ListCatalogItems(request.Context(), &pbCatalog.ListCatalogItemsRequest{})
	checkerr(writer, err)

	// Map with data to send to HTML file
	templateData := map[string]interface{}{
		"Title":    "Fanta Catalog",
		"Products": catalogRes.GetItems(), // List all the products from gRPC
	}

	checkerr(writer, s.templates.ExecuteTemplate(writer, "catalog.html", templateData))
}

func (s *WebServer) shoppingCartHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "shopping_cart.html", nil))
}

func (s *WebServer) accountHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "account.html", nil))
}

func (s *WebServer) registerHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "register.html", nil))
}

func (s *WebServer) loginHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "login.html", nil))
}

func main() {
	// Inizializza i client gRPC
	clientsRegistry, err := clients.InitClients()
	if err != nil {
		log.Fatal(err)
	}
	defer clientsRegistry.Close()

	// 3. Istanziamo il nostro WebServer con i template e i client
	server := &WebServer{
		templates: loadTemplates(),
		clients:   clientsRegistry,
	}

	// 4. Usiamo il Mux (consigliato dai tuoi commenti!)
	mux := http.NewServeMux()

	// Associamo i percorsi ai metodi dell'istanza 'server'
	mux.HandleFunc("/welcome", server.welcomeHandler)
	mux.HandleFunc("/catalog", server.productCatalogHandler)
	mux.HandleFunc("/shopping_cart", server.shoppingCartHandler)
	mux.HandleFunc("/account", server.accountHandler)
	mux.HandleFunc("/register", server.registerHandler)
	mux.HandleFunc("/login", server.loginHandler)

	log.Printf("The Web Server listening on %s", port)
	// Passiamo il mux personalizzato invece di nil
	log.Fatal(http.ListenAndServe(port, mux))
}
