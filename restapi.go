package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

// Very basic "Database" of posts, in form of a struct
type post struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	IsPublished bool `json:"ispublished"`
}
var eins post = post{ID: "1", Title: "Erster Post", Content: "Dies ist der erste Post.", IsPublished: true } 
var zwei post = post{ID: "2", Title: "Zweiter Post", Content: "Dies ist der zweite Post.", IsPublished: true }
var drei post = post{ID: "3", Title: "Letzter Post", Content: "Dies ist der letzte Post.", IsPublished: true }
var posts = []post{eins, zwei, drei}

// Sends all posts to client (in form of a json block)
func GetPosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(posts)
}

// Sends the selected post to client
func GetPost(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	for _, indexpost := range posts {
		if parameters["id"] == indexpost.ID {
			json.NewEncoder(w).Encode(indexpost)
		}
	}
}

// Deletes the post with corresponding ID in URL from the database
func DeletePost(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	l := len(posts)
	for i := 0; i < l; i++ {
		if posts[i].ID == parameters["id"] {
			posts = append(posts[:i], posts[i+1:]...)
			l--
			w.Write([]byte("Post successfully deleted!\n"))
			return
		}
	}
	w.Write([]byte("Post does not exist. Cannot delete.\n"))
}

// Decodes new post (json) from the client's request and appends it to the "database"
func CreatePost(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	for _, indexpost := range posts {
		if indexpost.ID == parameters["id"] {
			w.Write([]byte("Post with same ID already exists!\n"))
			return
		}
	}
	var newpost post
        json.NewDecoder(r.Body).Decode(&newpost)
	newpost.ID = parameters["id"]
        posts = append(posts, newpost)
        w.Write([]byte("Post successfully created!\n"))
}

// Updates post in database by replacing the struct with the corresponding ID from the slice
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	l := len(posts)
	for i := 0; i < l; i++ {
		if posts[i].ID == parameters["id"] {
			var updatedpost post
			json.NewDecoder(r.Body).Decode(&updatedpost)
			
			if updatedpost.ID != parameters["id"] {
				w.Write([]byte("Error. Trying to update post with wrong id!/n"))
				return
			}

			updatedleft := append(posts[:i], updatedpost) 
			posts = append(updatedleft, posts[i+1:]...)
			w.Write([]byte("Post successfully updated.\n"))
			l--
			return
		}
	}
	w.Write([]byte("Could not find post. Cannot update.\n"))
}

// main function
func main() {
	router := mux.NewRouter()
	
	router.HandleFunc("/post/", GetPosts).Methods("GET")
	router.HandleFunc("/post/{id}", GetPost).Methods("GET")
	router.HandleFunc("/post/{id}", DeletePost).Methods("DELETE")
	router.HandleFunc("/post/{id}", UpdatePost).Methods("PUT")
	router.HandleFunc("/post/{id}", CreatePost).Methods("POST")
	
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(router)
	}
}


