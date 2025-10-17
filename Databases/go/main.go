package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	_ "swagger/docs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

var jwtSecret []byte

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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

// usersHandler godoc
// @Summary Get all users
// @Description Returns a list of users from the database
// @Tags users
// @Success 200 {array} User
// @Failure 500 {string} string "Failed to connect to database" or "Query failed" or "Row scan failed"
// @Router /getUsers [get]
func usersHandler(w http.ResponseWriter, r *http.Request) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		http.Error(w, "DATABASE_URL not set", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT id, name, username FROM users")
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Username); err != nil {
			http.Error(w, "Row scan failed", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// createUserHandler godoc
// @Summary Create a new user
// @Description Create a new user and return the created record
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUser true "New user"
// @Success 201 {object} User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "DB error"
// @Router /createUser [post]
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var cu CreateUser
	if err := json.NewDecoder(r.Body).Decode(&cu); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		http.Error(w, "DATABASE_URL not set", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)

	var id int
	err = conn.QueryRow(ctx, "INSERT INTO users (name, username, password) VALUES ($1, $2, $3) RETURNING id", cu.Name, cu.Username, cu.Password).Scan(&id)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	user := User{ID: id, Name: cu.Name, Username: cu.Username}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
	http.Handle("/getUsers", jwtMiddleware(http.HandlerFunc(usersHandler)))
	http.Handle("/createUser", jwtMiddleware(http.HandlerFunc(createUserHandler)))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server at :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
