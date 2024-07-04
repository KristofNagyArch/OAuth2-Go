package handlers

import (
	"fmt"
	"log"
	"net/http"
	"oauth2-go/cache"
)

/*
 * Sample QBO API call to get CompanyInfo using OAuth2 tokens
 */
func GetTokenInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering GetTokenInfo ")
	accessToken := cache.GetFromCache("access_token")
	fmt.Fprintf(w, accessToken)
}
