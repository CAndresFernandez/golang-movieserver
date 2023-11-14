package main

// import any necessary packages for running the application, any external packages need to be installed with 'go get' only after creating a module with 'git init mod #projectname'
import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// a struct is like an entity, with properties to define as necessary
type Movie struct{
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	// the * points to the struct Director, it defines a relation between two structs
	Director *Director `json:"director"`
}

type Director struct{
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

// create a slice (array) and define what type of data will be in it
var movies []Movie

// findAll movies
// w - write, r - request
func getMovies(w http.ResponseWriter, r *http.Request) {
// set the header parameters, in this case content-type: json
	w.Header().Set("Content-Type", "application/json")
// encode the json, in this case encode the entire slice 'movies' since it's a findAll
	json.NewEncoder(w).Encode(movies)
}

// delete a movie by its {id}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// set some parameters
	params := mux.Vars(r)
	// loop over movies. index = key, item = value, range = the slice (array)
	for index, item := range movies {
		if item.ID == params["id"] {
			// the item that matches the ID passed in params will be replaced by everything that comes after it (movies[:index+1]...), effectively deleting the item in question
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	// return updated movie slice
	json.NewEncoder(w).Encode(movies)
}

// recuperates a movie by its {id}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// use the blank as a placeholder when we don't want to use index
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// creates a new movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// define variable of type Movie
	var movie Movie
	// we've sent data that needs to be decoded
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// set the id with a randomly generated integer converted to a string
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	// append the new Movie to the slice movies
	movies = append(movies, movie)
	// return json that shows the newly created Movie
	json.NewEncoder(w).Encode(movie)
}

// updates an existing movie (in this case we're just going to delete an existing movie and replacing it, which is not the practice when we're working with databases)
func updateMovie(w http.ResponseWriter, r *http.Request) {
// set json content type
w.Header().Set("Content-Type", "application/json")
// access params
params := mux.Vars(r)
// loop over the slice to find the movie to update
for index, item := range movies {
	if item.ID == params["id"] {
		// delete the movie like we did above
		movies = append(movies[:index], movies[index+1:]...)
		// create a new movie object from what we've sent in the body
		var movie Movie
		// decode it
		_ = json.NewDecoder(r.Body).Decode(&movie)
		// set the new movie object's id as the one we initially meant to update
		movie.ID = params["id"]
		// append it to the slice
		movies = append(movies, movie)
		// return the updated movie as json
		json.NewEncoder(w).Encode(movie)
		return
	}
}
}

// the main function for running the application
func main() {
	r := mux.NewRouter()

	// add some movies to the slice "movies" so we have some data to work with without a database connection
	// the & affirms the relation between the Movie we're adding to the slice 'movies' and the related the Director
	movies = append(movies, Movie{ID: "1", Isbn:"438227", Title:"Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "45455", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	// standard CRUD architecture for our API
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// something to print when the server is running
	fmt.Printf("Starting server at port 8000\n")
	// ListenAndServe is the request/response mechanism. define the port and the handler, in this case r (mux Router)
	log.Fatal(http.ListenAndServe(":8000", r))

}

