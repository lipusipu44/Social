package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string
	audience string
	issuer   string
}

//NewJWTAuthenticator
/*
This to be called in main.go class to create a JWTAuthenticator,
same flow like we did with NewStorage() to create a store instance
which has all concrete implementations of UserStorage, PostStorage etc..
*/
func NewJWTAuthenticator(secret string, audience string, issuer string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		audience: audience,
		issuer:   issuer,
	}
}

//GenerateToken
/*
Same flow like in storage UserStorage/PostStorage class, take the struct to call
the method and try to implement the authenticator interface.

so in api class of app struct we are going to use a field auth => auth.Authenticator
but which authenticator to be used can be decided in main.go, I learned the importance of
interface and how to change it in main.go as per usage is awesome

creates a token using jwt methods, I dont know much about it just trying to use it
*/
func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//ValidateToken
/*
Now actual implementation of valid token to be written here,
no idea to check in ChatGPT bit later
*/
func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(j.audience),
		jwt.WithIssuer(j.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
}
