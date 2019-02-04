package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	// [START imports]
	firebase "firebase.google.com/go"
	// [END imports]

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

//[START firebase config variable]
var (
	firebaseConfig = &firebase.Config{
		DatabaeURL:    "https://sma-cloud-storage.firebaseio.com",
		ProjectID:     "sma-cloud-storage",
		StorageBucket: "sma-cloud-storage.appspot.com",
	}
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

// [END firebase config variable]
// [START new_post_field]

type Post struct {
	Author  string
	UserID  string
	Message string
	Posted  time.Time
}

// [END new_post_field]

type templateParams struct {
	Notice  string
	Name    string
	Message string
	Posts   []Post
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	ctx := appengine.NewContext(r)
	params := templateParams{}

	q := datastore.NewQuery("Post").Order("-Posted").Limit(20)
	if _, err := q.GetAll(ctx, &params.Posts); err != nil {
		log.Errorf(ctx, "Getting posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't get latest posts. Refresh?"
		indexTemplate.Execute(w, params)
		return
	}

	if r.Method == "GET" {
		indexTemplate.Execute(w, params)
		return
	}
	// It's a POST request, so handle the form submission.

	// [START firebase_token]
	message := r.FormValue("message")

	// Create a new Firebase App.
	app, err := firebase.NewApp(ctx, firebaseConfig)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Create a new authenticator for the app.
	auth, err := app.Auth(ctx)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Verify the token passed in by the user is valid.
	tok, err := auth.VerifyIDTokenAndCheckRevoked(ctx, r.FormValue("token"))
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}
	// Use the validated token to get the user's information.
	user, err := auth.GetUser(ctx, tok.UID)
	if err != nil {
		params.Notice = "Couldn't authenticate. Try logging in again?"
		params.Message = message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	// [END firebase_token]

	// [START logged_in_post]
	post := Post{
		UserID:  user.UID, // Include UserID in case Author isn't unique.
		Author:  user.DisplayName,
		Message: message,
		Posted:  time.Now(),
	}
	// [END logged_in_post]

	params.Name = post.Author

	if post.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		params.Notice = "No message provided"
		indexTemplate.Execute(w, params)
		return
	}
	key := datastore.NewIncompleteKey(ctx, "Post", nil)
	if _, err := datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		params.Notice = "Couldn't add new post. Try again?"
		params.Message = post.Message // Preserve their message so they can try again.
		indexTemplate.Execute(w, params)
		return
	}

	// Prepend the post that was just added.
	params.Posts = append([]Post{post}, params.Posts...)
	params.Notice = fmt.Sprintf("Thank you for your submission, %s!", post.Author)
	indexTemplate.Execute(w, params)
}
