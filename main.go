package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "homepage Endpoint Hit")
}

type Article struct {
	Title string `json:"Title"`
	Desc string `json:"Desc"`
	Content string `json:"content"`
}

type Articles []Article

func allArticles(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Test Title", Desc: "Test Desc", Content: "Hello World"},
	}

	fmt.Println("Endpoint Hit: All articles Endpoint")
	json.NewEncoder(w).Encode(articles)
}

func testPostArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test post endpoint hit")
}

//******************************

func compute(value int) {
	for i := 0; i < value; i++ {
		time.Sleep(time.Second)
		fmt.Println(i)
	}
}

func sleepyCompute() {

	// regular calls
	compute(5)
	compute(5)

	// goroutine calls
	go compute(5)
	go compute(5)

	// pause the main so threads won't terminate early
	fmt.Scanln()
}

//******************************

func marshallingPb() {
	paul := &Person {
		Name: "Paul",
		Age: 31,
		SocialFollowers: &SocialFollowers{
			Twitter: 500,
			Youtube: 1000,
		},
	}

	data, err := proto.Marshal(paul)
	if (err) != nil {
		log.Fatal("Marshalling error", err)
	}

	fmt.Println(data)

	newPaul := &Person{}
	err = proto.Unmarshal(data, newPaul)
	if err != nil {
		log.Fatal("unmarshalling error: ", err)
	}

	fmt.Println(newPaul.GetName())
	fmt.Println(newPaul.GetAge())
	fmt.Println(newPaul.SocialFollowers.GetTwitter())
	fmt.Println(newPaul.SocialFollowers.GetYoutube())
}

//******************************

func addMockData() {

	// Mock data
	books = append(books, &Book {
		ID: "1",
		Isbn: "412523",
		Title: "Go Go Gadget Lang",
		Author: &Author{
			FirstName: "Paul",
			LastName: "Choi",
		},
	})
	books = append(books, &Book {
		ID: "2",
		Isbn: "847502",
		Title: "The Quick Fox",
		Author: &Author{
			FirstName: "Sweet",
			LastName: "Feet",
		},
	})
	books = append(books, &Book {
		ID: "3",
		Isbn: "563891",
		Title: "The Lazy Dog",
		Author: &Author{
			FirstName: "Brown",
			LastName: "Cow",
		},
	})
	books = append(books, &Book {
		ID: "4",
		Isbn: "273906",
		Title: "Mean Green Machine",
		Author: &Author{
			FirstName: "Artie",
			LastName: "Lang",
		},
	})
	books = append(books, &Book {
		ID: "5",
		Isbn: "969423",
		Title: "The Curious Case of the Missing Semicolon",
		Author: &Author{
			FirstName: "Proc",
			LastName: "Tologist",
		},
	})

}

var books []*Book

// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// get a single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, book := range books {
		if book.GetID() == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// create a book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if (err) != nil {
		log.Fatal("Error decoding book", err)
	}

	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, &book)

	json.NewEncoder(w).Encode(&book)
}

//update a book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// delete the book
	var book Book

	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {

			err := json.NewDecoder(r.Body).Decode(&book)
			if (err) != nil {
				log.Fatal("Error decoding book", err)
			}

			book.ID = item.ID
			books[index] = &book

			json.NewEncoder(w).Encode(&book)
		}
	}
}

//delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)

	router.HandleFunc("/articles", allArticles).Methods("GET")
	router.HandleFunc("/articles", testPostArticles).Methods("POST")

	booksSubrouter := router.PathPrefix("/api/v1/books").Subrouter()
	booksSubrouter.HandleFunc("/", getBooks).Methods("GET")
	booksSubrouter.HandleFunc("/{id}", getBook).Methods("GET")
	booksSubrouter.HandleFunc("/", createBook).Methods("POST")
	booksSubrouter.HandleFunc("/{id}", updateBook).Methods("PUT")
	booksSubrouter.HandleFunc("/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {

	addMockData()

	handleRequests()

	//sleepyCompute()

	//marshallingPb()
}
