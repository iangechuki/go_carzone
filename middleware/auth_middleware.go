package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)
var jwtKey = []byte("secret_key")
type Claims struct {
    UserName string `json:"username"`
    jwt.RegisteredClaims
}
func AuthMiddleware(next http.Handler)http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
        authHeader := r.Header.Get("Authorization")
        if authHeader == ""{
            http.Error(w,"Authorization header is missing",http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(authHeader,"Bearer ")
        if tokenString == ""{
            http.Error(w,"Authorization header is missing",http.StatusUnauthorized)
            return
        }
        claims := &Claims{}
        token,err := jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token)(interface{},error){
            return jwtKey,nil
        })
        if err != nil || !token.Valid{
            http.Error(w,"Invalid token",http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), "username", claims.UserName)
        r = r.WithContext(ctx)
        
        next.ServeHTTP(w,r)
    })
	
}