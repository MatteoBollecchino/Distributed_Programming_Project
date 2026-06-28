package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/clients"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/manager"
	"github.com/gorilla/sessions"
)

const sessionName = "ecommerce-session"

// ServerDependencies all of what handlers need
type ServerDependencies struct {
	Templates *template.Template
	Clients   *clients.ServiceClients
	Store     *sessions.CookieStore
	Manager   *manager.EventsManager
}

func checkerr(writer http.ResponseWriter, err error) bool {
	ok := true
	if err != nil {
		ok = false
		log.Printf("Error occured: %v", err.Error())
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	return ok
}

func checkIfUserIsLogged(s *ServerDependencies, request *http.Request, writer http.ResponseWriter) (*sessions.Session, bool) {
	session, err := s.Store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return nil, false
	}
	if loggedIn, ok := session.Values["logged_in"].(bool); !ok || !loggedIn {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return nil, false
	}
	return session, true
}
