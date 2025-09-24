package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	_ "swagger/docs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

var jwtSecret []byte

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// loginHandler godoc
// @Summary Generate JWT token
// @Description Returns a JWT token for a given username
// @Tags auth
// @Param username query string true "Username"
// @Success 200 {string} string "JWT token"
// @Failure 400 {string} string "Username required"
// @Failure 500 {string} string "Could not generate token"
// @Router /login [get]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}
	token, err := generateJWT(username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(token))
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtSecret)
}

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
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found or failed to load")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("JWT_SECRET environment variable not set!")
		return
	}
	jwtSecret = []byte(secret)

	http.Handle("/login", http.HandlerFunc(loginHandler))
	http.Handle("/okCode", jwtMiddleware(http.HandlerFunc(okCodeHandler)))
	http.Handle("/continueCode", jwtMiddleware(http.HandlerFunc(continueCodeHandler)))
	http.Handle("/movedPermanently", jwtMiddleware(http.HandlerFunc(movedPemanentlyHandler)))
	http.Handle("/badRequest", jwtMiddleware(http.HandlerFunc(badRequestHandler)))
	http.Handle("/forbidden", jwtMiddleware(http.HandlerFunc(forbiddenHandler)))
	http.Handle("/notFound", jwtMiddleware(http.HandlerFunc(notFoundHandler)))
	http.Handle("/proxyRequired", jwtMiddleware(http.HandlerFunc(proxyRequiredHandler)))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server at :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
