package main

import (
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

func checkIfUserIsLogged(s *WebServer, request *http.Request, writer http.ResponseWriter) (*sessions.Session, bool) {
	session, err := s.store.Get(request, sessionName)
	checkerr(writer, err)
	if loggedIn, ok := session.Values["logged_in"].(bool); !ok || !loggedIn {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return nil, false
	}
	return session, true
}

// WELCOME PAGE HANDLER ///////////////////////////////////////////////////////////////

func (s *WebServer) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	checkerr(writer, s.templates.ExecuteTemplate(writer, "welcome_page.html", nil))
}

// CATALOG PAGE HANDLER ///////////////////////////////////////////////////////////////

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

// CART PAGE HANDLERS ///////////////////////////////////////////////////////////////

func (s *WebServer) cartHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are accepted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve session
	session, err := s.store.Get(request, sessionName)
	checkerr(writer, err)

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
	var totalPrice float64
	for _, item := range cartRes.GetCart().GetItems() {
		totalPrice += float64(item.GetQuantity()) * float64(item.GetPrice())
	}

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"Items":      cartRes.GetCart().GetItems(),
		"TotalPrice": math.Trunc(totalPrice*100) / 100,
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
	if err != nil {
		log.Printf("Error in price conversion '%s': %v", priceStr, err)
		http.Error(writer, "Price not valid", http.StatusBadRequest)
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
	if err != nil {
		log.Printf("Impossible to retrieve shopping cart for %s: %v", username, err)
		http.Error(writer, "Error in retrieving user cart", http.StatusInternalServerError)
		return
	}

	// Calculate Totale price
	var totalPrice float64
	for _, item := range cartRes.GetCart().GetItems() {
		totalPrice += float64(item.GetQuantity()) * float64(item.GetPrice())
	}

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"Items":      cartRes.GetCart().GetItems(),
		"TotalPrice": math.Trunc(totalPrice*100) / 100,
	}

	checkerr(writer, s.templates.ExecuteTemplate(writer, "order.html", templateData))
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
	checkerr(writer, err)

	// OrderItems are created depending on CartItems
	var orderItems []*pbOrder.OrderItem
	for _, cartItem := range cartRes.GetCart().GetItems() {
		orderItems = append(orderItems, &pbOrder.OrderItem{
			ItemId:   cartItem.GetItemId(),
			Quantity: cartItem.GetQuantity(),
			Price:    cartItem.GetPrice(),
		})
	}

	// Creation of the Order
	orderRes, err := s.clients.Order.CreateOrder(request.Context(), &pbOrder.CreateOrderRequest{
		UserId:     username,
		OrderItems: orderItems,
	})
	if err != nil {
		log.Printf("Error creating order: %v", err)
		http.Error(writer, "Error creating order", http.StatusInternalServerError)
		return
	}
	log.Printf("Order successfully created for: %s", username)

	// Retrieve OrderId
	orderIdStr := orderRes.GetOrderId()

	// Retrieve total price order
	priceRes, err := s.clients.Order.GetOrderPrice(request.Context(), &pbOrder.GetOrderPriceRequest{
		OrderId: orderIdStr,
	})
	if err != nil {
		log.Printf("Error retrieving price for order %s: %v", orderIdStr, err)
		http.Error(writer, "Error retrieving price", http.StatusInternalServerError)
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
	checkerr(writer, err)

	// gRPC call at Payment service
	// Payment status is updated
	_, err = s.clients.Payment.ProcessPayment(request.Context(), &pbPayment.ProcessPaymentRequest{
		OrderId: orderId,
		Amount:  math.Trunc(amount*100) / 100,
	})
	if err != nil {
		log.Printf("Failed payment for order %s: %v", orderId, err)
		http.Error(writer, "Denied Transaction", http.StatusPaymentRequired)
		return
	}
	log.Printf("Successfull payment for: %s", username)

	// Update Order Status
	s.clients.Order.UpdateOrderStatus(request.Context(), &pbOrder.UpdateOrderStatusRequest{
		OrderId: orderId,
		Status:  pbOrder.OrderStatus_PROCESSING,
	})
	log.Printf("Processing order for: %s", username)

	// Update catalog updating  the available quantity of acquired items
	// 1. Recupera il carrello in modo sicuro
	cartRes, err := s.clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{Username: username})
	if err != nil {
		log.Printf("Failed to get cart for catalog update: %v", err)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return // FONDAMENTALE: blocca la funzione se c'è un errore
	}

	// 2. Ciclo sui prodotti del carrello
	for _, item := range cartRes.GetCart().GetItems() {
		// Recupera l'articolo dal catalogo per sapere la quantità attuale
		catalogItemRes, err := s.clients.Catalog.GetCatalogItem(request.Context(), &pbCatalog.GetCatalogItemRequest{
			ItemId: item.GetItemId(),
		})
		if err != nil {
			log.Printf("Failed to get catalog item %s: %v", item.GetItemId(), err)
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return // Evita il panic sulla riga successiva
		}

		// Calcola la nuova quantità disponibile
		currentStock := catalogItemRes.GetItem().GetQuantityAvailable()
		purchasedQty := item.GetQuantity()

		// Controllo di sicurezza per evitare che la quantità diventi negativa
		var newStock uint32
		if currentStock >= purchasedQty {
			newStock = currentStock - purchasedQty
		} else {
			newStock = 0 // O gestisci un errore di "prodotto esaurito" se necessario
		}

		// Aggiorna la quantità nel microservizio del catalogo
		_, err = s.clients.Catalog.UpdateQuantityAvailable(request.Context(), &pbCatalog.UpdateQuantityAvailableRequest{
			ItemId:   item.GetItemId(),
			Quantity: newStock,
		})
		if err != nil {
			log.Printf("Failed to update stock for item %s: %v", item.GetItemId(), err)
			http.Error(writer, "Failed to update catalog stock", http.StatusInternalServerError)
			return
		}
	}
	log.Printf("Catalog has been updated successfully after the purchase")

	// 3. Svuota il carrello dopo che tutto il resto è andato a buon fine
	_, err = s.clients.Cart.ClearCart(request.Context(), &pbCart.ClearCartRequest{Username: username})
	if err != nil {
		log.Printf("Warning: Failed to clear cart for %s: %v", username, err)
	} else {
		log.Printf("Cleared cart for: %s", username)
	}

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
	checkerr(writer, err)

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
			log.Printf("Failed registration for %s: passwords don't match", username)
			checkerr(writer, s.templates.ExecuteTemplate(writer, "register.html", "Passwords do not match"))
			return
		}

		// gRPC call at Auth service
		_, err := s.clients.Auth.Register(request.Context(), &pbAuth.RegisterRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			log.Printf("Failed Registration: %v", err)
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
			log.Printf("Failed Login: %v", err)
			checkerr(writer, s.templates.ExecuteTemplate(writer, "login.html", "Username or password are not valid"))
			//http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}

		// Session creation
		session, err := s.store.Get(request, sessionName)
		checkerr(writer, err)

		// Save user data in the session
		session.Values["username"] = authRes.GetUser().Username
		session.Values["role"] = authRes.GetUser().Role.String()
		session.Values["logged_in"] = true

		// Saving session
		if err = session.Save(request, writer); err != nil {
			http.Error(writer, "Saving Session Error", http.StatusInternalServerError)
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
	checkerr(writer, err)

	// Set MaxAge=-1 to tell the browser to delete the cookie
	session.Options.MaxAge = -1

	// Save update of the session
	if err := session.Save(request, writer); err != nil {
		http.Error(writer, "Errore during logout", http.StatusInternalServerError)
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

	// User GET request -> register page
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
	// DA FINIRE
}

// MAIN ///////////////////////////////////////////////////////////////

func main() {
	// Initialization gRPC clients
	clientsRegistry, err := clients.InitClients()
	if err != nil {
		log.Fatal(err)
	}
	defer clientsRegistry.Close()

	// Cookies creation
	cookieStore := sessions.NewCookieStore(cookieKey)

	// Web server creation
	server := &WebServer{
		templates: loadTemplates(),
		clients:   clientsRegistry,
		store:     cookieStore,
	}

	mux := http.NewServeMux()

	// Association of paths to correspondent handlers
	mux.HandleFunc("/welcome", server.welcomeHandler)
	mux.HandleFunc("/catalog", server.productCatalogHandler)
	mux.HandleFunc("/cart", server.cartHandler)
	mux.HandleFunc("/cart/add", server.addToCartHandler)
	mux.HandleFunc("/cart/remove", server.removeFromCartHandler)
	mux.HandleFunc("/order", server.orderHandler)
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
