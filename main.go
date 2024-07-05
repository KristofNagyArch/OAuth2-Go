package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"oauth2-go/cache"
	"oauth2-go/config"
	"oauth2-go/handlers"
)

func main() {

	InitRepositories(context.Background())

	//call discovery
	handlers.CallDiscoveryAPI()

	//register static routes
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/connected/", http.StripPrefix("/connected/", http.FileServer(http.Dir("static/connected/"))))

	//register handler routes
	http.HandleFunc("/oauth2redirect", handlers.CallBackFromOAuth)
	http.HandleFunc("/connectToQuickbooks", handlers.ConnectToQuickbooks)
	http.HandleFunc("/getToken", handlers.GetTokenInfo)

	http.HandleFunc("/getCompanyInfo", handlers.GetCompanyInfo)
	http.HandleFunc("/refreshToken", handlers.RefreshToken)
	http.HandleFunc("/revokeToken", handlers.RevokeToken)
	http.HandleFunc("/signInWithIntuit", handlers.SignInWithIntuit)
	http.HandleFunc("/getAppNow", handlers.GetAppNow)

	//log and start server
	log.Println("running server on ", config.OAuthConfig.Port)
	log.Fatal(http.ListenAndServe(config.OAuthConfig.Port, nil))

}

func createRedisConnection(rootCtx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.OAuthConfig.CacheAddress,
		DB:       config.OAuthConfig.CacheDatabaseNumber,
		PoolSize: config.OAuthConfig.CachePoolSize,
	})
	s := rdb.Ping(rootCtx)
	if s.Err() != nil {
		return nil, s.Err()
	}
	return rdb, nil
}

// InitRepositories is responsible for initialization of ticket and admission control repositories
func InitRepositories(ctx context.Context) {
	redisClient, err := createRedisConnection(ctx)
	if err != nil {
		panic(err)
	}
	cache.InitTokenRedisRepository(redisClient, config.OAuthConfig.ApplicationEnvironment, config.OAuthConfig.ApplicationName)
}
