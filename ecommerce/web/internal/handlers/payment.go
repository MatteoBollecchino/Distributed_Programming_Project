package handlers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	pbCart "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/cart"
	pbCatalog "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	pbOrder "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
	pbPayment "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/payment"
)

func (s *ServerDependencies) PaymentHandler(writer http.ResponseWriter, request *http.Request) {
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
	cartRes, err := s.Clients.Cart.GetCart(request.Context(), &pbCart.GetCartRequest{
		Username: username,
	})
	if !checkerr(writer, err) {
		return
	}

	// OrderItems are created depending on CartItems
	var orderItems []*pbOrder.OrderItem
	for _, cartItem := range cartRes.GetCart().GetItems() {

		// Before creating order, check if the items are still in the catalog

		// Retrieve catalog item that was in cart
		itemId := cartItem.GetItemId()
		getRes, err := s.Clients.Catalog.GetCatalogItem(request.Context(), &pbCatalog.GetCatalogItemRequest{
			ItemId: itemId,
		})

		// Catalog item removed -> item removed also from cart
		if err != nil || getRes.GetErrorMessage() != "" {
			_, err = s.Clients.Cart.RemoveItemFromCart(request.Context(), &pbCart.RemoveItemFromCartRequest{
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
			_, err = s.Clients.Cart.RemoveItemFromCart(request.Context(), &pbCart.RemoveItemFromCartRequest{
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
	orderRes, err := s.Clients.Order.CreateOrder(request.Context(), &pbOrder.CreateOrderRequest{
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
	priceRes, err := s.Clients.Order.GetOrderPrice(request.Context(), &pbOrder.GetOrderPriceRequest{
		OrderId: orderIdStr,
	})
	if !checkerr(writer, err) {
		return
	}

	// Creation of the Payment
	_, err = s.Clients.Payment.CreatePayment(request.Context(), &pbPayment.CreatePaymentRequest{
		OrderId: orderIdStr,
		Amount:  math.Trunc(priceRes.GetTotalPrice()*100) / 100,
	})
	log.Printf("Payment successfully created for: %s", username)

	// Mapping data for HTML file
	templateData := map[string]interface{}{
		"OrderID": orderIdStr,
		"Amount":  math.Trunc(priceRes.GetTotalPrice()*100) / 100,
	}

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "payment.html", templateData))
}

func (s *ServerDependencies) ProcessPaymentHandler(writer http.ResponseWriter, request *http.Request) {
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
	_, err = s.Clients.Payment.ProcessPayment(request.Context(), &pbPayment.ProcessPaymentRequest{
		OrderId: orderId,
		Amount:  math.Trunc(amount*100) / 100,
	})
	if !checkerr(writer, err) {
		return
	}

	// Verify payment status
	statusRes, err := s.Clients.Payment.GetPaymentStatus(request.Context(), &pbPayment.GetPaymentStatusRequest{
		OrderId: orderId,
	})
	if !checkerr(writer, err) || statusRes.GetStatus() != pbPayment.PaymentStatus_PAID {
		http.Redirect(writer, request, "/cart?error=payment_failed", http.StatusSeeOther)
		return
	}

	log.Printf("Successfull payment for: %s", username)

	// Update Order Status
	s.Clients.Order.UpdateOrderStatus(request.Context(), &pbOrder.UpdateOrderStatusRequest{
		OrderId: orderId,
		Status:  pbOrder.OrderStatus_PROCESSING,
	})
	log.Printf("Processing order for: %s", username)

	// Update catalog changing the available quantity of acquired items

	// Retrieve order
	orderRes, err := s.Clients.Order.GetOrder(request.Context(), &pbOrder.GetOrderRequest{OrderId: orderId})
	if !checkerr(writer, err) {
		return
	}

	for _, item := range orderRes.GetOrder().GetItems() {
		// Retrieve catalog item that was in the cart
		catalogItemRes, err := s.Clients.Catalog.GetCatalogItem(request.Context(), &pbCatalog.GetCatalogItemRequest{
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
		_, err = s.Clients.Catalog.UpdateQuantityAvailable(request.Context(), &pbCatalog.UpdateQuantityAvailableRequest{
			ItemId:   item.GetItemId(),
			Quantity: newStock,
		})
		if !checkerr(writer, err) {
			return
		}
	}
	log.Printf("Catalog has been updated successfully after the purchase")

	// Clear cart
	_, err = s.Clients.Cart.ClearCart(request.Context(), &pbCart.ClearCartRequest{Username: username})
	if !checkerr(writer, err) {
		return
	}

	log.Printf("Cleared cart for: %s", username)

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "process_payment.html", nil))
}
