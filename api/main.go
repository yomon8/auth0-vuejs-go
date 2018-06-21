package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	Auth0Domain        = "YOUR-DOMAIN.auth0.com"
	Auth0Audience      = "http://localhost:50000"
	Auth0RequiredScope = "read:messages"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func checkScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)

	claims, _ := token.Claims.(*CustomClaims)

	hasScope := false
	result := strings.Split(claims.Scope, " ")
	for i := range result {
		if result[i] == scope {
			hasScope = true
		}
	}
	return hasScope
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + Auth0Domain + "/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func public(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello public!\n"))
}

func idValidate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello jwt validated!\n"))
}

func scopeValidate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello jwt and scope validated!\n"))
}

func getValidateJwtMiddleware() func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(Auth0Audience, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "https://" + Auth0Domain + "/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				return token, err
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware.HandlerWithNext
}

func validateScopeMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]
	hasScope := checkScope(Auth0RequiredScope, token)
	if hasScope {
		next(rw, r)
	} else {
		http.Error(rw, "Not Authorized", 401)
	}
}

func main() {
	r := mux.NewRouter()

	r.Handle("/public", http.HandlerFunc(public))
	r.Handle("/private", negroni.New(
		negroni.HandlerFunc(getValidateJwtMiddleware()),
		negroni.Wrap(http.HandlerFunc(idValidate))))
	r.Handle("/private-scoped", negroni.New(
		negroni.HandlerFunc(getValidateJwtMiddleware()),
		negroni.HandlerFunc(validateScopeMiddleware),
		negroni.Wrap(http.HandlerFunc(scopeValidate))))

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization"})
	log.Fatal(http.ListenAndServe(":50000", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)))
}
