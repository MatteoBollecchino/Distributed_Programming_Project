package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	pbAuth "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	pbCatalog "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/clients"
	"github.com/gorilla/sessions"
)

var port = ":8080"
var cookieKey = []byte("FantaEcommerce2026")

const sessionName = "ecommerce-session"

type WebServer struct {
	templates *template.Template
	clients   *clients.ServiceClients
	store     *sessions.CookieStore
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

func (s *WebServer) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "welcome_page.html", nil))
}

func (s *WebServer) productCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// Request products from the catalog microservice via gRPC using s.clients
	catalogRes, err := s.clients.Catalog.ListCatalogItems(request.Context(), &pbCatalog.ListCatalogItemsRequest{})
	checkerr(writer, err)

	// Get current session
	session, err := s.store.Get(request, sessionName)
	checkerr(writer, err)

	isLoggedIn := false
	username := ""

	// checking if user is logged or not
	if loggedIn, ok := session.Values["logged_in"].(bool); ok && loggedIn {
		isLoggedIn = true
		username = session.Values["username"].(string)
	}

	// Map with data to send to HTML file
	templateData := map[string]interface{}{
		"Title":      "Fanta Catalog",
		"Products":   catalogRes.GetItems(), // List all the products from gRPC
		"IsLoggedIn": isLoggedIn,
		"Username":   username,
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
	// User request GET -> login page
	if request.Method == http.MethodGet {
		checkerr(writer, s.templates.ExecuteTemplate(writer, "login.html", nil))
		return
	}

	// User request POST, l'utente ha sottomesso il form
	if request.Method == http.MethodPost {
		username := request.FormValue("username")
		password := request.FormValue("password")

		// gRPC call at Auth service
		authRes, err := s.clients.Auth.Login(request.Context(), &pbAuth.LoginRequest{
			Username: username,
			Password: password,
		})

		// In case of errore -> redirection to login page
		if err != nil {
			// Nota: in un secondo momento potrai passare l'errore al template per mostrarlo all'utente
			log.Printf("Failed Login for %s: %v", username, err)
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}

		// Session creation
		session, _ := s.store.Get(request, sessionName)

		// Save user data in the session
		session.Values["username"] = authRes.GetUser().Username
		session.Values["role"] = authRes.GetUser().Role
		session.Values["logged_in"] = true

		// Saving session
		if err = session.Save(request, writer); err != nil {
			http.Error(writer, "Errore nel salvataggio della sessione", http.StatusInternalServerError)
			return
		}

		// Redirection to catalog page
		http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
	}
}

func main() {
	// Initialization gRPC clients
	clientsRegistry, err := clients.InitClients()
	if err != nil {
		log.Fatal(err)
	}
	defer clientsRegistry.Close()

	cookieStore := sessions.NewCookieStore(cookieKey)

	server := &WebServer{
		templates: loadTemplates(),
		clients:   clientsRegistry,
		store:     cookieStore,
	}

	mux := http.NewServeMux()

	// Association of paths to correspondent handlers
	mux.HandleFunc("/welcome", server.welcomeHandler)
	mux.HandleFunc("/catalog", server.productCatalogHandler)
	mux.HandleFunc("/shopping_cart", server.shoppingCartHandler)
	mux.HandleFunc("/account", server.accountHandler)
	mux.HandleFunc("/register", server.registerHandler)
	mux.HandleFunc("/login", server.loginHandler)

	log.Printf("The Web Server listening on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
