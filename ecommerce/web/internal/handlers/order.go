package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	pbCart "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	pbOrder "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
)

func (s *ServerDependencies) OrderHandler(writer http.ResponseWriter, request *http.Request) {
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
	cartRes, err := s.Clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{
		Username: username,
	})
	if !checkerr(writer, err) {
		return
	}

	// Calculate total price
	totalPriceRes, err := s.Clients.Cart.CalculateTotalPrice(request.Context(), &pbCart.CalculateTotalPriceRequest{
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

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "order.html", templateData))
}

func (s *ServerDependencies) UserOrdersHandler(writer http.ResponseWriter, request *http.Request) {
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
		orderRes, err := s.Clients.Order.ListOrdersByUser(request.Context(), &pbOrder.ListOrdersByUserRequest{
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

		checkerr(writer, s.Templates.ExecuteTemplate(writer, "user_orders.html", templateData))
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
		_, err = s.Clients.Order.UpdateOrderStatus(request.Context(), &pbOrder.UpdateOrderStatusRequest{
			OrderId: orderId,
			Status:  pbOrder.OrderStatus(newStatusValue),
		})
		if !checkerr(writer, err) {
			return
		}

		log.Printf("Status order successfully updated")

		http.Redirect(writer, request, "/list/users", http.StatusSeeOther)
	}
}
