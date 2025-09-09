package main

import (
	"fmt"
	"net/http"
)

func okCodeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Everything is awesome!")
}

func continueCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusContinue)
	fmt.Fprintln(w, "Continue processing...")
}

func movedPemanentlyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMovedPermanently)
	fmt.Fprintln(w, "This resource has been moved permanently.")
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request. Please check your input.")
}

func forbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintln(w, "Access forbidden. You don't have permission to access this resource.")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "Resource not found.")
}

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

	fmt.Println("Starting server at :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
