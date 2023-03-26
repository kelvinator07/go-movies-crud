package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	lastMovie := movies[len(movies) - 1]
	lastMovieID, _ := strconv.Atoi(lastMovie.ID)
	movie.ID = strconv.Itoa(lastMovieID + 1)
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movieIndex int
	var movie Movie

	_ = json.NewDecoder(r.Body).Decode(&movie)
	params := mux.Vars(r)

	for index, item := range movies {
		if params["id"] == item.ID {
			movieIndex = index
			break
		}
	}

	m := &movies[movieIndex]
	m.Isbn = movie.Isbn
	m.Title = movie.Title
	m.Director = movie.Director
	
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	movieId := params["id"]
	for index, movie := range movies {
		if movie.ID == movieId {
			movies = append(movies[:index], movies[index + 1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode("Movie with id: " + movieId + " successfully deleted")
}

func index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Go Movies CRUD API")
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "12345", Title: "John Wick Part One", Director: &Director{Firstname: "Kelvin", Lastname: "Geeky"}})
	movies = append(movies, Movie{ID: "2", Isbn: "67890", Title: "John Wick Part Two", Director: &Director{Firstname: "Steven", Lastname: "Speilberg"}})

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
