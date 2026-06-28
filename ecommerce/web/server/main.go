package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/clients"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/handlers"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/manager"
	"github.com/gorilla/sessions"
)

var port = ":8080"
var authKey = []byte("FantaEcommerce2026SecureAuthKey1") // For authetication
var encKey = []byte("FantaEcommerce2026EncryptionKey1")  // For encryption

type WebServer struct {
	dep *handlers.ServerDependencies
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

// WELCOME PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.WelcomeHandler(writer, request)
}

// CATALOG PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) catalogHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.CatalogHandler(writer, request)
}

func (s *WebServer) updateCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.UpdateCatalogHandler(writer, request)
}

func (s *WebServer) addToCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.AddToCatalogHandler(writer, request)
}

func (s *WebServer) removeFromCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.RemoveFromCatalogHandler(writer, request)
}

func (s *WebServer) updatePriceCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.UpdatePriceCatalogHandler(writer, request)
}

func (s *WebServer) updateQuantityCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.UpdateQuantityCatalogHandler(writer, request)
}

// CART PAGE HANDLERS ///////////////////////////////////////////////////////////////

func (s *WebServer) cartHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.CartHandler(writer, request)
}

func (s *WebServer) addToCartHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.AddToCartHandler(writer, request)
}

func (s *WebServer) removeFromCartHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.RemoveFromCartHandler(writer, request)
}

func (s *WebServer) updateQuantityCartHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.UpdateQuantityCartHandler(writer, request)
}

// ORDER PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) orderHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.OrderHandler(writer, request)
}

func (s *WebServer) userOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.UserOrdersHandler(writer, request)
}

// PAYMENT PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) paymentHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.PaymentHandler(writer, request)
}

func (s *WebServer) processPaymentHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.ProcessPaymentHandler(writer, request)
}

// AUTHETIFICATION PAGE HANDLERS ///////////////////////////////////////////////////////////////

func (s *WebServer) accountHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.AccountHandler(writer, request)
}

func (s *WebServer) registerHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.RegisterHandler(writer, request)
}

func (s *WebServer) loginHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.LoginHandler(writer, request)
}

func (s *WebServer) logoutHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.LogoutHandler(writer, request)
}

func (s *WebServer) changePasswordHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.ChangePasswordHandler(writer, request)
}

func (s *WebServer) listAllUsersHandler(writer http.ResponseWriter, request *http.Request) {
	s.dep.ListAllUsersHandler(writer, request)
}

// MAIN ///////////////////////////////////////////////////////////////

func main() {
	// Initialization gRPC clients
	clientsRegistry, err := clients.InitClients()
	if err != nil {
		log.Fatal(err)
	}
	defer clientsRegistry.Close()

	// Manager of the synchronization between server and client/browser
	eventsManager := manager.NewEventsManager()

	// Cookies creation
	cookieStore := sessions.NewCookieStore(authKey, encKey)

	// Cookies settings
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	// Dependencies creation
	dependecies := &handlers.ServerDependencies{
		Templates: loadTemplates(),
		Clients:   clientsRegistry,
		Store:     cookieStore,
		Manager:   eventsManager,
	}

	// Web server creation
	server := &WebServer{
		dep: dependecies,
	}

	mux := http.NewServeMux()

	// Association of paths to correspondent handlers
	mux.HandleFunc("/welcome", server.welcomeHandler)
	mux.HandleFunc("/events", server.dep.Manager.HandleEvents)
	mux.HandleFunc("/catalog", server.catalogHandler)
	mux.HandleFunc("/catalog/add", server.addToCatalogHandler)
	mux.HandleFunc("/catalog/remove", server.removeFromCatalogHandler)
	mux.HandleFunc("/catalog/update/price", server.updatePriceCatalogHandler)
	mux.HandleFunc("/catalog/update/quantity", server.updateQuantityCatalogHandler)
	mux.HandleFunc("/update/catalog", server.updateCatalogHandler)
	mux.HandleFunc("/cart", server.cartHandler)
	mux.HandleFunc("/cart/add", server.addToCartHandler)
	mux.HandleFunc("/cart/remove", server.removeFromCartHandler)
	mux.HandleFunc("/cart/update", server.updateQuantityCartHandler)
	mux.HandleFunc("/order", server.orderHandler)
	mux.HandleFunc("/user/orders", server.userOrdersHandler)
	mux.HandleFunc("/payment", server.paymentHandler)
	mux.HandleFunc("/payment/process", server.processPaymentHandler)
	mux.HandleFunc("/account", server.accountHandler)
	mux.HandleFunc("/register", server.registerHandler)
	mux.HandleFunc("/login", server.loginHandler)
	mux.HandleFunc("/logout", server.logoutHandler)
	mux.HandleFunc("/change/password", server.changePasswordHandler)
	mux.HandleFunc("/list/users", server.listAllUsersHandler)

	log.Printf("The Web Server listening on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
