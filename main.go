package main

import (
	"encoding/json"
	"fmt"
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
	bookId, _ := strconv.Atoi(val)
	index := indexByID(books, bookId)
	if index < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books[index]); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	val := vars["id"]
	bookId, _ := strconv.Atoi(val)
	index := indexByID(books, bookId)
	if index < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	books = append(books[:index], books[index+1:]...)
	w.WriteHeader(http.StatusOK)
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	bookId, _ := strconv.Atoi(id)
	index := indexByID(books, bookId)
	if index < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b := Book{}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b.Id = bookId
	books[index] = b
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&b); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func initializeBooks() {
	books = Books{
		Book{Id: 1, Name: "Book 1", Description: "Book one"},
		Book{Id: 2, Name: "Book 2", Description: "Book two"},
	}
}
func indexByID(books []Book, id int) int {
	for i := 0; i < len(books); i++ {
		if books[i].Id == id {
			return i
		}
	}
	return -1
}
func main() {
	initializeBooks()
	r := mux.NewRouter()
	booksRouter := r.PathPrefix("/books").Subrouter()
	// Routes consist of a path and a handler function.
	booksRouter.HandleFunc("", ListBooks).Methods(http.MethodGet)
	booksRouter.HandleFunc("", AddBook).Methods(http.MethodPost)
	booksRouter.HandleFunc("/{id}", ShowBook).Methods(http.MethodGet)
	booksRouter.HandleFunc("/{id}", DeleteBook).Methods(http.MethodDelete)
	booksRouter.HandleFunc("/{id}", UpdateBook).Methods(http.MethodPut)
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
