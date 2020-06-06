package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/casbin/casbin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	bearerFormat    string = "Bearer %s"
	tokenExpiration        = 60
)

func IssueToken(w http.ResponseWriter, req *http.Request) {
	user, pass, _ := req.BasicAuth()
	if user != "test" || pass != "test" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expTime := time.Now().Add(tokenExpiration * time.Minute)
	claims := jwt.MapClaims{}
	claims["exp"] = expTime.Unix()
	claims["user"] = user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Authorization", fmt.Sprintf(bearerFormat, tokenString))
}

func AuthenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authzHeader := req.Header.Get("Authorization")
		if authzHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Absent header")
			return
		}
		bearerToken := strings.Split(authzHeader, " ")
		if len(bearerToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Invalid bearer token")
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			log.Print("Got token ", token)
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}
		var user string
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user = claims["user"].(string)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "sad")
			return
		}
		ctx := context.WithValue(req.Context(), "user", user)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func AuthorMiddleware(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			user := ctx.Value("user")
			log.Println(user)
			res, err := e.Enforce(user, req.URL.Path, req.Method)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !res {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(fn)
	}
}

func Ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "PONG PING")
}

func main() {
	var port = flag.Int("port", 3000, "port to listen on")
	flag.Parse()
	policyEngine, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		log.Fatal("Unable to create casbin policy engine: ", err)
	}
	router := mux.NewRouter()
	pubRouter := router.PathPrefix("/public").Subrouter()
	pubRouter.HandleFunc("/token", IssueToken).Methods("POST")
	privRouter := router.PathPrefix("/private").Subrouter()
	privRouter.Use(AuthenMiddleware, AuthorMiddleware(policyEngine))
	privRouter.HandleFunc("/secret", Ping).Methods("GET")
	log.Println("Starting the server...")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), router))
}
