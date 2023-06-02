// Declare the package main
package main

// Import the necessary libraries
import (
	"encoding/base64" // to encode and decode strings to/from base64
	"html/template"
	"log"      // to log errors
	"net/http" // to handle HTTP requests
	"path/filepath"
)

// Constants that will be used throughout the application
const (
	Flag        = "aituctf{Ch4n91n9_c00k13_v3ry_345y!}"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    // the flag to display
	SecretUser  = "00110110 00110101 00100000 00110110 00110101 00100000 00110011 00110001 00100000 00110011 00110001 00100000 00110110 00110011 00100000 00110110 00110010 00100000 00110110 00110010 00100000 00110011 00110001 00100000 00110011 00111001 00100000 00110011 00110000 00100000 00110011 00110101 00100000 00110011 00110010 00100000 00110110 00110101 00100000 00110011 00110100 00100000 00110011 00110000 00100000 00110110 00110010 00100000 00110011 00110000 00100000 00110011 00110111 00100000 00110110 00110001 00100000 00110110 00110001 00100000 00110110 00110011 00100000 00110011 00110000 00100000 00110110 00110011 00100000 00110110 00110001 00100000 00110011 00110000 00100000 00110011 00110110 00100000 00110011 00110000 00100000 00110110 00110011 00100000 00110011 00110010 00100000 00110011 00110011 00100000 00110110 00110101 00100000 00110110 00110101" // default user role
	SecretAdmin = "00110011 00110010 00100000 00110011 00110001 00100000 00110011 00110010 00100000 00110011 00110011 00100000 00110011 00110010 00100000 00110110 00110110 00100000 00110011 00110010 00100000 00110011 00111001 00100000 00110011 00110111 00100000 00110110 00110001 00100000 00110011 00110101 00100000 00110011 00110111 00100000 00110110 00110001 00100000 00110011 00110101 00100000 00110110 00110001 00100000 00110011 00110111 00100000 00110011 00110100 00100000 00110011 00110011 00100000 00110011 00111000 00100000 00110011 00111001 00100000 00110011 00110100 00100000 00110110 00110001 00100000 00110011 00110000 00100000 00110110 00110101 00100000 00110011 00110100 00100000 00110110 00110001 00100000 00110011 00111000 00100000 00110011 00110000 00100000 00110011 00110001 00100000 00110110 00110110 00100000 00110110 00110011 00100000 00110011 00110011" // admin role
	CookieName  = "Role"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   // the name of the cookie to check
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
		message := "Access Denied"
		statusCode := http.StatusUnauthorized

		if err != nil || cookie.Value != s.secretAdmin {
			cookie := &http.Cookie{
				Name:  CookieName,
				Value: s.secretUser,
			}
			http.SetCookie(w, cookie)
		} else {
			message = s.flag
			statusCode = http.StatusOK
		}

		templatesDir := filepath.Join("templates", "index.html")
		tmpl, err := template.ParseFiles(templatesDir)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(statusCode)
		data := struct {
			Message string
		}{
			Message: message,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// Method of the Server struct that starts the server
func (s *Server) Run(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex())

	// Wrap the mux with the logRequest middleware
	loggedRouter := logRequest(mux)

	// Start the server and log any errors that occur
	log.Fatal(http.ListenAndServe(addr, loggedRouter))
}
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Main function that creates a new Server instance and runs it
func main() {

	server := NewServer()
	server.Run(":1337")
}
