package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Books []Book

var books Books

func ListBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.Id = len(books) + 1
	books = append(books, book)
	json.NewEncoder(w).Encode(&book)
}
func ShowBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	val := vars["id"]
	bookId, err := strconv.Atoi(val)
	if err != nil {
		json.NewEncoder(w).Encode(nil)
		return
	}
	for _, b := range books {
		if b.Id == bookId {
			json.NewEncoder(w).Encode(b)
			return
		}
	}
	json.NewEncoder(w).Encode(nil)
}
func initializeBooks() {
	books = Books{
		Book{Id: 1, Name: "Book 1", Description: "Book one"},
		Book{Id: 2, Name: "Book 2", Description: "Book two"},
	}
}
func main() {
	initializeBooks()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/books", ListBooks).Methods("GET")
	r.HandleFunc("/books", AddBook).Methods("POST")
	r.HandleFunc("/getBook/{id}", ShowBook).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
