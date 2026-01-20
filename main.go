package main

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/bigspaceships/circlejerk/auth"
	"github.com/computersciencehouse/csh-auth"

	// "github.com/gin-contrib/static"
	// "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed static
var server embed.FS

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	csh := csh_auth.CSHAuth{}
	csh.Init(
		os.Getenv("OIDC_CLIENT_ID"),
		os.Getenv("OIDC_CLIENT_SECRET"),
		os.Getenv("JWT_SECRET"),
		os.Getenv("STATE"),
		os.Getenv("HOST"),
		os.Getenv("HOST")+"/auth/callback",
		os.Getenv("HOST")+"/auth/login",
		[]string{"profile", "groups"},
	)

	myAuth := auth.Config{
		ClientId: os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		State: os.Getenv("STATE"),
		RedirectURI: os.Getenv("HOST")+"/auth/callback",
	}

	myAuth.SetupAuth()

	http.HandleFunc("/auth/login", myAuth.LoginRequest)
	http.HandleFunc("/auth/callback", myAuth.LoginCallback)

	// fs, err := static.EmbedFolder(server, "static")
	// if err != nil {
	// 	panic(err)
	// }


	// r.Use(csh.AuthWrapper(static.Serve("/", fs)))

	// r.GET("/api/ping", csh.AuthWrapper(func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// }))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
