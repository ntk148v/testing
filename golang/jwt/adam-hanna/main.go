// Copyright 2021 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"./randomStrings"
	"./templates"
	"github.com/adam-hanna/jwt-auth/jwt"

	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type randomStringStruct struct {
	Secret string `json:"secret"`
}

var restrictedRoute jwt.Auth

var myUnauthorizedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "I Pitty the fool who is Unauthorized", 401)
	return
})

var refreshSecretHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	csrfSecret := w.Header().Get("X-CSRF-Token")
	claims, err := restrictedRoute.GrabTokenClaims(r)
	log.Println("CSRF secret:", csrfSecret)
	log.Println("Claims", claims)

	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var randomString randomStringStruct
	newString, err := randomstrings.GenerateRandomString(16)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	randomString.Secret = newString
	json.NewEncoder(w).Encode(randomString)
	return
})

var loginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.RenderTemplate(w, "login", &templates.LoginPage{})

	case "POST":
		r.ParseForm()
		log.Println("Username: " + strings.Join(r.Form["username"], ""))

		if strings.Join(r.Form["username"], "") == "testUser" && strings.Join(r.Form["password"], "") == "testPassword" {
			claims := jwt.ClaimsType{}
			claims.CustomClaims = make(map[string]interface{})
			claims.CustomClaims["User"] = r.Form["username"]
			claims.CustomClaims["Role"] = "user"

			err := restrictedRoute.IssueNewTokens(w, &claims)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			var randomString randomStringStruct
			randomString.Secret, err = randomstrings.GenerateRandomString(16)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				return
			}
			json.NewEncoder(w).Encode(randomString)
			return

		}

		http.Error(w, "Unauthorized", 401)
		return

	default:
		http.Error(w, "Method Not Allowed", 405)
	}
})

var logoutHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := restrictedRoute.NullifyTokens(w, r)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method Not Allowed", 405)
	}
})

func main() {
	authErr := jwt.New(&restrictedRoute, jwt.Options{
		SigningMethodString:   "RS256",
		PrivateKeyLocation:    "keys/app.rsa",     // `$ openssl genrsa -out app.rsa 2048`
		PublicKeyLocation:     "keys/app.rsa.pub", // `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`
		RefreshTokenValidTime: 10 * time.Minute,
		AuthTokenValidTime:    5 * time.Minute,
		Debug:                 true,
		BearerTokens:          true,
	})
	if authErr != nil {
		log.Println("Error initializing the JWT's!")
		log.Fatal(authErr)
	}

	restrictedRoute.SetUnauthorizedHandler(myUnauthorizedHandler)

	http.HandleFunc("/", loginHandler)
	http.Handle("/refreshSecret", restrictedRoute.Handler(refreshSecretHandler))
	http.Handle("/logout", restrictedRoute.Handler(logoutHandler))

	log.Println("Listening on localhost:3000")
	http.ListenAndServe("127.0.0.1:3000", nil)
}
