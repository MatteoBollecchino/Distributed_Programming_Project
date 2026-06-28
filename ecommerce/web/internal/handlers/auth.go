package handlers

import (
	"log"
	"net/http"

	pbAuth "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/auth"
	pbOrder "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/order"
)

func (s *ServerDependencies) AccountHandler(writer http.ResponseWriter, request *http.Request) {
	// Only GET requests are permitted
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve session
	session, err := s.Store.Get(request, sessionName)
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

	ordersRes, err := s.Clients.Order.ListOrdersByUser(request.Context(), &pbOrder.ListOrdersByUserRequest{
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

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "account.html", templateData))
}

func (s *ServerDependencies) RegisterHandler(writer http.ResponseWriter, request *http.Request) {
	// User GET request -> register page
	if request.Method == http.MethodGet {
		checkerr(writer, s.Templates.ExecuteTemplate(writer, "register.html", nil))
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
			checkerr(writer, s.Templates.ExecuteTemplate(writer, "register.html", "Passwords do not match"))
			return
		}

		// gRPC call at Auth service
		_, err := s.Clients.Auth.Register(request.Context(), &pbAuth.RegisterRequest{
			Username: username,
			Password: password,
		})

		if err != nil {
			// log.Printf("Failed Registration: %v", err)
			checkerr(writer, s.Templates.ExecuteTemplate(writer, "register.html", "Username already exists or invalid data"))
			return
		}

		log.Printf("New user correctly registerd: %s", username)

		// Redirection to login page
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
	}
}

func (s *ServerDependencies) LoginHandler(writer http.ResponseWriter, request *http.Request) {
	// User GET request -> login page
	if request.Method == http.MethodGet {
		checkerr(writer, s.Templates.ExecuteTemplate(writer, "login.html", nil))
		return
	}

	// User POST request -> sends credentials
	if request.Method == http.MethodPost {
		username := request.FormValue("username")
		password := request.FormValue("password")

		// gRPC call at Auth service
		authRes, err := s.Clients.Auth.Login(request.Context(), &pbAuth.LoginRequest{
			Username: username,
			Password: password,
		})

		// In case of error -> redirection to login page
		if err != nil {
			checkerr(writer, s.Templates.ExecuteTemplate(writer, "login.html", "Username or password are not valid"))
			return
		}

		// Session creation
		session, err := s.Store.Get(request, sessionName)
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

func (s *ServerDependencies) LogoutHandler(writer http.ResponseWriter, request *http.Request) {
	// Retrieve current session
	session, err := s.Store.Get(request, sessionName)
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

func (s *ServerDependencies) ChangePasswordHandler(writer http.ResponseWriter, request *http.Request) {
	// User must be logged
	session, ok := checkIfUserIsLogged(s, request, writer)
	if !ok {
		return
	}

	// User GET request -> change password page
	if request.Method == http.MethodGet {
		checkerr(writer, s.Templates.ExecuteTemplate(writer, "change_password.html", nil))
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
			checkerr(writer, s.Templates.ExecuteTemplate(writer, "change_password.html", "Passwords do not match"))
			return
		}

		// gRPC call at Auth service
		_, err := s.Clients.Auth.ChangePassword(request.Context(), &pbAuth.ChangePasswordRequest{
			Username:    username.(string),
			OldPassword: currentPassword,
			NewPassword: newPassword,
		})
		if err != nil {
			log.Printf("Failed Changing Passowrds: %v", err)
			checkerr(writer, s.Templates.ExecuteTemplate(writer, "change_password.html", "Invalid Password"))
			return
		}

		log.Printf("User %v successfully changed password", session.Values["username"])

		// Redirection to account page
		http.Redirect(writer, request, "/account", http.StatusSeeOther)
	}
}

func (s *ServerDependencies) ListAllUsersHandler(writer http.ResponseWriter, request *http.Request) {
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

	users, err := s.Clients.Auth.GetAllUsers(request.Context(), &pbAuth.GetAllUsersRequest{})
	if !checkerr(writer, err) {
		return
	}

	// Preparing user data for HTML template
	templateData := map[string]interface{}{
		"Users": users.GetUsers(),
	}

	log.Printf("Users successfully listed")

	checkerr(writer, s.Templates.ExecuteTemplate(writer, "list_users.html", templateData))
}
