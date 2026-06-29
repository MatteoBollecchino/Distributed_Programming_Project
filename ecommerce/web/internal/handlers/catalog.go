package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	pbCatalog "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
)

func (s *ServerDependencies) CatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// Calling catalog service via gRPC
	catalogRes, err := s.Clients.Catalog.ListCatalogItems(request.Context(), &pbCatalog.ListCatalogItemsRequest{})
	if !checkerr(writer, err) {
		return
	}

	// Get current session
	session, err := s.Store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return
	}

	// Retrieving if user is logged in or not
	isLoggedIn, _ := session.Values["logged_in"].(bool)

	// Map with data to send to HTML file
	templateData := map[string]interface{}{
		"Title":      "Fanta Catalog",
		"Products":   catalogRes.GetItems(), // List all the products from gRPC
		"IsLoggedIn": isLoggedIn,
	}

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "catalog.html", templateData))
}

func (s *ServerDependencies) UpdateCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}
	role := session.Values["role"].(string)

	templateData := map[string]interface{}{
		"Role":  role,
		"Admin": "ADMIN",
	}

	// Only GET requests are accepted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "update_catalog.html", templateData))
}

func (s *ServerDependencies) AddToCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}
	username := session.Values["username"]
	role := session.Values["role"].(string)

	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is an admin
	if role != "ADMIN" {
		checkerr(writer, errors.New("User Must be an admin to do this operation"))
		return
	}

	// Retrieve item data
	itemId := request.FormValue("item_id")
	description := request.FormValue("description")
	priceStr := request.FormValue("price")
	quantityStr := request.FormValue("quantity")

	price, err := strconv.ParseFloat(priceStr, 64)
	if !checkerr(writer, err) {
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if !checkerr(writer, err) {
		return
	}

	// Creating catalog item
	item := pbCatalog.CatalogItem{
		ItemId:            itemId,
		Description:       description,
		Price:             price,
		QuantityAvailable: uint32(quantity),
	}

	// Calling catalog service via gRPC
	_, err = s.Clients.Catalog.AddCatalogItem(request.Context(), &pbCatalog.AddCatalogItemRequest{
		Item: &item,
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.Manager.NotifyCatalogUpdate()

	log.Printf("New Item successfully added to catalog by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

func (s *ServerDependencies) RemoveFromCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}
	username := session.Values["username"]
	role := session.Values["role"].(string)

	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is an admin
	if role != "ADMIN" {
		checkerr(writer, errors.New("User Must be an admin to do this operation"))
		return
	}

	// Retrieve item data
	itemId := request.FormValue("item_id")

	// Calling catalog service via gRPC
	_, err := s.Clients.Catalog.RemoveCatalogItem(request.Context(), &pbCatalog.RemoveCatalogItemRequest{
		ItemId: itemId,
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.Manager.NotifyCatalogUpdate()

	log.Printf("Item successfully removed from catalog by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

func (s *ServerDependencies) UpdatePriceCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}
	username := session.Values["username"]
	role := session.Values["role"].(string)

	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is an admin
	if role != "ADMIN" {
		checkerr(writer, errors.New("User Must be an admin to do this operation"))
		return
	}

	// Retrieve item data
	itemId := request.FormValue("item_id")
	priceStr := request.FormValue("price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if !checkerr(writer, err) {
		return
	}

	// Calling catalog service via gRPC
	_, err = s.Clients.Catalog.UpdatePrice(request.Context(), &pbCatalog.UpdatePriceRequest{
		ItemId: itemId,
		Price:  price,
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.Manager.NotifyCatalogUpdate()

	log.Printf("Item Price successfully updated by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

func (s *ServerDependencies) UpdateQuantityCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}
	username := session.Values["username"]
	role := session.Values["role"].(string)

	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is an admin
	if role != "ADMIN" {
		checkerr(writer, errors.New("User Must be an admin to do this operation"))
		return
	}

	// Retrieve item data
	itemId := request.FormValue("item_id")
	quantityStr := request.FormValue("quantity")

	quantity, err := strconv.Atoi(quantityStr)
	if !checkerr(writer, err) {
		return
	}

	// Calling catalog service via gRPC
	_, err = s.Clients.Catalog.UpdateQuantityAvailable(request.Context(), &pbCatalog.UpdateQuantityAvailableRequest{
		ItemId:   itemId,
		Quantity: uint32(quantity),
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.Manager.NotifyCatalogUpdate()

	log.Printf("Item Quantity successfully updated by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}
