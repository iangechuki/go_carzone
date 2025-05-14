package login

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iangechuki/go_carzone/models"
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credientials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	valid := (credentials.UserName == "admin") && (credentials.Password == "admin")
	if !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString,err := GenerateToken(credentials.UserName)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Println("Error: ",err)
		return
	}
	response := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject: username,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}