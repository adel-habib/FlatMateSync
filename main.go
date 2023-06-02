package main

import (
	"FlatMateSync/config"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig *oauth2.Config

func setOAuthClient(cfg config.Config) {
	oauthConfig = &oauth2.Config{
		ClientID:     cfg.Oauth.Client_Id,
		ClientSecret: cfg.Oauth.Client_Secret,
		RedirectURL:  "http://localhost:8080/oauth/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func RunMigrations(cfg *config.Config) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	fmt.Println(cfg.Database.Migrations_Url)
	// Run the migrations
	m, err := migrate.New(cfg.Database.Migrations_Url, dbURL)
	if err != nil {
		log.Fatalf("migration failed: %s", err.Error())
		return
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no database change necessary")
			return
		} else {
			log.Fatalf("migration failed: %s", err.Error())
		}
	}

	log.Println("Migration ran successfully")
}

func main() {
	fmt.Println("HI STARTED!!!")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	setOAuthClient(cfg)
	RunMigrations(&cfg)
	router := gin.Default()

	store := memstore.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	router.Use(sessions.Sessions("sid", store))
	router.GET("/", indexHandler)
	router.GET("/login", loginHandler)
	router.GET("/oauth/callback", callbackHandler)
	authorized := router.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/protected", protectedHandler)
	}

	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}

func indexHandler(c *gin.Context) {
	const html = `
    <!DOCTYPE html>
    <html>
        <head>
            <title>My Site</title>
        </head>
        <body>
			<h1> Login with Google!</h1>
			<h2> this is all fun and games </h2>
            <button onclick="location.href='/login'" type="button">Login</button>
        </body>
    </html>
    `
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func loginHandler(c *gin.Context) {
	session := sessions.Default(c)
	state := generateState()
	fmt.Println(state)
	codeVerifier := generateRandomString(64)
	codeChallenge := generateCodeChallenge(codeVerifier)
	session.Set("code_verifier", codeVerifier)
	session.Set("state", state)
	err := session.Save()
	if err != nil {
		fmt.Println("Error saving state: ", err.Error())
	}

	fmt.Println("State is", session.Get("state"))
	url := oauthConfig.AuthCodeURL(state, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	c.Redirect(http.StatusFound, url)
}

func callbackHandler(c *gin.Context) {

	state := c.Query("state")
	session := sessions.Default(c)
	savedState := session.Get("state")

	if state != savedState {
		fmt.Println("Given state ", state, "\nSaved state ", savedState)
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid state"))
		return
	}
	code := c.Query("code")
	codeVerifier := session.Get("code_verifier").(string)
	token, err := oauthConfig.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		session.Delete("state")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session.Set("token", token.AccessToken)
	session.Delete("state")
	session.Save()
	c.Redirect(http.StatusFound, "/protected")
}

func protectedHandler(c *gin.Context) {
	session := sessions.Default(c)
	token := session.Get("token")
	c.String(http.StatusOK, fmt.Sprintf("Hello, you're logged in with token: %s", token))
}

func generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func generateRandomString(size int) string {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token")
		if token == nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
