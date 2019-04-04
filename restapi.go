package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type post struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	IsPublished bool `json:"ispublished"`
}

// total beschissene Art, drei posts in ein slice von Posts zu tun, aber zum testen genuegts
var eins post = post{ID: "1", Title: "Erster Post", Content: "Dies ist der erste Post.", IsPublished: true } 
var zwei post = post{ID: "2", Title: "Zweiter Post", Content: "Dies ist der zweite Post.", IsPublished: true }
var drei post = post{ID: "3", Title: "Letzter Post", Content: "Dies ist der letzte Post.", IsPublished: true }
var posts = []post{eins, zwei, drei}

// Liesst alle Posts und gibt ein json mit allen posts als response aus
func GetPosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(posts)
}

// Liesst alle posts, vergleicht ob Id stimmt, gibt post mit passender Id als json aus
func GetPost(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	for _, indexpost := range posts {
		if parameters["id"] == indexpost.ID {
			json.NewEncoder(w).Encode(indexpost)
		}
	}
}

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

