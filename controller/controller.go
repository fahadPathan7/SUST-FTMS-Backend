package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ftms/models"
	"strconv"

	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

// connecting to mysql database
func CreateDbConnection() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/ftms")
	// port 3306 is the default port for mysql in xampp
	// here ftms is the database name

	if err != nil {
		fmt.Println("Error connecting databse!")
		panic(err.Error())
	}

	// defer db.Close()
	fmt.Println("Successfully connected to mysql database")
}

// verifications.

// check if player exists in database
func playerExists(playerRegNo int) bool {
	var player models.Player
	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if team exists in database
func teamExists(tournamentId string, deptCode int) bool {
	var team models.Team
	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManager, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19])

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if tournament exists in database
func tournamentExists(tournamentId string) bool {
	var tournament models.Tournament
	err := db.QueryRow("SELECT * FROM tbltournament WHERE tournamentId = ?", tournamentId).Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.TournamentYear)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if dept exists in database
func deptExists(deptCode int) bool {
	var dept models.Dept
	err := db.QueryRow("SELECT * FROM tbldept WHERE deptCode = ?", deptCode).Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptHeadName, &dept.DeptShortName)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if referee exists in database
func refereeExists(refereeID int) bool {
	var referee models.Referee
	err := db.QueryRow("SELECT * FROM tblreferee WHERE refereeID = ?", refereeID).Scan(&referee.RefereeID, &referee.RefereeName, &referee.RefereeInstitute)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if match exists in database
func matchExists(tournamentId string, matchId string) bool {
	var match models.Match
	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check if tiebreaker exists in database
func tiebreakerExists(tournamentId string, matchId string) bool {
	var tiebreaker models.Tiebreaker
	err := db.QueryRow("SELECT * FROM tbltiebreaker WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}

// check a team exists in a tournament or not
func teamExistsInATournament(tournamentId string, deptCode int) bool {
	var team models.Team
	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.DeptCode, &team.TeamSubmissionDate, &team.TeamManager, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19])

	if err != nil {
		return false
	}

	return true
}







// insert operations

// insert dept info into database
func insertNewDept(dept models.Dept) {
	// dept.DeptCode is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbldept VALUES (?, ?, ?, ?)", dept.DeptCode, dept.DeptName, dept.DeptHeadName, dept.DeptShortName)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new dept
func InsertNewDept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var dept models.Dept
	_ = json.NewDecoder(r.Body).Decode(&dept)

	if !deptExists(dept.DeptCode) {
		insertNewDept(dept)
		json.NewEncoder(w).Encode(dept)
	} else {
		json.NewEncoder(w).Encode("Dept already exists!")
	}
}





// insert player info into database
func insertNewPlayer(player models.Player) {
	// player.PlayerRegNo is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblplayer VALUES (?, ?, ?, ?, ?)", player.PlayerRegNo, player.PlayerSession, player.PlayerSemester, player.PlayerName, player.PlayerDeptCode)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new player
func InsertNewPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	// check if player already exists
	if playerExists(player.PlayerRegNo) {
		json.NewEncoder(w).Encode("Player already exists!")
		return
	}

	// check if dept exists
	if !deptExists(player.PlayerDeptCode) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// insert new player
	insertNewPlayer(player)
	json.NewEncoder(w).Encode(player)
}





// insert team info into database
func insertNewTeam(team models.Team) {
	// team.TournamentId is int type. and team.deptCode is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tblteam VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", team.TournamentId, team.TeamSubmissionDate, team.DeptCode, team.TeamManager, team.TeamCaptainRegID, team.PlayerRegNo[0], team.PlayerRegNo[1], team.PlayerRegNo[2], team.PlayerRegNo[3], team.PlayerRegNo[4], team.PlayerRegNo[5], team.PlayerRegNo[6], team.PlayerRegNo[7], team.PlayerRegNo[8], team.PlayerRegNo[9], team.PlayerRegNo[10], team.PlayerRegNo[11], team.PlayerRegNo[12], team.PlayerRegNo[13], team.PlayerRegNo[14], team.PlayerRegNo[15], team.PlayerRegNo[16], team.PlayerRegNo[17], team.PlayerRegNo[18], team.PlayerRegNo[19])

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
}

// return player's dept code from tblplayer in database
func getPlayerDeptCode(playerRegNo int) int {
	var player models.Player
	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

	if err != nil {
		panic(err.Error())
	}

	return player.PlayerDeptCode
}

// controller function to insert new team
func InsertNewTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var team models.Team
	_ = json.NewDecoder(r.Body).Decode(&team)

	// check if team already exists
	if teamExists(team.TournamentId, team.DeptCode) {
		json.NewEncoder(w).Encode("Team already exists!")
		return
	}

	// check if tournament exists
	if !tournamentExists(team.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if dept exists
	if !deptExists(team.DeptCode) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// check if all players exist
	if !playerExists(team.TeamCaptainRegID) {
		json.NewEncoder(w).Encode("Team captain doesn't exist!")
		return
	}
	for i := 0; i < 20; i++ {
		if !playerExists(team.PlayerRegNo[i]) {
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " doesn't exist!")
			return
		}
	}

	// check if all players are from same dept.
	if getPlayerDeptCode(team.TeamCaptainRegID) != team.DeptCode {
		json.NewEncoder(w).Encode("Team captain is not from this dept!")
		return
	}
	for i := 0; i < 20; i++ {
		if getPlayerDeptCode(team.PlayerRegNo[i]) != team.DeptCode {
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " is not from this dept!")
			return
		}
	}

	// check if team captain is in player list or not
	var captainFound bool = false
	for i := 0; i < 20; i++ {
		if team.PlayerRegNo[i] == team.TeamCaptainRegID {
			captainFound = true
			break
		}
	}
	if !captainFound {
		json.NewEncoder(w).Encode("Team captain is not in player list!")
		return
	}

	// check if player list has duplicate players
	for i := 0; i < 20 - 1; i++ {
		for j := i + 1; j < 20; j++ {
			if team.PlayerRegNo[i] == team.PlayerRegNo[j] {
				json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " and Player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}


	// insert new team
	insertNewTeam(team)
	json.NewEncoder(w).Encode(team)
}





// insert tournament info into database
func insertNewTournament(tournament models.Tournament) {
	// tournament.TournamentId is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbltournament VALUES(?, ?, ?)", tournament.TournamentId, tournament.TournamentName, tournament.TournamentYear)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new tournament
func InsertNewTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var tournament models.Tournament
	_ = json.NewDecoder(r.Body).Decode(&tournament)

	if !tournamentExists(tournament.TournamentId) {
		insertNewTournament(tournament)
		json.NewEncoder(w).Encode(tournament)
	} else {
		json.NewEncoder(w).Encode("Tournament already exists!")
	}
}





// insert referee info into database
func insertNewReferee(referee models.Referee) {
	// referee.RefereeID is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblreferee VALUES (?, ?, ?)", referee.RefereeID, referee.RefereeName, referee.RefereeInstitute)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new referee
func InsertNewReferee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var referee models.Referee
	_ = json.NewDecoder(r.Body).Decode(&referee)

	if refereeExists(referee.RefereeID) {
		json.NewEncoder(w).Encode("Referee already exists!")
		return
	}

	insertNewReferee(referee)
	json.NewEncoder(w).Encode(referee)
}





// insert match info into database
func insertNewMatch(match models.Match) {
	// match.TournamentId is int type. and match.MatchId is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tblmatch VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", match.TournamentId, match.MatchId, match.MatchDate, match.Team1DeptCode, match.Team2DeptCode, match.Team1Score, match.Team2Score, match.WinnerTeamDeptCode, match.MatchRefereeID, match.MatchLinesman1ID, match.MatchLinesman2ID, match.MatchFourthRefereeID)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new match
func InsertNewMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var match models.Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	// check if match already exists
	if matchExists(match.TournamentId, match.MatchId) {
		json.NewEncoder(w).Encode("Match already exists!")
		return
	}

	// check if tournament exists
	if !tournamentExists(match.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if team1 exists
	if !teamExists(match.TournamentId, match.Team1DeptCode) {
		json.NewEncoder(w).Encode("Team1 doesn't exist!")
		return
	}

	// check if team2 exists
	if !teamExists(match.TournamentId, match.Team2DeptCode) {
		json.NewEncoder(w).Encode("Team2 doesn't exist!")
		return
	}

	// check if referee exists
	if !refereeExists(match.MatchRefereeID) {
		json.NewEncoder(w).Encode("Referee doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman1ID) {
		json.NewEncoder(w).Encode("Linesman1 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman2ID) {
		json.NewEncoder(w).Encode("Linesman2 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchFourthRefereeID) {
		json.NewEncoder(w).Encode("Fourth referee doesn't exist!")
		return
	}

	// check if the winner team is one of the two teams
	if match.WinnerTeamDeptCode != match.Team1DeptCode && match.WinnerTeamDeptCode != match.Team2DeptCode {
		json.NewEncoder(w).Encode("Winner team is not one of the two teams!")
		return
	}

	// insert new match
	insertNewMatch(match)
	json.NewEncoder(w).Encode(match)
}





// insert tiebreaker info into database
func insertNewTiebreaker(tiebreaker models.Tiebreaker) {
	// tiebreaker.TournamentId is int type. and tiebreaker.MatchId is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tbltiebreaker VALUES (?, ?, ?, ?, ?, ?)", tiebreaker.TournamentId, tiebreaker.MatchId, tiebreaker.Team1DeptCode, tiebreaker.Team2DeptCode, tiebreaker.Team1TieBreakerScore, tiebreaker.Team2TieBreakerScore)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new tiebreaker
func InsertNewTiebreaker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var tiebreaker models.Tiebreaker
	_ = json.NewDecoder(r.Body).Decode(&tiebreaker)

	// check if tournament exists
	if !tournamentExists(tiebreaker.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(tiebreaker.TournamentId, tiebreaker.MatchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team1 exists
	if !teamExists(tiebreaker.TournamentId, tiebreaker.Team1DeptCode) {
		json.NewEncoder(w).Encode("Team1 doesn't exist!")
		return
	}

	// check if team2 exists
	if !teamExists(tiebreaker.TournamentId, tiebreaker.Team2DeptCode) {
		json.NewEncoder(w).Encode("Team2 doesn't exist!")
		return
	}

	// insert new tiebreaker
	insertNewTiebreaker(tiebreaker)
	json.NewEncoder(w).Encode(tiebreaker)
}





// insert individual score info into database
func insertNewIndividualScore(individualScore models.IndividualScore) {
	// individualScore.TournamentId is int type. and individualScore.MatchId is int type. and individualScore.PlayerRegNo is int type. and all are primary key.
	insert, err := db.Query("INSERT INTO tblindividualscore VALUES (?, ?, ?, ?, ?)", individualScore.TournamentId, individualScore.MatchId, individualScore.PlayerRegNo, individualScore.TeamDeptCode, individualScore.Goals)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new individual score
func InsertNewIndividualScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var individualScore models.IndividualScore
	_ = json.NewDecoder(r.Body).Decode(&individualScore)

	// check if tournament exists
	if !tournamentExists(individualScore.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(individualScore.TournamentId, individualScore.MatchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team exists
	if !teamExists(individualScore.TournamentId, individualScore.TeamDeptCode) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// check if player exists
	if !playerExists(individualScore.PlayerRegNo) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// insert new individual score
	insertNewIndividualScore(individualScore)
	json.NewEncoder(w).Encode(individualScore)
}





// insert individual punishment info into database
func insertNewIndividualPunishment(individualPunishment models.IndividualPunishment) {
	// individualPunishment.TournamentId is int type. and individualPunishment.MatchId is int type. and individualPunishment.PlayerRegNo is int type. and all are primary key.
	insert, err := db.Query("INSERT INTO tblindividualpunishment VALUES (?, ?, ?, ?, ?)", individualPunishment.TournamentId, individualPunishment.MatchId, individualPunishment.PlayerRegNo, individualPunishment.TeamDeptCode, individualPunishment.PunishmentType)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new individual punishment
func InsertNewIndividualPunishment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var individualPunishment models.IndividualPunishment
	_ = json.NewDecoder(r.Body).Decode(&individualPunishment)

	// check if tournament exists
	if !tournamentExists(individualPunishment.TournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(individualPunishment.TournamentId, individualPunishment.MatchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team exists
	if !teamExists(individualPunishment.TournamentId, individualPunishment.TeamDeptCode) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// check if player exists
	if !playerExists(individualPunishment.PlayerRegNo) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// insert new individual punishment
	insertNewIndividualPunishment(individualPunishment)
	json.NewEncoder(w).Encode(individualPunishment)
}










// getting info from database

// get all depts from database
func getAllDepts() []models.Dept {
	var dept models.Dept
	var depts []models.Dept

	result, err := db.Query("SELECT * FROM tbldept")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptHeadName, &dept.DeptShortName)

		if err != nil {
			panic(err.Error())
		}

		depts = append(depts, dept)
	}

	return depts
}

// controller function to get all depts
func GetAllDepts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var depts []models.Dept
	depts = getAllDepts()

	json.NewEncoder(w).Encode(depts)
}





// get a dept from database
func getADept(deptCode int) models.Dept {
	var dept models.Dept

	err := db.QueryRow("SELECT * FROM tbldept WHERE deptCode = ?", deptCode).Scan(&dept.DeptCode, &dept.DeptName, &dept.DeptHeadName, &dept.DeptShortName)

	if err != nil {
		//panic(err.Error())
		return models.Dept{}
	}

	return dept
}

// controller function to get a dept
func GetADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/dept/{deptCode}", controller.GetADept).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// dept exists or not
	if !deptExists(id) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	var dept models.Dept
	dept = getADept(id)

	json.NewEncoder(w).Encode(dept)
}





// get all the tournaments
func getAllTournaments() []models.Tournament {
	var tournament models.Tournament
	var tournaments []models.Tournament

	result, err := db.Query("SELECT * FROM tbltournament")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.TournamentYear)

		if err != nil {
			panic(err.Error())
		}

		tournaments = append(tournaments, tournament)
	}

	return tournaments
}

// controller function to get all tournaments
func GetAllTournaments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	var tournaments []models.Tournament
	tournaments = getAllTournaments()

	json.NewEncoder(w).Encode(tournaments)
}





// get a tournament
func getATournament(tournamentId string) models.Tournament {
	var tournament models.Tournament

	err := db.QueryRow("SELECT * FROM tbltournament WHERE tournamentId = ?", tournamentId).Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.TournamentYear)

	if err != nil {
		//panic(err.Error())
		return models.Tournament{}
	}

	return tournament
}

// controller function to get a tournament
func GetATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/{tournamentId}", controller.GetATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var tournament models.Tournament
	tournament = getATournament(id)

	json.NewEncoder(w).Encode(tournament)
}





// get all teams of a tournament
func getAllTeamsOfATournament(tournamentId string) []models.Team {
	var team models.Team
	var teams []models.Team

	result, err := db.Query("SELECT * FROM tblteam WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManager, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19])

		if err != nil {
			panic(err.Error())
		}

		teams = append(teams, team)
	}

	return teams
}

// controller function to get all teams of a tournament
func GetAllTeamsOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/teams/{tournamentId}", controller.GetAllTeamsOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var teams []models.Team
	teams = getAllTeamsOfATournament(id)

	json.NewEncoder(w).Encode(teams)
}





// get players of a dept
func getPlayersOfADept(deptCode int) []models.Player {
	var player models.Player
	var players []models.Player

	result, err := db.Query("SELECT * FROM tblplayer WHERE playerDeptCode = ?", deptCode)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

		if err != nil {
			panic(err.Error())
		}

		players = append(players, player)
	}

	return players
}

// controller function to get players of a dept
func GetPlayersOfADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/dept/players/{deptCode}", controller.GetPlayersOfADept).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// dept exists or not
	if !deptExists(id) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	var players []models.Player
	players = getPlayersOfADept(id)

	json.NewEncoder(w).Encode(players)
}





// get a team of a tournament
func getATeamOfATournament(tournamentId string, deptCode int) models.Team {
	var team models.Team

	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManager, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19])

	if err != nil {
		//panic(err.Error())
		return models.Team{}
	}

	return team
}

// controller function to get a team of a tournament
func GetATeamOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.GetATeamOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	tournamentId, _ := params["tournamentId"]

	deptCode, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// team exists or not
	if !teamExists(tournamentId, deptCode) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	var team models.Team
	team = getATeamOfATournament(tournamentId, deptCode)

	json.NewEncoder(w).Encode(team)
}





// get all matches of a tournament
func getAllMatchesOfATournament(tournamentId string) []models.Match {
	var match models.Match
	var matches []models.Match

	result, err := db.Query("SELECT * FROM tblmatch WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID)

		if err != nil {
			panic(err.Error())
		}

		matches = append(matches, match)
	}

	return matches
}

// controller function to get all matches of a tournament
func GetAllMatchesOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/matches/{tournamentId}", controller.GetAllMatchesOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var matches []models.Match
	matches = getAllMatchesOfATournament(id)

	json.NewEncoder(w).Encode(matches)
}





// get a match of a tournament
func getAMatchOfATournament(tournamentId string, matchId string) models.Match {
	var match models.Match

	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID)

	if err != nil {
		//panic(err.Error())
		return models.Match{}
	}

	return match
}

// controller function to get a match of a tournament
func GetAMatchOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.GetAMatchOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]


	// match exists or not
	if !matchExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	var match models.Match
	match = getAMatchOfATournament(tournamentId, matchId)

	json.NewEncoder(w).Encode(match)
}





// get a player
func getAPlayer(playerRegNo int) models.Player {
	var player models.Player

	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode)

	if err != nil {
		//panic(err.Error())
		return models.Player{}
	}

	return player
}

// controller function to get a player
func GetAPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/player/{playerRegNo}", controller.GetAPlayer).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["playerRegNo"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// player exists or not
	if !playerExists(id) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	var player models.Player
	player = getAPlayer(id)

	json.NewEncoder(w).Encode(player)
}





// get all referees
func getAllReferees() []models.Referee {
	var referee models.Referee
	var referees []models.Referee

	result, err := db.Query("SELECT * FROM tblreferee")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&referee.RefereeID, &referee.RefereeName, &referee.RefereeInstitute)

		if err != nil {
			panic(err.Error())
		}

		referees = append(referees, referee)
	}

	return referees
}

// controller function to get all referees
func GetAllReferees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var referees []models.Referee
	referees = getAllReferees()

	json.NewEncoder(w).Encode(referees)
}





// get a referee
func getAReferee(refereeId int) models.Referee {
	var referee models.Referee

	err := db.QueryRow("SELECT * FROM tblreferee WHERE refereeID = ?", refereeId).Scan(&referee.RefereeID, &referee.RefereeName, &referee.RefereeInstitute)

	if err != nil {
		//panic(err.Error())
		return models.Referee{}
	}

	return referee
}

// controller function to get a referee
func GetAReferee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/referee/{refereeId}", controller.GetAReferee).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["refereeId"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// referee exists or not
	if !refereeExists(id) {
		json.NewEncoder(w).Encode("Referee doesn't exist!")
		return
	}

	var referee models.Referee
	referee = getAReferee(id)

	json.NewEncoder(w).Encode(referee)
}





// get all tiebreakers of a tournament
func getAllTiebreakersOfATournament(tournamentId string) []models.Tiebreaker {
	var tiebreaker models.Tiebreaker
	var tiebreakers []models.Tiebreaker

	result, err := db.Query("SELECT * FROM tbltiebreaker WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

		if err != nil {
			panic(err.Error())
		}

		tiebreakers = append(tiebreakers, tiebreaker)
	}

	return tiebreakers
}

// controller function to get all tiebreakers of a tournament
func GetAllTiebreakersOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/tiebreakers/{tournamentId}", controller.GetAllTiebreakersOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var tiebreakers []models.Tiebreaker
	tiebreakers = getAllTiebreakersOfATournament(id)

	json.NewEncoder(w).Encode(tiebreakers)
}





// get a tiebreaker of a tournament
func getATiebreakerOfATournament(tournamentId string, matchId string) models.Tiebreaker {
	var tiebreaker models.Tiebreaker

	err := db.QueryRow("SELECT * FROM tbltiebreaker WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId).Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

	if err != nil {
		//panic(err.Error())
		return models.Tiebreaker{}
	}

	return tiebreaker
}

// controller function to get a tiebreaker of a tournament
func GetATiebreakerOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/tournament/tiebreaker/{tournamentId}/{matchId}", controller.GetATiebreakerOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// tiebreaker exists or not
	if !tiebreakerExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Tiebreaker doesn't exist!")
		return
	}

	var tiebreaker models.Tiebreaker
	tiebreaker = getATiebreakerOfATournament(tournamentId, matchId)

	json.NewEncoder(w).Encode(tiebreaker)
}





// get all individual scores of a tournament
func getAllIndividualScoresOfATournament(tournamentId string) []models.IndividualScore {
	var individualScore models.IndividualScore
	var individualScores []models.IndividualScore

	result, err := db.Query("SELECT * FROM tblindividualscore WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&individualScore.TournamentId, &individualScore.MatchId, &individualScore.PlayerRegNo, &individualScore.TeamDeptCode, &individualScore.Goals)

		if err != nil {
			panic(err.Error())
		}

		individualScores = append(individualScores, individualScore)
	}

	return individualScores
}

// controller function to get all individual scores of a tournament
func GetAllIndividualScoresOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/individualscores/{tournamentId}", controller.GetAllIndividualScoresOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var individualScores []models.IndividualScore
	individualScores = getAllIndividualScoresOfATournament(id)

	json.NewEncoder(w).Encode(individualScores)
}





// get all individual scores of a match
func getAllIndividualScoresOfAMatch(tournamentId string, matchId string) []models.IndividualScore {
	var individualScore models.IndividualScore
	var individualScores []models.IndividualScore

	result, err := db.Query("SELECT * FROM tblindividualscore WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&individualScore.TournamentId, &individualScore.MatchId, &individualScore.PlayerRegNo, &individualScore.TeamDeptCode, &individualScore.Goals)

		if err != nil {
			panic(err.Error())
		}

		individualScores = append(individualScores, individualScore)
	}

	return individualScores
}

// controller function to get all individual scores of a match
func GetAllIndividualScoresOfAMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/match/individualscores/{tournamentId}/{matchId}", controller.GetAllIndividualScoresOfAMatch).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	var individualScores []models.IndividualScore
	individualScores = getAllIndividualScoresOfAMatch(tournamentId, matchId)

	json.NewEncoder(w).Encode(individualScores)
}





// get all individual scores of a player in a tournament
func getAllIndividualScoresOfAPlayerInATournament(tournamentId string, playerRegNo int) []models.IndividualScore {
	var individualScore models.IndividualScore
	var individualScores []models.IndividualScore

	result, err := db.Query("SELECT * FROM tblindividualscore WHERE tournamentId = ? AND playerRegNo = ?", tournamentId, playerRegNo)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&individualScore.TournamentId, &individualScore.MatchId, &individualScore.PlayerRegNo, &individualScore.TeamDeptCode, &individualScore.Goals)

		if err != nil {
			panic(err.Error())
		}

		individualScores = append(individualScores, individualScore)
	}

	return individualScores
}

// controller function to get all individual scores of a player in a tournament
func GetAllIndividualScoresOfAPlayerInATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/player/individualscores/{tournamentId}/{playerRegNo}", controller.GetAllIndividualScoresOfAPlayerInATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and playerRegNo from url
	tournamentId, _ := params["tournamentId"]
	playerRegNo, _ := params["playerRegNo"]

	// convert playerRegNo from string to int
	playerRegNoInt, err := strconv.Atoi(playerRegNo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// tournament exists or not
	if !tournamentExists(tournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// player exists or not
	if !playerExists(playerRegNoInt) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	var individualScores []models.IndividualScore
	individualScores = getAllIndividualScoresOfAPlayerInATournament(tournamentId, playerRegNoInt)

	json.NewEncoder(w).Encode(individualScores)
}





// get all individual punishments of a tournament
func getAllIndividualPunishmentsOfATournament(tournamentId string) []models.IndividualPunishment {
	var individualPunishment models.IndividualPunishment
	var individualPunishments []models.IndividualPunishment

	result, err := db.Query("SELECT * FROM tblindividualpunishment WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&individualPunishment.TournamentId, &individualPunishment.MatchId, &individualPunishment.PlayerRegNo, &individualPunishment.TeamDeptCode, &individualPunishment.PunishmentType)

		if err != nil {
			panic(err.Error())
		}

		individualPunishments = append(individualPunishments, individualPunishment)
	}

	return individualPunishments
}

// controller function to get all individual punishments of a tournament
func GetAllIndividualPunishmentsOfATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/individualpunishments/{tournamentId}", controller.GetAllIndividualPunishmentsOfATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var individualPunishments []models.IndividualPunishment
	individualPunishments = getAllIndividualPunishmentsOfATournament(id)

	json.NewEncoder(w).Encode(individualPunishments)
}





// get all individual punishments of a match
func getAllIndividualPunishmentsOfAMatch(tournamentId string, matchId string) []models.IndividualPunishment {
	var individualPunishment models.IndividualPunishment
	var individualPunishments []models.IndividualPunishment

	result, err := db.Query("SELECT * FROM tblindividualpunishment WHERE tournamentId = ? AND matchId = ?", tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&individualPunishment.TournamentId, &individualPunishment.MatchId, &individualPunishment.PlayerRegNo, &individualPunishment.TeamDeptCode, &individualPunishment.PunishmentType)

		if err != nil {
			panic(err.Error())
		}

		individualPunishments = append(individualPunishments, individualPunishment)
	}

	return individualPunishments
}

// controller function to get all individual punishments of a match
func GetAllIndividualPunishmentsOfAMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/match/individualpunishments/{tournamentId}/{matchId}", controller.GetAllIndividualPunishmentsOfAMatch).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	var individualPunishments []models.IndividualPunishment
	individualPunishments = getAllIndividualPunishmentsOfAMatch(tournamentId, matchId)

	json.NewEncoder(w).Encode(individualPunishments)
}





// get all individual punishments of a player in a tournament
func getAllIndividualPunishmentsOfAPlayerInATournament(tournamentId string, playerRegNo int) []models.IndividualPunishment {
	var individualPunishment models.IndividualPunishment
	var individualPunishments []models.IndividualPunishment

	result, err := db.Query("SELECT * FROM tblindividualpunishment WHERE tournamentId = ? AND playerRegNo = ?", tournamentId, playerRegNo)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&individualPunishment.TournamentId, &individualPunishment.MatchId, &individualPunishment.PlayerRegNo, &individualPunishment.TeamDeptCode, &individualPunishment.PunishmentType)

		if err != nil {
			panic(err.Error())
		}

		individualPunishments = append(individualPunishments, individualPunishment)
	}

	return individualPunishments
}

// controller function to get all individual punishments of a player in a tournament
func GetAllIndividualPunishmentsOfAPlayerInATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/player/individualpunishments/{tournamentId}/{playerRegNo}", controller.GetAllIndividualPunishmentsOfAPlayerInATournament).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and playerRegNo from url
	tournamentId, _ := params["tournamentId"]
	playerRegNo, _ := params["playerRegNo"]

	// convert playerRegNo from string to int
	playerRegNoInt, err := strconv.Atoi(playerRegNo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// tournament exists or not
	if !tournamentExists(tournamentId) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// player exists or not
	if !playerExists(playerRegNoInt) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	var individualPunishments []models.IndividualPunishment
	individualPunishments = getAllIndividualPunishmentsOfAPlayerInATournament(tournamentId, playerRegNoInt)

	json.NewEncoder(w).Encode(individualPunishments)
}










// put operations

// update a tournament
func updateATournament(tournamentId string, tournament models.Tournament) {
	_, err := db.Query("UPDATE tbltournament SET tournamentName = ?, tournamentYear = ? WHERE tournamentId = ?", tournament.TournamentName, tournament.TournamentYear, tournamentId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a tournament
func UpdateATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/{tournamentId}", controller.UpdateATournament).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// get tournament from body
	var tournament models.Tournament
	_ = json.NewDecoder(r.Body).Decode(&tournament)

	updateATournament(id, tournament)

	json.NewEncoder(w).Encode(tournament)
}





// update a player
func updateAPlayer(playerRegNo int, player models.Player) {
	_, err := db.Query("UPDATE tblplayer SET playerSession = ?, playerSemester = ?, playerName = ?, playerDeptCode = ? WHERE playerRegNo = ?", player.PlayerSession, player.PlayerSemester, player.PlayerName, player.PlayerDeptCode, playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a player
func UpdateAPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/player/{playerRegNo}", controller.UpdateAPlayer).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["playerRegNo"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// player exists or not
	if !playerExists(id) {
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// get player from body
	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	updateAPlayer(id, player)

	json.NewEncoder(w).Encode(player)
}





// update a dept
func updateADept(deptCode int, dept models.Dept) {
	_, err := db.Query("UPDATE tbldept SET deptName = ?, deptShortName = ?, deptHeadName = ?, WHERE deptCode = ?", dept.DeptName, dept.DeptShortName, dept.DeptHeadName, deptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a dept
func UpdateADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/dept/{deptCode}", controller.UpdateADept).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// dept exists or not
	if !deptExists(id) {
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// get dept from body
	var dept models.Dept
	_ = json.NewDecoder(r.Body).Decode(&dept)

	updateADept(id, dept)

	json.NewEncoder(w).Encode(dept)
}





// update a team
func updateATeam(tournamentId string, deptCode int, team models.Team) {
	_, err := db.Query("UPDATE tblteam SET teamSubmissionDate = ?, teamManager = ?, teamCaptainRegID = ?, player1RegNo = ?, player2RegNo, player3RegNo = ?, player4RegNo = ?, player5RegNo = ?, player6RegNo = ?, player7RegNo = ?, player8RegNo = ?, player9RegNo = ?, player10RegNo = ?, player11RegNo = ?, player12RegNo = ?, player13RegNo = ?, player14RegNo = ?, player15RegNo = ?, player16RegNo = ?, player17RegNo = ?, player18RegNo = ?, player19RegNo = ?, player20RegNo = ? WHERE tournamentId = ? AND deptCode = ?", team.TeamSubmissionDate, team.TeamManager, team.TeamCaptainRegID, team.PlayerRegNo[0], team.PlayerRegNo[1], team.PlayerRegNo[2], team.PlayerRegNo[3], team.PlayerRegNo[4], team.PlayerRegNo[5], team.PlayerRegNo[6], team.PlayerRegNo[7], team.PlayerRegNo[8], team.PlayerRegNo[9], team.PlayerRegNo[10], team.PlayerRegNo[11], team.PlayerRegNo[12], team.PlayerRegNo[13], team.PlayerRegNo[14], team.PlayerRegNo[15], team.PlayerRegNo[16], team.PlayerRegNo[17], team.PlayerRegNo[18], team.PlayerRegNo[19], tournamentId, deptCode)
	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a team
func UpdateATeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.UpdateATeam).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and deptCode from url
	tournamentId, _ := params["tournamentId"]
	deptCode, _ := params["deptCode"]

	// convert deptCode from string to int
	deptCodeInt, err := strconv.Atoi(deptCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// team exists or not
	if !teamExists(tournamentId, deptCodeInt) {
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// get team from body
	var team models.Team
	_ = json.NewDecoder(r.Body).Decode(&team)


	// check if all players exist
	if !playerExists(team.TeamCaptainRegID) {
		json.NewEncoder(w).Encode("Team captain doesn't exist!")
		return
	}
	for i := 0; i < 20; i++ {
		if !playerExists(team.PlayerRegNo[i]) {
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " doesn't exist!")
			return
		}
	}

	// check if all players are from same dept.
	if getPlayerDeptCode(team.TeamCaptainRegID) != team.DeptCode {
		json.NewEncoder(w).Encode("Team captain is not from this dept!")
		return
	}
	for i := 0; i < 20; i++ {
		if getPlayerDeptCode(team.PlayerRegNo[i]) != team.DeptCode {
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " is not from this dept!")
			return
		}
	}

	// check if team captain is in player list or not
	var captainFound bool = false
	for i := 0; i < 20; i++ {
		if team.PlayerRegNo[i] == team.TeamCaptainRegID {
			captainFound = true
			break
		}
	}
	if !captainFound {
		json.NewEncoder(w).Encode("Team captain is not in player list!")
		return
	}

	// check if player list has duplicate players
	for i := 0; i < 20 - 1; i++ {
		for j := i + 1; j < 20; j++ {
			if team.PlayerRegNo[i] == team.PlayerRegNo[j] {
				json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " and Player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	updateATeam(tournamentId, deptCodeInt, team)

	json.NewEncoder(w).Encode(team)
}





// update a match
func updateAMatch(tournamentId string, matchId string, match models.Match) {
	_, err := db.Query("UPDATE tblmatch SET matchDate = ?, team1_deptCode = ?, team2_deptCode = ?, team1_goal_number = ?, team2_goal_number = ?, winner_team = ?, matchRefereeID = ?, matchLineman1ID = ?, matchLineman2ID = ?, matchFourthRefereeID = ? WHERE tournamentId = ? AND matchID = ?", match.MatchDate, match.Team1DeptCode, match.Team2DeptCode, match.Team1Score, match.Team2Score, match.WinnerTeamDeptCode, match.MatchRefereeID, match.MatchLinesman1ID, match.MatchLinesman2ID, match.MatchFourthRefereeID, tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a match
func UpdateAMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.UpdateAMatch).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// get match from body
	var match models.Match
	_ = json.NewDecoder(r.Body).Decode(&match)


	// check both teams participates in this tournament or not
	if !teamExistsInATournament(tournamentId, match.Team1DeptCode) {
		json.NewEncoder(w).Encode("Team1 doesn't participate!")
		return
	}
	if !teamExistsInATournament(tournamentId, match.Team2DeptCode) {
		json.NewEncoder(w).Encode("Team2 doesn't participate!")
		return
	}

	// check if all referee exists
	if !refereeExists(match.MatchRefereeID) {
		json.NewEncoder(w).Encode("Match referee doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman1ID) {
		json.NewEncoder(w).Encode("Match linesman1 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman2ID) {
		json.NewEncoder(w).Encode("Match linesman2 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchFourthRefereeID) {
		json.NewEncoder(w).Encode("Match fourth referee doesn't exist!")
		return
	}

	// check if the winner team is one of the two teams
	if match.WinnerTeamDeptCode != match.Team1DeptCode && match.WinnerTeamDeptCode != match.Team2DeptCode {
		json.NewEncoder(w).Encode("Winner team is not one of the two teams!")
		return
	}

	updateAMatch(tournamentId, matchId, match)

	json.NewEncoder(w).Encode(match)
}