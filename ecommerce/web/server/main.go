package main

import (
	"errors"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	pbAuth "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	pbCart "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	pbCatalog "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	pbOrder "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
	pbPayment "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/clients"
	manager "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/web/internal/manager"
	"github.com/gorilla/sessions"
)

var port = ":8080"
var cookieKey = []byte("FantaEcommerce2026")

const sessionName = "ecommerce-session"

type WebServer struct {
	templates *template.Template
	clients   *clients.ServiceClients
	store     *sessions.CookieStore
	manager   *manager.EventsManager
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

func loadTemplates() *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	basePath := filepath.Dir(wd)
	templatesPath := filepath.Join(basePath, "templates", "*.html")

	return template.Must(template.ParseGlob(templatesPath))
}

func checkIfUserIsLogged(s *WebServer, request *http.Request, writer http.ResponseWriter) (*sessions.Session, bool) {
	session, err := s.store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return nil, false
	}
	if loggedIn, ok := session.Values["logged_in"].(bool); !ok || !loggedIn {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return nil, false
	}
	return session, true
}

// WELCOME PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	if !checkerr(writer, s.templates.ExecuteTemplate(writer, "welcome_page.html", nil)) {
		return
	}
}

// CATALOG PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) catalogHandler(writer http.ResponseWriter, request *http.Request) {
	// Request products from the catalog microservice via gRPC using s.clients
	catalogRes, err := s.clients.Catalog.ListCatalogItems(request.Context(), &pbCatalog.ListCatalogItemsRequest{})
	if !checkerr(writer, err) {
		return
	}

	// Get current session
	session, err := s.store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return
	}

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

func (s *WebServer) updateCatalogHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be
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

	checkerr(writer, s.templates.ExecuteTemplate(writer, "update_catalog.html", templateData))
}

func (s *WebServer) addToCatalogHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err = s.clients.Catalog.AddCatalogItem(request.Context(), &pbCatalog.AddCatalogItemRequest{
		Item: &item,
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.manager.NotifyCatalogUpdate()

	log.Printf("New Item successfully added to catalog by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

func (s *WebServer) removeFromCatalogHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err := s.clients.Catalog.RemoveCatalogItem(request.Context(), &pbCatalog.RemoveCatalogItemRequest{
		ItemId: itemId,
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.manager.NotifyCatalogUpdate()

	log.Printf("Item successfully removed from catalog by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

func (s *WebServer) updatePriceCatalogHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err = s.clients.Catalog.UpdatePrice(request.Context(), &pbCatalog.UpdatePriceRequest{
		ItemId: itemId,
		Price:  price,
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.manager.NotifyCatalogUpdate()

	log.Printf("Item Price successfully updated by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

func (s *WebServer) updateQuantityCatalogHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err = s.clients.Catalog.UpdateQuantityAvailable(request.Context(), &pbCatalog.UpdateQuantityAvailableRequest{
		ItemId:   itemId,
		Quantity: uint32(quantity),
	})
	if !checkerr(writer, err) {
		return
	}

	// Notification that the catalog has changed
	s.manager.NotifyCatalogUpdate()

	log.Printf("Item Price successfully updated by %s", username)

	// Redirection to catalog page
	http.Redirect(writer, request, "/catalog", http.StatusSeeOther)
}

// CART PAGE HANDLERS ///////////////////////////////////////////////////////////////

func (s *WebServer) cartHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are accepted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve session
	session, err := s.store.Get(request, sessionName)
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
	cartRes, err := s.clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{
		Username: username,
	})
	if err != nil {
		log.Printf("Impossible to retrieve shopping cart for %s: %v", username, err)
		checkerr(writer, s.templates.ExecuteTemplate(writer, "cart.html", nil))
		return
	}

	// Calculating total price
	totalPriceRes, err := s.clients.Cart.CalculateTotalPrice(request.Context(), &pbCart.CalculateTotalPriceRequest{
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

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"Items":      cartRes.GetCart().GetItems(),
		"TotalPrice": math.Trunc(totalPriceRes.GetTotalPrice()*100) / 100,
		"Error":      errorMessage,
	}

	checkerr(writer, s.templates.ExecuteTemplate(writer, "cart.html", templateData))
}

func (s *WebServer) addToCartHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err = s.clients.Cart.AddItemToCart(request.Context(), &pbCart.AddItemToCartRequest{
		Username: username,
		CartItem: &pbCart.CartItem{
			ItemId:   productId,
			Price:    float64(price),
			Quantity: uint32(quantity)},
	})

	if err != nil {
		log.Printf("Failed add item to cart: %v", err)
		checkerr(writer, s.templates.ExecuteTemplate(writer, "cart.html", "Failed add item to cart"))
		return
	}

	log.Printf("Product added successfully to cart")

	// Redirection to shopping cart page
	http.Redirect(writer, request, "/cart", http.StatusSeeOther)
}

func (s *WebServer) removeFromCartHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err := s.clients.Cart.RemoveItemFromCart(request.Context(), &pbCart.RemoveItemFromCartRequest{
		Username: username,
		ItemId:   productId,
	})

	if err != nil {
		log.Printf("Failed remove item from cart: %v", err)
		checkerr(writer, s.templates.ExecuteTemplate(writer, "cart.html", "Failed remove item from cart"))
		return
	}

	log.Printf("Product removed successfully from cart")

	// Redirection to shopping cart page
	http.Redirect(writer, request, "/cart", http.StatusSeeOther)
}

func (s *WebServer) updateQuantityCartHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err = s.clients.Cart.UpdateItemQuantity(request.Context(), &pbCart.UpdateItemQuantityRequest{
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

// ORDER PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) orderHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are accepted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}
	username := session.Values["username"].(string)

	// gRPC call at Cart service to retrieve the cart
	cartRes, err := s.clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{
		Username: username,
	})
	if !checkerr(writer, err) {
		return
	}

	// Calculate total price
	totalPriceRes, err := s.clients.Cart.CalculateTotalPrice(request.Context(), &pbCart.CalculateTotalPriceRequest{
		Username: username,
	})
	if !checkerr(writer, err) {
		return
	}

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"Items":      cartRes.GetCart().GetItems(),
		"TotalPrice": math.Trunc(totalPriceRes.GetTotalPrice()*100) / 100,
	}

	if !checkerr(writer, s.templates.ExecuteTemplate(writer, "order.html", templateData)) {
		return
	}
}

func (s *WebServer) userOrdersHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	_, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	// User GET request -> user's orders page
	if request.Method == http.MethodGet {

		// Retrieving username from URL
		username := request.URL.Query().Get("username")

		// gRPC call at Order service to retrieve all the order for a certain user
		orderRes, err := s.clients.Order.ListOrdersByUser(request.Context(), &pbOrder.ListOrdersByUserRequest{
			UserId: username,
		})
		if !checkerr(writer, err) {
			return
		}

		// Mapping data for HTML file
		templateData := map[string]interface{}{
			"Username": username,
			"Orders":   orderRes.GetOrders(),
		}

		log.Printf("List of %s's orders successfully retrieved", username)

		checkerr(writer, s.templates.ExecuteTemplate(writer, "user_orders.html", templateData))
	}

	// User POST request -> update order status
	if request.Method == http.MethodPost {

		// Retrieving the order info
		newStatusStr := request.FormValue("new_status")
		newStatusValue, err := strconv.Atoi(newStatusStr)
		if !checkerr(writer, err) {
			return
		}

		orderId := request.FormValue("order_id")

		// gRPC call at Order service to update the status of the order
		_, err = s.clients.Order.UpdateOrderStatus(request.Context(), &pbOrder.UpdateOrderStatusRequest{
			OrderId: orderId,
			Status:  pbOrder.OrderStatus(newStatusValue),
		})
		if !checkerr(writer, err) {
			return
		}

		log.Printf("Status order successfullt updated")

		http.Redirect(writer, request, "/list/users", http.StatusSeeOther)
	}
}

// PAYMENT PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) paymentHandler(writer http.ResponseWriter, request *http.Request) {
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
	username := session.Values["username"].(string)

	// Retrieve cart to create the order
	cartRes, err := s.clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{Username: username})
	if !checkerr(writer, err) {
		return
	}

	// OrderItems are created depending on CartItems
	var orderItems []*pbOrder.OrderItem
	for _, cartItem := range cartRes.GetCart().GetItems() {

		// Before creating order, check if the items are still in the catalog

		// Retrieve catalog item that was in cart
		itemId := cartItem.GetItemId()
		getRes, err := s.clients.Catalog.GetCatalogItem(request.Context(), &pbCatalog.GetCatalogItemRequest{
			ItemId: itemId,
		})

		// Catalog item removed -> item removed also from cart
		if err != nil || getRes.GetErrorMessage() != "" {
			_, err = s.clients.Cart.RemoveItemFromCart(request.Context(), &pbCart.RemoveItemFromCartRequest{
				Username: username,
				ItemId:   itemId,
			})
			if !checkerr(writer, err) {
				return
			}

			// Redirection to cart page
			http.Redirect(writer, request, "/cart?error=catalog_changed", http.StatusSeeOther)
			return
		}

		// Catalog item price or quantity updated -> Item removed from cart and error message
		catalogItem := getRes.GetItem()
		cartQuantity := cartItem.GetQuantity()
		cartPrice := cartItem.GetPrice()
		// Quantity of the item in the cart greater than the quantity in the catalog -> error
		if cartQuantity > catalogItem.GetQuantityAvailable() || cartPrice != catalogItem.GetPrice() {

			// In case of error item is removed from cart
			_, err = s.clients.Cart.RemoveItemFromCart(request.Context(), &pbCart.RemoveItemFromCartRequest{
				Username: username,
				ItemId:   itemId,
			})
			if !checkerr(writer, err) {
				return
			}

			// Redirection to cart page
			http.Redirect(writer, request, "/cart?error=catalog_changed", http.StatusSeeOther)
			return
		}

		orderItems = append(orderItems, &pbOrder.OrderItem{
			ItemId:   itemId,
			Quantity: cartQuantity,
			Price:    cartPrice,
		})
	}

	// Creation of the Order
	orderRes, err := s.clients.Order.CreateOrder(request.Context(), &pbOrder.CreateOrderRequest{
		UserId:     username,
		OrderItems: orderItems,
	})
	if !checkerr(writer, err) {
		return
	}

	log.Printf("Order successfully created for: %s", username)

	// Retrieve OrderId
	orderIdStr := orderRes.GetOrderId()

	// Retrieve total price order
	priceRes, err := s.clients.Order.GetOrderPrice(request.Context(), &pbOrder.GetOrderPriceRequest{
		OrderId: orderIdStr,
	})
	if !checkerr(writer, err) {
		return
	}

	// Creation of the Payment
	_, err = s.clients.Payment.CreatePayment(request.Context(), &pbPayment.CreatePaymentRequest{
		OrderId: orderIdStr,
		Amount:  math.Trunc(priceRes.GetTotalPrice()*100) / 100,
	})
	log.Printf("Payment successfully created for: %s", username)

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"OrderID": orderIdStr,
		"Amount":  math.Trunc(priceRes.GetTotalPrice()*100) / 100,
	}

	checkerr(writer, s.templates.ExecuteTemplate(writer, "payment.html", templateData))
}

func (s *WebServer) processPaymentHandler(writer http.ResponseWriter, request *http.Request) {
	// Only Post requests are accepted
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	username := session.Values["username"].(string)
	orderId := request.FormValue("order_id")
	amountStr := request.FormValue("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if !checkerr(writer, err) {
		return
	}

	// gRPC call at Payment service
	// Payment status is updated
	_, err = s.clients.Payment.ProcessPayment(request.Context(), &pbPayment.ProcessPaymentRequest{
		OrderId: orderId,
		Amount:  math.Trunc(amount*100) / 100,
	})
	if !checkerr(writer, err) {
		return
	}

	log.Printf("Successfull payment for: %s", username)

	// Update Order Status
	s.clients.Order.UpdateOrderStatus(request.Context(), &pbOrder.UpdateOrderStatusRequest{
		OrderId: orderId,
		Status:  pbOrder.OrderStatus_PROCESSING,
	})
	log.Printf("Processing order for: %s", username)

	// Update catalog changing  the available quantity of acquired items
	// Retrieve cart
	cartRes, err := s.clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{Username: username})
	if !checkerr(writer, err) {
		return
	}

	for _, item := range cartRes.GetCart().GetItems() {
		// Retrieve catalog item that was in cart
		catalogItemRes, err := s.clients.Catalog.GetCatalogItem(request.Context(), &pbCatalog.GetCatalogItemRequest{
			ItemId: item.GetItemId(),
		})
		if !checkerr(writer, err) {
			return
		}

		// Find new available quantity for that catalog item
		currentStock := catalogItemRes.GetItem().GetQuantityAvailable()
		purchasedQty := item.GetQuantity()

		var newStock uint32
		if currentStock >= purchasedQty {
			newStock = currentStock - purchasedQty
		} else {
			newStock = 0
		}

		// Update available quantity for that item in the catalog
		_, err = s.clients.Catalog.UpdateQuantityAvailable(request.Context(), &pbCatalog.UpdateQuantityAvailableRequest{
			ItemId:   item.GetItemId(),
			Quantity: newStock,
		})
		if !checkerr(writer, err) {
			return
		}
	}
	log.Printf("Catalog has been updated successfully after the purchase")

	// Clear cart
	_, err = s.clients.Cart.ClearCart(request.Context(), &pbCart.ClearCartRequest{Username: username})
	if !checkerr(writer, err) {
		return
	}

	log.Printf("Cleared cart for: %s", username)

	checkerr(writer, s.templates.ExecuteTemplate(writer, "process_payment.html", nil))
}

// AUTHETIFICATION PAGE HANDLERS ///////////////////////////////////////////////////////////////

func (s *WebServer) accountHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are permitted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve session
	session, err := s.store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return
	}

	// Check if user is logged
	if loggedIn, ok := session.Values["logged_in"].(bool); !ok || !loggedIn {
		// User not logged -> Redirection to login page
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	// Extraction of user data
	username, _ := session.Values["username"].(string)
	role, _ := session.Values["role"].(string)

	ordersRes, err := s.clients.Order.ListOrdersByUser(request.Context(), &pbOrder.ListOrdersByUserRequest{
		UserId: username,
	})
	if err != nil {
		log.Printf("Error retrieving orders for %s: %v", username, err)
		// In case of error -> empty list
		ordersRes = &pbOrder.ListOrdersByUserResponse{}
	}

	// Preparing user data for HTML template
	templateData := map[string]interface{}{
		"Username": username,
		"Role":     role,
		"Orders":   ordersRes.GetOrders(),
		"Admin":    "ADMIN",
	}

	checkerr(writer, s.templates.ExecuteTemplate(writer, "account.html", templateData))
}

func (s *WebServer) registerHandler(writer http.ResponseWriter, request *http.Request) {
	// User GET request -> register page
	if request.Method == http.MethodGet {
		checkerr(writer, s.templates.ExecuteTemplate(writer, "register.html", nil))
		return
	}

	// User POST request -> user sends credentials
	if request.Method == http.MethodPost {
		username := request.FormValue("username")
		password := request.FormValue("password")
		confirmPassword := request.FormValue("confirm_password")

		// Password validation
		if password != confirmPassword {
			// log.Printf("Failed registration for %s: passwords don't match", username)
			checkerr(writer, s.templates.ExecuteTemplate(writer, "register.html", "Passwords do not match"))
			return
		}

		// gRPC call at Auth service
		_, err := s.clients.Auth.Register(request.Context(), &pbAuth.RegisterRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			// log.Printf("Failed Registration: %v", err)
			checkerr(writer, s.templates.ExecuteTemplate(writer, "register.html", "Username already exists or invalid data"))
			return
		}

		log.Printf("New user correctly registerd: %s", username)

		// Redirection to login page
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
	}
}

func (s *WebServer) loginHandler(writer http.ResponseWriter, request *http.Request) {
	// User GET request -> login page
	if request.Method == http.MethodGet {
		checkerr(writer, s.templates.ExecuteTemplate(writer, "login.html", nil))
		return
	}

	// User POST request -> sends credentials
	if request.Method == http.MethodPost {
		username := request.FormValue("username")
		password := request.FormValue("password")

		// gRPC call at Auth service
		authRes, err := s.clients.Auth.Login(request.Context(), &pbAuth.LoginRequest{
			Username: username,
			Password: password,
		})

		// In case of error -> redirection to login page
		if err != nil {
			checkerr(writer, s.templates.ExecuteTemplate(writer, "login.html", "Username or password are not valid"))
			return
		}

		// Session creation
		session, err := s.store.Get(request, sessionName)
		if !checkerr(writer, err) {
			return
		}

		// Save user data in the session
		session.Values["username"] = authRes.GetUser().GetUsername()
		session.Values["role"] = authRes.GetUser().GetRole().String()
		session.Values["logged_in"] = true

		// Saving session
		err = session.Save(request, writer)
		if !checkerr(writer, err) {
			return
		}

		log.Printf("User %v successfully logged in", session.Values["username"])

		// Redirection to catalog page
		http.Redirect(writer, request, "/cart", http.StatusSeeOther)
	}
}

func (s *WebServer) logoutHandler(writer http.ResponseWriter, request *http.Request) {
	// Retrieve current session
	session, err := s.store.Get(request, sessionName)
	if !checkerr(writer, err) {
		return
	}

	// Set MaxAge=-1 to tell the browser to delete the cookie
	session.Options.MaxAge = -1

	// Save update of the session
	err = session.Save(request, writer)
	if !checkerr(writer, err) {
		return
	}

	log.Printf("User %v successfully logged out", session.Values["username"])

	// Redirecting user to welcome page
	http.Redirect(writer, request, "/welcome", http.StatusSeeOther)
}

func (s *WebServer) changePasswordHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	// User GET request -> change password page
	if request.Method == http.MethodGet {
		checkerr(writer, s.templates.ExecuteTemplate(writer, "change_password.html", nil))
		return
	}

	// User POST request -> user sends credentials
	if request.Method == http.MethodPost {
		username := session.Values["username"]
		currentPassword := request.FormValue("current_password")
		newPassword := request.FormValue("new_password")
		confirmNewPassword := request.FormValue("confirm_new_password")

		// Password validation
		if newPassword != confirmNewPassword {
			log.Printf("Failed registration for %s: passwords don't match", username)
			checkerr(writer, s.templates.ExecuteTemplate(writer, "change_password.html", "Passwords do not match"))
			return
		}

		// gRPC call at Auth service
		_, err := s.clients.Auth.ChangePassword(request.Context(), &pbAuth.ChangePasswordRequest{
			Username:    username.(string),
			OldPassword: currentPassword,
			NewPassword: newPassword,
		})
		if err != nil {
			log.Printf("Failed Changing Passowrds: %v", err)
			checkerr(writer, s.templates.ExecuteTemplate(writer, "change_password.html", "Invalid Password"))
			return
		}

		log.Printf("User %v successfully changed password", session.Values["username"])

		// Redirection to account page
		http.Redirect(writer, request, "/account", http.StatusSeeOther)
	}
}

func (s *WebServer) listAllUsersHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are accepted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// User must be logged
	_, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	users, err := s.clients.Auth.GetAllUsers(request.Context(), &pbAuth.GetAllUsersRequest{})
	if !checkerr(writer, err) {
		return
	}

	// Preparing user data for HTML template
	templateData := map[string]interface{}{
		"Users": users.GetUsers(),
	}

	log.Printf("Users successfully listed")

	checkerr(writer, s.templates.ExecuteTemplate(writer, "list_users.html", templateData))
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
	cookieStore := sessions.NewCookieStore(cookieKey)

	// Web server creation
	server := &WebServer{
		templates: loadTemplates(),
		clients:   clientsRegistry,
		store:     cookieStore,
		manager:   eventsManager,
	}

	mux := http.NewServeMux()

	// Association of paths to correspondent handlers
	mux.HandleFunc("/welcome", server.welcomeHandler)
	mux.HandleFunc("/events", server.manager.HandleEvents)
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
