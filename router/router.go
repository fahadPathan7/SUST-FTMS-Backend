package router

import (
	"ftms/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter() // create mux router. it will be used to register routes

	// register routes. we are calling controller functions for each route
	// router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	// router.HandleFunc("/api/movie/{id}", controller.GetOneMovie).Methods("GET")
	router.HandleFunc("/api/dept", controller.InsertNewDept).Methods("POST")
	// router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT")
	// router.HandleFunc("/api/movie/{id}", controller.DeleteOneMovie).Methods("DELETE")
	// router.HandleFunc("/api/dmovies", controller.DeleteAllMovies).Methods("DELETE")

	return router
}