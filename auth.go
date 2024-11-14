package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/redis/go-redis/v9"
	"html/template"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY is not set in .env file")
	}

	store := sessions.NewCookieStore([]byte(sessionKey))
	gothic.Store = store
}

func auth(client *redis.Client, ctx context.Context) {
	var userTemplate = `
	<p><a href="/logout/{{.Provider}}">logout</a></p>
	<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
	<p>Email: {{.Email}}</p>
	<p>NickName: {{.NickName}}</p>
	<p>Location: {{.Location}}</p>
	<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
	<p>Description: {{.Description}}</p>
	<p>UserID: {{.UserID}}</p>
	<p>AccessToken: {{.AccessToken}}</p>
	<p>ExpiresAt: {{.ExpiresAt}}</p>
	<p>RefreshToken: {{.RefreshToken}}</p>
	`

	var indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login</title>
</head>
<body style="font-family: Arial, sans-serif; background-color: #f2f2f2; margin: 0; padding: 0; display: flex; justify-content: center; align-items: center; height: 100vh; background-image: url('https://images.unsplash.com/photo-1501594907359-97c6cd3e40eb?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=MnwzNjY0OHwwfDF8c2VhY2h8Mnx8Y29tcHV0ZXIlMjBsaWZlfGVufDB8fHx8fDE2OTk1MTM3NzE&ixlib=rb-1.2.1&q=80&w=1080'); background-size: cover; background-position: center;">
    
    <div style="text-align: center; background-color: rgba(0, 0, 0, 0.6); padding: 30px 50px; border-radius: 8px; width: 100%; max-width: 400px; color: white;">
        <h1 style="font-size: 2.5em; margin-bottom: 20px;">Welcome to Our Platform</h1>
        <p style="font-size: 1.2em; margin-bottom: 20px;">Please log in to continue</p>
        <a href="/auth/google" style="background-color: #4285F4; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-size: 1.2em; transition: background-color 0.3s ease;">Log in with Google</a>
    </div>

    <div style="position: absolute; bottom: 20px; width: 100%; text-align: center; font-size: 0.9em; color: white;">
        <p>Powered by YourCompany</p>
    </div>
</body>
</html>
`

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)

	// Create a Gin router
	r := gin.Default()

	// Callback route for handling the response from Google
	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		// Add provider to context to ensure gothic can use it
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
		userData(client, ctx, user)
		fmt.Println(user)
	})

	// Route to log out
	r.GET("/logout/:provider", func(c *gin.Context) {
		// Add provider to context
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))
		err := gothic.Logout(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	r.GET("/auth/:provider", func(c *gin.Context) {
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", c.Param("provider")))

		if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			t, _ := template.New("user").Parse(userTemplate)
			t.Execute(c.Writer, gothUser)
		} else {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		}
	})

	r.GET("/", func(c *gin.Context) {
		t, _ := template.New("index").Parse(indexTemplate)
		t.Execute(c.Writer, nil)
	})

	log.Println("Starting server on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func userData(client *redis.Client, ctx context.Context, user goth.User) {
	userKey := "user" + user.Email
	err := client.HSet(ctx, userKey, "name",
		user.Name, "email",
		user.Email, "nickname", user.NickName, "location",
		user.Location, "avatar_url", user.AvatarURL, "description",
		user.Description, "user_id",
		user.UserID, "access_token",
		user.AccessToken, "expires_at",
		user.ExpiresAt, "refresh_token",
		user.RefreshToken).Err()

	if err != nil {
		fmt.Println("Could not store user data:", err)
		return
	}
}
