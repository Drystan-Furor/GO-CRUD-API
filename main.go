package main

import (
	"encoding/json" //for coding data into json for sending to Postman
	"fmt"           //for string printing out stuff
	"log"           // for logging errors
	"math/rand"     // random for ID generation
	"net/http"      //allowes a server in Golang
	"strconv"       // to convert math/rand, using itoa

	"github.com/gorilla/mux" // the external library
)

//structs and slices
// struct/class == type of object with key value pairs

//struct
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"` //Director is a pointer
	// it will have the same values that i'll define inside Director
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	// sending to postman I'll be using: isbn, id, title and director smallcaps
}

//slice of movies
var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content-Type", "application/json")
	//params
	params := mux.Vars(r)
	//loop over the movies range
	///delete the movie with the i.d. that you've sent
	//add new movie - the movie that we send in the body of postman

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter() //r now is new router

	// append a movie, a slice from Movie, define the movie's properties
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "45455", Title: "Movie Two", Director: &Director{Firstname: "Steven", Lastname: "Smith"}})
	// this means we basically have 2 movies to test if the API is working correctly because I won't be using a database connection here

	// HandleFunc = if my handle is "/string" USE funcion
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/id", getMovie).Methods("GET")
	r.HandleFunc("/movies/", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at localhost:8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
