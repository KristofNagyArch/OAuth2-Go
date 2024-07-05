package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"oauth2-go/cache"
	"oauth2-go/domain"
)

/*
 * Call the refresh endpoint to generate new tokens
 */
func RefreshToken(w http.ResponseWriter, r *http.Request) {

	log.Println("Entering RefreshToken ")
	client := &http.Client{}
	data := url.Values{}

	//add parameters
	data.Set("grant_type", "refresh_token")
	refreshToken := cache.GetFromCache("refresh_token")
	refreshToken = domain.RedisRepositoryImpl.FindRefreshTokenByPartnerID(context.Background(), "123")
	data.Add("refresh_token", refreshToken)

	tokenEndpoint := cache.GetFromCache("token_endpoint")
	request, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//set the headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth())

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bearerTokenResponse, err := getBearerTokenResponse([]byte(body))
	//add the tokens to cache - in real app store in database
	cache.AddToCache("access_token", bearerTokenResponse.AccessToken)
	domain.RedisRepositoryImpl.StoreAccessTokenToPartnerID(context.Background(), "123", bearerTokenResponse.AccessToken)
	cache.AddToCache("refresh_token", bearerTokenResponse.RefreshToken)
	domain.RedisRepositoryImpl.StoreRefreshTokenToPartnerID(context.Background(), "123", bearerTokenResponse.RefreshToken)

	responseString := string(body)
	log.Println("Exiting RefreshToken ")
	fmt.Fprintf(w, responseString)

}
