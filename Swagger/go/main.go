package main

import (
	"fmt"
	"net/http"
	_ "swagger/docs"

	"github.com/swaggo/http-swagger"
)

// okCodeHandler godoc
// @Summary Returns OK status
// @Description Responds with HTTP 200 and a message
// @Tags codes
// @Success 200 {string} string "Everything is awesome!"
// @Router /okCode [get]
func okCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Everything is awesome!")
}

// continueCodeHandler godoc
// @Summary Returns Continue status
// @Description Responds with HTTP 100 and a message
// @Tags codes
// @Success 100 {string} string "Continue processing..."
// @Router /continueCode [get]
func continueCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusContinue)
	fmt.Fprintln(w, "Continue processing...")
}

// movedPemanentlyHandler godoc
// @Summary Returns Moved Permanently status
// @Description Responds with HTTP 301 and a message
// @Tags codes
// @Success 301 {string} string "This resource has been moved permanently."
// @Router /movedPermanently [get]
func movedPemanentlyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMovedPermanently)
	fmt.Fprintln(w, "This resource has been moved permanently.")
}

// badRequestHandler godoc
// @Summary Returns Bad Request status
// @Description Responds with HTTP 400 and a message
// @Tags codes
// @Success 400 {string} string "Bad request. Please check your input."
// @Router /badRequest [get]
func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request. Please check your input.")
}

// forbiddenHandler godoc
// @Summary Returns Forbidden status
// @Description Responds with HTTP 403 and a message
// @Tags codes
// @Success 403 {string} string "Access forbidden. You don't have permission to access this resource."
// @Router /forbidden [get]
func forbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintln(w, "Access forbidden. You don't have permission to access this resource.")
}

// notFoundHandler godoc
// @Summary Returns Not Found status
// @Description Responds with HTTP 404 and a message
// @Tags codes
// @Success 404 {string} string "Resource not found."
// @Router /notFound [get]
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "Resource not found.")
}

// proxyRequiredHandler godoc
// @Summary Returns Proxy Authentication Required status
// @Description Responds with HTTP 407 and a message
// @Tags codes
// @Success 407 {string} string "Proxy authentication required."
// @Router /proxyRequired [get]
func proxyRequiredHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusProxyAuthRequired)
	fmt.Fprintln(w, "Proxy authentication required.")
}

func main() {
	http.HandleFunc("/okCode", okCodeHandler)
	http.HandleFunc("/continueCode", continueCodeHandler)
	http.HandleFunc("/movedPermanently", movedPemanentlyHandler)
	http.HandleFunc("/badRequest", badRequestHandler)
	http.HandleFunc("/forbidden", forbiddenHandler)
	http.HandleFunc("/notFound", notFoundHandler)
	http.HandleFunc("/proxyRequired", proxyRequiredHandler)

	// Serve Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server at :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
