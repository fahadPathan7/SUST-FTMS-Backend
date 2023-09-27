package router

import (
	"ftms/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter() // create mux router. it will be used to register routes

	// POST operations
	router.HandleFunc("/api/dept", controller.InsertNewDept).Methods("POST") // to insert new dept
	router.HandleFunc("/api/player", controller.InsertNewPlayer).Methods("POST") // to insert new player
	router.HandleFunc("/api/team", controller.InsertNewTeam).Methods("POST") // to insert new team
	router.HandleFunc("/api/tournament", controller.InsertNewTournament).Methods("POST") // to insert new tournament
	router.HandleFunc("/api/referee", controller.InsertNewReferee).Methods("POST") // to insert new referee
	router.HandleFunc("/api/match", controller.InsertNewMatch).Methods("POST") // to insert new match
	router.HandleFunc("/api/tiebreaker", controller.InsertNewTiebreaker).Methods("POST") // to insert new tiebreaker
	router.HandleFunc("/api/individualpunishment", controller.InsertNewIndividualPunishment).Methods("POST") // to insert new individualpunishment
	router.HandleFunc("/api/individualscore", controller.InsertNewIndividualScore).Methods("POST") // to insert new individualscore

	// GET operations
	router.HandleFunc("/api/depts", controller.GetAllDepts).Methods("GET") // to get all depts
	router.HandleFunc("/api/dept/{deptCode}", controller.GetADept).Methods("GET") // to get a specific dept
	router.HandleFunc("/api/tournaments", controller.GetAllTournaments).Methods("GET") // to get all tournaments
	router.HandleFunc("/api/tournament/teams/{tournamentId}", controller.GetAllTeamsOfATournament).Methods("GET") // to get all teams of a tournament
	router.HandleFunc("/api/dept/players/{deptCode}", controller.GetPlayersOfADept).Methods("GET") // to get all players of a dept
	router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.GetATeamOfATournament).Methods("GET") // to get a team of a tournament
	router.HandleFunc("/api/tournament/matches/{tournamentId}", controller.GetAllMatchesOfATournament).Methods("GET") // to get all matches of a tournament

	return router
}