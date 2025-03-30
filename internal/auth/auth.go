package auth

import "github.com/golang-jwt/jwt/v5"

//Authenticator
/*
The concept is same as store interface we have created for SQL,
this interface will work as an Authentication base, now its JWT stateless
tomorrow it can be anything else who can implement these methods.

This InterFace is going to be called from app like store has been called,
like in store.go NewStorage() is used to create concrete instance in main.go, similarly for this
also it will be created and to be passed in app instance in main.go, so actual implementation can be seen
*/
type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
