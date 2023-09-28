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
	router.HandleFunc("/api/player/{playerRegNo}", controller.GetAPlayer).Methods("GET") // to get a specific player
	router.HandleFunc("/api/depts", controller.GetAllDepts).Methods("GET") // to get all depts
	router.HandleFunc("/api/dept/{deptCode}", controller.GetADept).Methods("GET") // to get a specific dept
	router.HandleFunc("/api/tournaments", controller.GetAllTournaments).Methods("GET") // to get all tournaments
	router.HandleFunc("/api/tournament/{tournamentId}", controller.GetATournament).Methods("GET") // to get a specific tournament
	router.HandleFunc("/api/tournament/teams/{tournamentId}", controller.GetAllTeamsOfATournament).Methods("GET") // to get all teams of a tournament
	router.HandleFunc("/api/dept/players/{deptCode}", controller.GetPlayersOfADept).Methods("GET") // to get all players of a dept
	router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.GetATeamOfATournament).Methods("GET") // to get a team of a tournament
	router.HandleFunc("/api/tournament/matches/{tournamentId}", controller.GetAllMatchesOfATournament).Methods("GET") // to get all matches of a tournament
	router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.GetAMatchOfATournament).Methods("GET") // to get a match of a tournament
	router.HandleFunc("/api/referees", controller.GetAllReferees).Methods("GET") // to get all referees
	router.HandleFunc("/api/referee/{refereeId}", controller.GetAReferee).Methods("GET") // to get a specific referee
	router.HandleFunc("/api/tournament/tiebreakers/{tournamentId}", controller.GetAllTiebreakersOfATournament).Methods("GET") // to get all tiebreakers of a tournament
	router.HandleFunc("/api/tournament/tiebreaker/{tournamentId}/{matchId}", controller.GetATiebreakerOfATournament).Methods("GET") // to get a tiebreaker of a tournament
	router.HandleFunc("/api/tournament/individualscores/{tournamentId}", controller.GetAllIndividualScoresOfATournament).Methods("GET") // to get all individualscores (all players) of a tournament
	router.HandleFunc("/api/tournament/player/individualscores/{tournamentId}/{playerRegNo}", controller.GetAllIndividualScoresOfAPlayerInATournament).Methods("GET") // to get all individualscores of a player in a tournament
	router.HandleFunc("/api/tournament/match/individualscores/{tournamentId}/{matchId}", controller.GetAllIndividualScoresOfAMatch).Methods("GET") // to get all individualscores of a match
	router.HandleFunc("/api/tournament/individualpunishments/{tournamentId}", controller.GetAllIndividualPunishmentsOfATournament).Methods("GET") // to get all individualpunishments (all players) of a tournament
	router.HandleFunc("/api/tournament/match/individualpunishments/{tournamentId}/{matchId}", controller.GetAllIndividualPunishmentsOfAMatch).Methods("GET") // to get all individualpunishments of a match
	router.HandleFunc("/api/tournament/player/individualpunishments/{tournamentId}/{playerRegNo}", controller.GetAllIndividualPunishmentsOfAPlayerInATournament).Methods("GET") // to get all individual punishments of a player in a tournament

	// PUT operations
	router.HandleFunc("/api/player/{playerRegNo}", controller.UpdateAPlayer).Methods("PUT") // to update a player
	router.HandleFunc("/api/dept/{deptCode}", controller.UpdateADept).Methods("PUT") // to update a dept
	router.HandleFunc("/api/tournament/{tournamentId}", controller.UpdateATournament).Methods("PUT") // to update a tournament
	router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.UpdateATeam).Methods("PUT") // to update a team
	router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.UpdateAMatch).Methods("PUT") // to update a match
	router.HandleFunc("/api/referee/{refereeId}", controller.UpdateAReferee).Methods("PUT") // to update a referee

	return router
}