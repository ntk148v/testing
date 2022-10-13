// https://cloud.google.com/translate/docs/basic/translating-text#translating_text
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const GoogleTransApiUrl = "https://translation.googleapis.com/language/translate/v2"

type ResponseBody struct {
	Data struct {
		Translations []Translation `json:"translations"`
	} `json:"data"`
}

type Translation struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
}

type RequestBody struct {
	Source string   `json:"source,omitempty"`
	Target string   `json:"target"`
	Query  []string `json:"q"`
}

func main() {
	// https://cloud.google.com/docs/authentication/provide-credentials-adc#local-key
	token := flag.String("token", "", "Google access token (defaults to $GOOGLE_TOKEN)")
	target := flag.String("target", "", "Destination language (two-letter code); defaults to VN")
	source := flag.String("source", "", "Source language (two-letter code); auto-detected by default")

	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *token == "" {
		*token = os.Getenv("GOOGLE_TOKEN")
	}
	if *token == "" {
		log.Fatal("Missing required Google token")
	}

	if *target == "" {
		*target = "vn"
	}

	// Create body
	reqBody := RequestBody{
		Target: *target,
		Query:  flag.Args(),
	}
	if *source == "" {
		reqBody.Source = *source
	}

	reqRaw, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	bodyReader := bytes.NewReader(reqRaw)

	req, err := http.NewRequest("POST", GoogleTransApiUrl, bodyReader)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+*token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	respRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse response body
	var respBody ResponseBody
	if err := json.Unmarshal(respRaw, &respBody); err != nil {
		log.Fatal(err)
	}

	for _, t := range respBody.Data.Translations {
		fmt.Printf("%s (%s)\n", html.UnescapeString(t.TranslatedText), t.DetectedSourceLanguage)
	}
}
