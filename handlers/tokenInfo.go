package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"oauth2-go/domain"
)

/*
 * Sample QBO API call to get CompanyInfo using OAuth2 tokens
 */
func GetTokenInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering GetTokenInfo ")
	accessToken := domain.RedisRepositoryImpl.FindAccessTokenByPartnerID(context.Background(), "123")
	fmt.Fprintf(w, accessToken)
}
