package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	pbCart "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
)

func (s *ServerDependencies) CartHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are accepted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve session
	session, err := s.Store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return
	}

	// Checking if user is logged
	if loggedIn, ok := session.Values["logged_in"].(bool); !ok || !loggedIn {
		// Redirection to login page
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}
	username := session.Values["username"].(string)

	// Calling cart service via gRPC
	cartRes, err := s.Clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{
		Username: username,
	})
	if err != nil {
		log.Printf("Impossible to retrieve shopping cart for %s: %v", username, err)
		checkerr(writer, s.Templates.ExecuteTemplate(writer, "cart.html", nil))
		return
	}

	// Calculating total price
	totalPriceRes, err := s.Clients.Cart.CalculateTotalPrice(request.Context(), &pbCart.CalculateTotalPriceRequest{
		Username: username,
	})
	if !checkerr(writer, err) {
		return
	}

	errorMessage := ""
	queryError := request.URL.Query().Get("error")

	if queryError == "catalog_changed" {
		errorMessage = "Catalog has been updated. Items could have been changed or removed."
	}
	if queryError == "payment_failed" {
		errorMessage = "Failed payment: the amount provided was insufficient"
	}

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"Items":      cartRes.GetCart().GetItems(),
		"TotalPrice": math.Trunc(totalPriceRes.GetTotalPrice()*100) / 100,
		"Error":      errorMessage,
	}

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "cart.html", templateData))
}

func (s *ServerDependencies) AddToCartHandler(writer http.ResponseWriter, request *http.Request) {
	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	// Retrieve user and product data
	productId := request.FormValue("product_id")
	priceStr := request.FormValue("price")
	username := session.Values["username"].(string)
	quantityStr := request.FormValue("quantity")

	// Price conversion
	price, err := strconv.ParseFloat(priceStr, 64)
	if !checkerr(writer, err) {
		return
	}

	// Quantity Conversion
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		log.Printf("Error in quantity conversion '%s': %v", quantityStr, err)
		http.Error(writer, "Quantity not valid", http.StatusBadRequest)
		return
	}

	// gRPC call at Cart service
	_, err = s.Clients.Cart.AddItemToCart(request.Context(), &pbCart.AddItemToCartRequest{
		Username: username,
		CartItem: &pbCart.CartItem{
			ItemId:   productId,
			Price:    float64(price),
			Quantity: uint32(quantity)},
	})

	if err != nil {
		log.Printf("Failed add item to cart: %v", err)
		checkerr(writer, s.Templates.ExecuteTemplate(writer, "cart.html", "Failed add item to cart"))
		return
	}

	log.Printf("Product added successfully to cart")

	// Redirection to shopping cart page
	http.Redirect(writer, request, "/cart", http.StatusSeeOther)
}

func (s *ServerDependencies) RemoveFromCartHandler(writer http.ResponseWriter, request *http.Request) {
	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	// Retrieve user and product data
	productId := request.FormValue("product_id")
	username := session.Values["username"].(string)

	// gRPC call at Cart service
	_, err := s.Clients.Cart.RemoveItemFromCart(request.Context(), &pbCart.RemoveItemFromCartRequest{
		Username: username,
		ItemId:   productId,
	})

	if err != nil {
		log.Printf("Failed remove item from cart: %v", err)
		checkerr(writer, s.Templates.ExecuteTemplate(writer, "cart.html", "Failed remove item from cart"))
		return
	}

	log.Printf("Product removed successfully from cart")

	// Redirection to shopping cart page
	http.Redirect(writer, request, "/cart", http.StatusSeeOther)
}

func (s *ServerDependencies) UpdateQuantityCartHandler(writer http.ResponseWriter, request *http.Request) {
	// Only POST requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	// Retrieve user and product data
	productId := request.FormValue("product_id")
	username := session.Values["username"].(string)
	quantityStr := request.FormValue("quantity")

	quantity, err := strconv.Atoi(quantityStr)
	if !checkerr(writer, err) {
		return
	}

	// gRPC call at Cart service
	_, err = s.Clients.Cart.UpdateItemQuantity(request.Context(), &pbCart.UpdateItemQuantityRequest{
		Username: username,
		ItemId:   productId,
		Quantity: uint32(quantity),
	})
	if !checkerr(writer, err) {
		return
	}

	log.Printf("Product quantity successfully updated")

	// Redirection to shopping cart page
	http.Redirect(writer, request, "/cart", http.StatusSeeOther)
}
