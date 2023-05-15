// Declare the package main
package main

// Import the necessary libraries
import (
	"encoding/base64" // to encode and decode strings to/from base64
	"log"             // to log errors
	"net/http"        // to handle HTTP requests
)

// Constants that will be used throughout the application
const (
	Flag        = "CTF{4514n_L0x}"                               // the flag to display
	SecretUser  = "01110101 01110011 01100101 01110010"          // default user role
	SecretAdmin = "01100001 01100100 01101101 01101001 01101110" // admin role
	CookieName  = "Role"                                         // the name of the cookie to check
)

// Server struct that holds the state of the server
type Server struct {
	secretUser  string // base64 encoded "user"
	secretAdmin string // base64 encoded "admin"
	flag        string // the flag to display
}

// Constructor function for creating a new Server instance
func NewServer() *Server {
	return &Server{
		secretUser:  base64.StdEncoding.EncodeToString([]byte(SecretUser)),
		secretAdmin: base64.StdEncoding.EncodeToString([]byte(SecretAdmin)),
		flag:        Flag,
	}
}

// Method of the Server struct that handles requests to the index route ("/")
func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to read the "Role" cookie from the request
		cookie, err := r.Cookie(CookieName)

		// If the cookie does not exist or its value is not the base64 encoding of "admin",
		// a new cookie is created with the value being the base64 encoding of "user",
		// and an "Access Denied" message is displayed.
		if err != nil || cookie.Value != s.secretAdmin {
			cookie := &http.Cookie{
				Name:  CookieName,
				Value: s.secretUser,
			}
			http.SetCookie(w, cookie)
			http.Error(w, "Access Denied", http.StatusUnauthorized)
			return
		}
		// If the cookie exists and its value is the base64 encoding of "admin",
		// the flag is displayed.
		w.Write([]byte(s.flag))
	}
}

// Method of the Server struct that starts the server
func (s *Server) Run(addr string) {
	// Handle requests to the index route ("/") with the handleIndex method
	http.HandleFunc("/", s.handleIndex())

	// Start the server and log any errors that occur
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Main function that creates a new Server instance and runs it
func main() {
	server := NewServer()
	server.Run(":1337")
}
