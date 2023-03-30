package main

import (
	"fmt"
	"net/http"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Get the auth cookie from the request
	cookie, err := r.Cookie("auth")
	if err != nil {
		// Cookie not set
		http.Error(w, "You are not admin", http.StatusUnauthorized)
		return
	}

	// Check if the cookie value is "admin"
	if cookie.Value != "admin" {
		http.Error(w, "You are not admin", http.StatusUnauthorized)
		return
	}

	// Set the "user" cookie for the new user
	userCookie := &http.Cookie{
		Name:  "user",
		Value: "user",
	}
	http.SetCookie(w, userCookie)

	// Display the "/admin" content
	fmt.Fprintf(w, "Welcome to home")
}

func main() {
	http.HandleFunc("/admin", adminHandler)
	http.ListenAndServe(":8000", nil)
}
