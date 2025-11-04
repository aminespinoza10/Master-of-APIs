package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	_ "swagger/docs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte
var emailKey []byte

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type EmailResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
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

// getEmailHandler godoc
// @Summary Get decrypted email by username
// @Description Returns the decrypted email for the given username
// @Tags users
// @Param username query string true "Username"
// @Success 200 {object} EmailResponse
// @Failure 400 {string} string "Username required"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "DB error" or "Decryption failed"
// @Router /getEmail [get]
func getEmailHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if strings.TrimSpace(username) == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
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

	var encEmail string
	row := conn.QueryRow(ctx, "SELECT email FROM users WHERE username = $1 LIMIT 1", username)
	if err := row.Scan(&encEmail); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	email := ""
	if strings.TrimSpace(encEmail) != "" {
		email, err = decryptEmail(encEmail)
		if err != nil {
			fmt.Println("decryptEmail failed:", err)
			http.Error(w, "Decryption failed", http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(EmailResponse{Username: username, Email: email})
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

	if strings.TrimSpace(cu.Password) == "" {
		http.Error(w, "Password required", http.StatusBadRequest)
		return
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(cu.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	encEmail := ""
	if strings.TrimSpace(cu.Email) != "" {
		encEmail, err = encryptEmail(cu.Email)
		if err != nil {
			http.Error(w, "Failed to encrypt email", http.StatusInternalServerError)
			return
		}
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
	err = conn.QueryRow(ctx, "INSERT INTO users (name, username, password, email) VALUES ($1, $2, $3, $4) RETURNING id", cu.Name, cu.Username, string(hashedPw), encEmail).Scan(&id)
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

	encKeyB64 := os.Getenv("EMAIL_ENC_KEY")
	if encKeyB64 == "" {
		fmt.Println("EMAIL_ENC_KEY environment variable not set!")
		return
	}
	key, err := base64.StdEncoding.DecodeString(encKeyB64)
	if err != nil {
		fmt.Println("EMAIL_ENC_KEY must be base64-encoded:", err)
		return
	}
	if l := len(key); l != 16 && l != 24 && l != 32 {
		fmt.Println("EMAIL_ENC_KEY must decode to 16, 24, or 32 bytes (AES-128/192/256)")
		return
	}
	emailKey = key

	http.Handle("/login", http.HandlerFunc(loginHandler))
	http.Handle("/okCode", jwtMiddleware(http.HandlerFunc(okCodeHandler)))
	http.Handle("/getUsers", jwtMiddleware(http.HandlerFunc(usersHandler)))
	http.Handle("/createUser", jwtMiddleware(http.HandlerFunc(createUserHandler)))
	http.Handle("/getEmail", jwtMiddleware(http.HandlerFunc(getEmailHandler)))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server at :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}

func encryptEmail(plain string) (string, error) {
	block, err := aes.NewCipher(emailKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptEmail(b64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(emailKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
