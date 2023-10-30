package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ftms/models"
	"strconv"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var host = "http://localhost:5050"

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

	// Ping the database to ensure the connection is valid.
	if err := db.Ping(); err != nil {
		fmt.Printf("Could not connect to the database: %v", err)
	}

	//defer db.Close()
	fmt.Println("Successfully connected to mysql database")
}

// credentials for login
var jwtKey = []byte("secret_key")

// valid token check from cookies
func isTokenValid(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("jwtToken")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			//json.NewEncoder(w).Encode("No cookie found!")
			return false
		}

		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	tokenString := cookie.Value

	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			//json.NewEncoder(w).Encode("Invalid token!")
			return false
		}

		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		//json.NewEncoder(w).Encode("Invalid token!")
		return false
	}

	return true
}

// controller function to check if token is valid or not
func IsTokenValid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	if isTokenValid(w, r) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Valid token"})
		return
	} else {
		//w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Valid token"})
		return
	}
}

// generate token for login
func generateToken(w http.ResponseWriter, r *http.Request, userEmail string) bool {
	// json web token
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &models.Claims{
		Email: userEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		// set response header as internal server error
		w.WriteHeader(http.StatusInternalServerError)
		//json.NewEncoder(w).Encode("Couldn't generate token string!")
		return false
	}

	http.SetCookie(w, &http.Cookie{
		// set cookie name
		Name: "jwtToken",
		// set cookie value
		Value: tokenString,
		// set cookie expiration time
		Expires: expirationTime,
		// access from all urls
		Path: "/", // root directory
	})

	return true
}

// controller function to generate token
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// url is /api/token/generate/{userEmail}
	// so we need to get userEmail from url
	params := mux.Vars(r)
	userEmail := params["userEmail"]

	if generateToken(w, r, userEmail) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "token generated"})
		return
	} else {
		//w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "token not generated"})
		return
	}
}

// check if the user is valid or not from database
func checkUser(userEmail string, password string) bool {
	//db, _ := sql.Open("mysql", "root:@tcp(localhost:3306)/ftms")
	//defer db.Close()
	var operator models.Operator
	err := db.QueryRow("SELECT * FROM tbloperator WHERE email = ?", userEmail).Scan(&operator.Email, &operator.Password, &operator.Name, &operator.Office)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	if operator.Password != password {
		return false
	}

	return true
}

// login function for admin
func Login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("1 login successful!")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var operator models.Operator
	_ = json.NewDecoder(r.Body).Decode(&operator)

	// null check
	if operator.Email == "" || operator.Password == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		//fmt.Println("2 login successful!")
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check credentials
	if checkUser(operator.Email, operator.Password) == false {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		//fmt.Println("3 login successful!")
		json.NewEncoder(w).Encode("Invalid credentials!")
		return
	}

	// set response header as ok
	w.WriteHeader(http.StatusOK)
	//fmt.Println("4 login successful!")
	json.NewEncoder(w).Encode("Login successful!")
}

// verifications.

// check if player exists in database
func playerExists(playerRegNo int) bool {
	var player models.Player
	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode, &player.PlayerJerseyNo)

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
	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManagerEmail, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19], &team.IsKnockedOut)

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
	err := db.QueryRow("SELECT * FROM tbltournament WHERE tournamentId = ?", tournamentId).Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.StartingDate, &tournament.EndingDate)

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
	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tournamentId, matchId).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)

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
	err := db.QueryRow("SELECT * FROM tbltiebreaker WHERE tournamentId = ? AND matchID = ?", tournamentId, matchId).Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)

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
	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManagerEmail, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19], &team.IsKnockedOut)

	if err != nil {
		return false
	}

	return true
}

// check if the team is playing in a match or not
func teamIsPlayingInAMatchOfATournament(tournamentId string, matchId string, deptCode int) bool {
	var match models.Match
	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchID = ? AND (team1DeptCode = ? OR team2DeptCode = ?)", tournamentId, matchId, deptCode, deptCode).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)

	if err != nil {
		return false
	}

	return true
}

// check player is playing in a match of a tournament or not
func playerIsPlayingInAMatchOfATournament(tournamentId string, matchId string, playerRegNo int) bool {
	var team1DeptCode int
	var team2DeptCode int

	// get team1DeptCode and team2DeptCode from match
	query := "SELECT team1DeptCode, team2DeptCode FROM tblmatch WHERE tournamentId = ? AND matchID = ?"
	result, err := db.Query(query, tournamentId, matchId)

	if err != nil {
		return false
	}

	for result.Next() {
		err = result.Scan(&team1DeptCode, &team2DeptCode)

		if err != nil {
			return false
		}
	}

	// check if player is in team1 or team2
	if playerIsInATeamOfATournament(tournamentId, team1DeptCode, playerRegNo) || playerIsInATeamOfATournament(tournamentId, team2DeptCode, playerRegNo) {
		return true
	}

	return false
}

// check player is in a team of a tournament or not
func playerIsInATeamOfATournament(tournamentId string, deptCode int, playerRegNo int) bool {
	var playerRegNoFromDB [20]int

	// get playerRegNo from team
	result, err := db.Query("SELECT player1RegNo, player2RegNo, player3RegNo, player4RegNo, player5RegNo, player6RegNo, player7RegNo, player8RegNo, player9RegNo, player10RegNo, player11RegNo, player12RegNo, player13RegNo, player14RegNo, player15RegNo, player16RegNo, player17RegNo, player18RegNo, player19RegNo, player20RegNo FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&playerRegNoFromDB[0], &playerRegNoFromDB[1], &playerRegNoFromDB[2], &playerRegNoFromDB[3], &playerRegNoFromDB[4], &playerRegNoFromDB[5], &playerRegNoFromDB[6], &playerRegNoFromDB[7], &playerRegNoFromDB[8], &playerRegNoFromDB[9], &playerRegNoFromDB[10], &playerRegNoFromDB[11], &playerRegNoFromDB[12], &playerRegNoFromDB[13], &playerRegNoFromDB[14], &playerRegNoFromDB[15], &playerRegNoFromDB[16], &playerRegNoFromDB[17], &playerRegNoFromDB[18], &playerRegNoFromDB[19])

		if err != nil {
			panic(err.Error())
		}

		// check if playerRegNo is in playerRegNoFromDB
		for i := 0; i < 20; i++ {
			if playerRegNo == playerRegNoFromDB[i] {
				return true
			}
		}
	}

	return false
}

// insert operations

// insert a teacher into database
func insertNewTeacher(teacher models.Teacher) {
	// teacher.TeacherID is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblteacher VALUES (?, ?, ?, ?)", teacher.Email, teacher.Name, teacher.DeptCode, teacher.Title)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new teacher
func InsertNewTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var teacher models.Teacher
	_ = json.NewDecoder(r.Body).Decode(&teacher)

	// null check
	if teacher.Email == "" || teacher.Name == "" || teacher.DeptCode == 0 || teacher.Title == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if teacher already exists
	if teacherExists(teacher.Email) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Teacher already exists!")
		return
	}

	// check if dept exists
	if !deptExists(teacher.DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// insert new teacher
	insertNewTeacher(teacher)
	json.NewEncoder(w).Encode(teacher)
}

// teacher exists in database or not
func teacherExists(teacherEmail string) bool {
	var teacher models.Teacher
	err := db.QueryRow("SELECT * FROM tblteacher WHERE email = ?", teacherEmail).Scan(&teacher.Email, &teacher.Name, &teacher.DeptCode, &teacher.Title)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}




// insert a teamManager into database
func insertNewTeamManager(teamManager models.TeamManager) {
	// teamManager.TeamManagerRegNo is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblteammanager VALUES (?, ?)", teamManager.Email, teamManager.TournamentId)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new teamManager
func InsertNewTeamManager(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var teamManager models.TeamManager
	_ = json.NewDecoder(r.Body).Decode(&teamManager)

	// null check
	if teamManager.Email == "" || teamManager.TournamentId == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if teamManager already exists
	if teamManagerExists(teamManager.Email) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team manager already exists!")
		return
	}

	// insert new teamManager
	insertNewTeamManager(teamManager)
	json.NewEncoder(w).Encode(teamManager)
}

// teamManager exists in database or not
func teamManagerExists(teamManagerEmail string) bool {
	var teamManager models.TeamManager
	err := db.QueryRow("SELECT * FROM tblteammanager WHERE email = ?", teamManagerEmail).Scan(&teamManager.Email, &teamManager.TournamentId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err.Error())
		}
	}

	return true
}





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

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var dept models.Dept
	_ = json.NewDecoder(r.Body).Decode(&dept)

	// null check
	if dept.DeptCode == 0 || dept.DeptName == "" || dept.DeptHeadName == "" || dept.DeptShortName == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// check if dept already exists
	if deptExists(dept.DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Dept already exists!")
		return
	}

	// insert new dept
	insertNewDept(dept)
	json.NewEncoder(w).Encode(dept)
}

// insert player info into database
func insertNewPlayer(player models.Player) {
	// player.PlayerRegNo is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tblplayer VALUES (?, ?, ?, ?, ?, ?)", player.PlayerRegNo, player.PlayerSession, player.PlayerSemester, player.PlayerName, player.PlayerDeptCode, player.PlayerJerseyNo)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new player
func InsertNewPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	// null check
	if player.PlayerRegNo == 0 || player.PlayerSession == "" || player.PlayerSemester == 0 || player.PlayerName == "" || player.PlayerDeptCode == 0 || player.PlayerJerseyNo == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if player already exists
	if playerExists(player.PlayerRegNo) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player already exists!")
		return
	}

	// check if dept exists
	if !deptExists(player.PlayerDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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
	insert, err := db.Query("INSERT INTO tblteam VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", team.TournamentId, team.TeamSubmissionDate, team.DeptCode, team.TeamManagerEmail, team.TeamCaptainRegID, team.PlayerRegNo[0], team.PlayerRegNo[1], team.PlayerRegNo[2], team.PlayerRegNo[3], team.PlayerRegNo[4], team.PlayerRegNo[5], team.PlayerRegNo[6], team.PlayerRegNo[7], team.PlayerRegNo[8], team.PlayerRegNo[9], team.PlayerRegNo[10], team.PlayerRegNo[11], team.PlayerRegNo[12], team.PlayerRegNo[13], team.PlayerRegNo[14], team.PlayerRegNo[15], team.PlayerRegNo[16], team.PlayerRegNo[17], team.PlayerRegNo[18], team.PlayerRegNo[19], team.IsKnockedOut)

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
}

// return player's dept code from tblplayer in database
func getPlayerDeptCode(playerRegNo int) int {
	var player models.Player
	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode, &player.PlayerJerseyNo)

	if err != nil {
		panic(err.Error())
	}

	return player.PlayerDeptCode
}

// controller function to insert new team
func InsertNewTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var team models.Team
	_ = json.NewDecoder(r.Body).Decode(&team)

	// null check
	if team.TournamentId == "" || team.TeamSubmissionDate == "" || team.DeptCode == 0 || team.TeamManagerEmail == "" || team.TeamCaptainRegID == 0 || team.PlayerRegNo[0] == 0 || team.PlayerRegNo[1] == 0 || team.PlayerRegNo[2] == 0 || team.PlayerRegNo[3] == 0 || team.PlayerRegNo[4] == 0 || team.PlayerRegNo[5] == 0 || team.PlayerRegNo[6] == 0 || team.PlayerRegNo[7] == 0 || team.PlayerRegNo[8] == 0 || team.PlayerRegNo[9] == 0 || team.PlayerRegNo[10] == 0 || team.PlayerRegNo[11] == 0 || team.PlayerRegNo[12] == 0 || team.PlayerRegNo[13] == 0 || team.PlayerRegNo[14] == 0 || team.PlayerRegNo[15] == 0 || team.PlayerRegNo[16] == 0 || team.PlayerRegNo[17] == 0 || team.PlayerRegNo[18] == 0 || team.PlayerRegNo[19] == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if team already exists
	if teamExists(team.TournamentId, team.DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team already exists!")
		return
	}

	// check if tournament exists
	if !tournamentExists(team.TournamentId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if dept exists
	if !deptExists(team.DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// check if all players exist
	if !playerExists(team.TeamCaptainRegID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team captain doesn't exist!")
		return
	}
	for i := 0; i < 20; i++ {
		if !playerExists(team.PlayerRegNo[i]) {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " doesn't exist!")
			return
		}
	}

	// check if all players are from same dept.
	if getPlayerDeptCode(team.TeamCaptainRegID) != team.DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team captain is not from this dept!")
		return
	}
	for i := 0; i < 20; i++ {
		if getPlayerDeptCode(team.PlayerRegNo[i]) != team.DeptCode {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team captain is not in player list!")
		return
	}

	// check if player list has duplicate players
	for i := 0; i < 20-1; i++ {
		for j := i + 1; j < 20; j++ {
			if team.PlayerRegNo[i] == team.PlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " and Player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check if team manager exists
	if !teamManagerExists(team.TeamManagerEmail) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team manager doesn't exist!")
		return
	}

	// insert new team
	insertNewTeam(team)
	json.NewEncoder(w).Encode(team)
}

// insert tournament info into database
func insertNewTournament(tournament models.Tournament) {
	// tournament.TournamentId is int type. and it is primary key.
	insert, err := db.Query("INSERT INTO tbltournament VALUES(?, ?, ?, ?)", tournament.TournamentId, tournament.TournamentName, tournament.StartingDate, tournament.EndingDate)

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

	// null check
	if tournament.TournamentId == "" || tournament.TournamentName == "" || tournament.StartingDate == "" || tournament.EndingDate == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	if tournamentExists(tournament.TournamentId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament already exists!")
		return
	}

	// insert new tournament
	insertNewTournament(tournament)
	json.NewEncoder(w).Encode(tournament)
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

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var referee models.Referee
	_ = json.NewDecoder(r.Body).Decode(&referee)

	// null check
	if referee.RefereeID == 0 || referee.RefereeName == "" || referee.RefereeInstitute == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	if refereeExists(referee.RefereeID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Referee already exists!")
		return
	}

	insertNewReferee(referee)
	json.NewEncoder(w).Encode(referee)
}

// insert match info into database
func insertNewMatch(match models.Match) {
	// match.TournamentId is int type. and match.MatchId is int type. and both are primary key.
	insert, err := db.Query("INSERT INTO tblmatch VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", match.TournamentId, match.MatchId, match.MatchDate, match.Team1DeptCode, match.Team2DeptCode, match.Team1Score, match.Team2Score, match.WinnerTeamDeptCode, match.MatchRefereeID, match.MatchLinesman1ID, match.MatchLinesman2ID, match.MatchFourthRefereeID, match.Venue)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new match
func InsertNewMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var match models.Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	// // null check
	// if match.TournamentId == "" || match.MatchId == "" || match.MatchDate == "" || match.Team1DeptCode == 0 || match.Team2DeptCode == 0 || match.MatchRefereeID == 0 || match.MatchLinesman1ID == 0 || match.MatchLinesman2ID == 0 || match.MatchFourthRefereeID == 0 {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("All fields are required!")
	// 	return
	// }

	// check if match already exists
	if matchExists(match.TournamentId, match.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match already exists!")
		return
	}

	// check if tournament exists
	if !tournamentExists(match.TournamentId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// // check if team1 exists
	// if !teamExists(match.TournamentId, match.Team1DeptCode) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Team1 doesn't exist!")
	// 	return
	// }

	// // check if team2 exists
	// if !teamExists(match.TournamentId, match.Team2DeptCode) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Team2 doesn't exist!")
	// 	return
	// }

	// // check if team1 and team2 are different
	// if match.Team1DeptCode == match.Team2DeptCode {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Team1 and Team2 are same!")
	// 	return
	// }

	// // check if team1 and team2 are playing in the tournament or not
	// if !teamExistsInATournament(match.TournamentId, match.Team1DeptCode) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Team1 is not playing in the tournament!")
	// 	return
	// }
	// if !teamExistsInATournament(match.TournamentId, match.Team2DeptCode) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Team2 is not playing in the tournament!")
	// 	return
	// }

	// // check if referee exists
	// if !refereeExists(match.MatchRefereeID) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Referee doesn't exist!")
	// 	return
	// }
	// if !refereeExists(match.MatchLinesman1ID) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Linesman1 doesn't exist!")
	// 	return
	// }
	// if !refereeExists(match.MatchLinesman2ID) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Linesman2 doesn't exist!")
	// 	return
	// }
	// if !refereeExists(match.MatchFourthRefereeID) {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Fourth referee doesn't exist!")
	// 	return
	// }

	// // check if the winner team is one of the two teams
	// if match.WinnerTeamDeptCode != match.Team1DeptCode && match.WinnerTeamDeptCode != match.Team2DeptCode {
	// 	// set response header as forbidden
	// 	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Winner team is not one of the two teams!")
	// 	return
	// }

	// insert new match
	insertNewMatch(match)
	json.NewEncoder(w).Encode(match)
}

// insert starting eleven info into database
func insertNewStartingEleven(startingEleven models.StartingEleven) {
	// startingEleven.TournamentId is int type. and startingEleven.MatchId is int type. and startingEleven.DeptCode is int type. and all are primary key.
	insert, err := db.Query("INSERT INTO tblplaying11 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", startingEleven.TournamentId, startingEleven.MatchId, startingEleven.TeamDeptCode, startingEleven.StartingPlayerRegNo[0], startingEleven.StartingPlayerRegNo[1], startingEleven.StartingPlayerRegNo[2], startingEleven.StartingPlayerRegNo[3], startingEleven.StartingPlayerRegNo[4], startingEleven.StartingPlayerRegNo[5], startingEleven.StartingPlayerRegNo[6], startingEleven.StartingPlayerRegNo[7], startingEleven.StartingPlayerRegNo[8], startingEleven.StartingPlayerRegNo[9], startingEleven.StartingPlayerRegNo[10], startingEleven.SubstitutePlayerRegNo[0], startingEleven.SubstitutedPlayerRegNo[0], startingEleven.SubstitutePlayerRegNo[1], startingEleven.SubstitutedPlayerRegNo[1], startingEleven.SubstitutePlayerRegNo[2], startingEleven.SubstitutedPlayerRegNo[2])

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

// controller function to insert new starting eleven
func InsertNewStartingEleven(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var startingEleven models.StartingEleven
	_ = json.NewDecoder(r.Body).Decode(&startingEleven)

	// null check
	if startingEleven.TournamentId == "" || startingEleven.MatchId == "" || startingEleven.TeamDeptCode == 0 || startingEleven.StartingPlayerRegNo[0] == 0 || startingEleven.StartingPlayerRegNo[1] == 0 || startingEleven.StartingPlayerRegNo[2] == 0 || startingEleven.StartingPlayerRegNo[3] == 0 || startingEleven.StartingPlayerRegNo[4] == 0 || startingEleven.StartingPlayerRegNo[5] == 0 || startingEleven.StartingPlayerRegNo[6] == 0 || startingEleven.StartingPlayerRegNo[7] == 0 || startingEleven.StartingPlayerRegNo[8] == 0 || startingEleven.StartingPlayerRegNo[9] == 0 || startingEleven.StartingPlayerRegNo[10] == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if starting eleven already exists
	if startingElevenExists(startingEleven.TournamentId, startingEleven.MatchId, startingEleven.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Starting eleven already exists!")
		return
	}

	// check if match exists
	if !matchExists(startingEleven.TournamentId, startingEleven.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team is playing in the match or not
	if !teamIsPlayingInAMatchOfATournament(startingEleven.TournamentId, startingEleven.MatchId, startingEleven.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	// check players are from the team or not
	for i := 0; i < 11; i++ {
		if !playerIsInATeamOfATournament(startingEleven.TournamentId, startingEleven.TeamDeptCode, startingEleven.StartingPlayerRegNo[i]) {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " is not from the team!")
			return
		}
	}

	// check if substitute players are from the team or not
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutePlayerRegNo[i] != 0 {
			if !playerIsInATeamOfATournament(startingEleven.TournamentId, startingEleven.TeamDeptCode, startingEleven.SubstitutePlayerRegNo[i]) {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " is not from the team!")
				return
			}
		}
	}

	// check if for all the substitute players, there is a substituted player
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutePlayerRegNo[i] == 0 {
			if startingEleven.SubstitutedPlayerRegNo[i] != 0 {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " needed!")
				return
			}
			continue
		}
		if startingEleven.SubstitutedPlayerRegNo[i] == 0 {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Substituted player" + strconv.Itoa(i+1) + " needed!")
			return
		}
	}

	// check duplicate players
	for i := 0; i < 11; i++ {
		for j := i + 1; j < 11; j++ {
			if startingEleven.StartingPlayerRegNo[i] == startingEleven.StartingPlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " and Player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check duplicate substitute players
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if startingEleven.SubstitutePlayerRegNo[i] != 0 && startingEleven.SubstitutePlayerRegNo[i] == startingEleven.SubstitutePlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " and Substitute player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check duplicate substituted players
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if startingEleven.SubstitutedPlayerRegNo[i] != 0 && startingEleven.SubstitutedPlayerRegNo[i] == startingEleven.SubstitutedPlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substituted player" + strconv.Itoa(i+1) + " and Substituted player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check if the substitute players are from starting eleven or not
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutePlayerRegNo[i] != 0 {
			var found bool = false
			for j := 0; j < 11; j++ {
				if startingEleven.SubstitutePlayerRegNo[i] == startingEleven.StartingPlayerRegNo[j] {
					found = true
					break
				}
			}
			if found {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " is from starting eleven!")
				return
			}
		}
	}

	// check if the substitued players are from starting eleven or not
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutedPlayerRegNo[i] != 0 {
			var found bool = false
			for j := 0; j < 11; j++ {
				if startingEleven.SubstitutedPlayerRegNo[i] == startingEleven.StartingPlayerRegNo[j] {
					found = true
					break
				}
			}
			if !found {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substituted player" + strconv.Itoa(i+1) + " is not from starting eleven!")
				return
			}
		}
	}

	insertNewStartingEleven(startingEleven)
	json.NewEncoder(w).Encode(startingEleven)
}

// check if starting eleven already exists
func startingElevenExists(tournamentId string, matchId string, teamDeptCode int) bool {
	var startingEleven models.StartingEleven
	err := db.QueryRow("SELECT * FROM tblplaying11 WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ?", tournamentId, matchId, teamDeptCode).Scan(&startingEleven.TournamentId, &startingEleven.MatchId, &startingEleven.TeamDeptCode, &startingEleven.StartingPlayerRegNo[0], &startingEleven.StartingPlayerRegNo[1], &startingEleven.StartingPlayerRegNo[2], &startingEleven.StartingPlayerRegNo[3], &startingEleven.StartingPlayerRegNo[4], &startingEleven.StartingPlayerRegNo[5], &startingEleven.StartingPlayerRegNo[6], &startingEleven.StartingPlayerRegNo[7], &startingEleven.StartingPlayerRegNo[8], &startingEleven.StartingPlayerRegNo[9], &startingEleven.StartingPlayerRegNo[10], &startingEleven.SubstitutePlayerRegNo[0], &startingEleven.SubstitutedPlayerRegNo[0], &startingEleven.SubstitutePlayerRegNo[1], &startingEleven.SubstitutedPlayerRegNo[1], &startingEleven.SubstitutePlayerRegNo[2], &startingEleven.SubstitutedPlayerRegNo[2])

	if err != nil {
		return false
	}

	return true
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

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var tiebreaker models.Tiebreaker
	_ = json.NewDecoder(r.Body).Decode(&tiebreaker)

	// null check
	if tiebreaker.TournamentId == "" || tiebreaker.MatchId == "" || tiebreaker.Team1DeptCode == 0 || tiebreaker.Team2DeptCode == 0 || tiebreaker.Team1TieBreakerScore == 0 || tiebreaker.Team2TieBreakerScore == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if tournament exists
	if !tournamentExists(tiebreaker.TournamentId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(tiebreaker.TournamentId, tiebreaker.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team1 exists
	if !teamExists(tiebreaker.TournamentId, tiebreaker.Team1DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 doesn't exist!")
		return
	}

	// check if team2 exists
	if !teamExists(tiebreaker.TournamentId, tiebreaker.Team2DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team2 doesn't exist!")
		return
	}

	// check if team1 and team2 are playing in the match or not
	if !teamIsPlayingInAMatchOfATournament(tiebreaker.TournamentId, tiebreaker.MatchId, tiebreaker.Team1DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 is not playing in the match!")
		return
	}
	if !teamIsPlayingInAMatchOfATournament(tiebreaker.TournamentId, tiebreaker.MatchId, tiebreaker.Team2DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team2 is not playing in the match!")
		return
	}

	// check if team1 and team2 are different
	if tiebreaker.Team1DeptCode == tiebreaker.Team2DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 and Team2 are same!")
		return
	}

	// check if team1 matches with the team1 of the match
	var team1DeptCode int
	err := db.QueryRow("SELECT team1DeptCode FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tiebreaker.TournamentId, tiebreaker.MatchId).Scan(&team1DeptCode)
	if err != nil || tiebreaker.Team1DeptCode != team1DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 and Team2 are misplaced!")
		return
	}

	// check tie breaker eligibility
	var team1Score int
	var team2Score int
	err = db.QueryRow("SELECT team1Score, team2Score FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tiebreaker.TournamentId, tiebreaker.MatchId).Scan(&team1Score, &team2Score)
	if err != nil {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Error in getting match score!")
		return
	}
	if team1Score != team2Score {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tie breaker is not eligible for this match!")
		return
	}

	// check if tiebreaker already exists
	if tiebreakerExists(tiebreaker.TournamentId, tiebreaker.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tiebreaker already exists!")
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

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var individualScore models.IndividualScore
	_ = json.NewDecoder(r.Body).Decode(&individualScore)

	// null check
	if individualScore.TournamentId == "" || individualScore.MatchId == "" || individualScore.PlayerRegNo == 0 || individualScore.TeamDeptCode == 0 || individualScore.Goals == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if tournament exists
	if !tournamentExists(individualScore.TournamentId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(individualScore.TournamentId, individualScore.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team exists
	if !teamExists(individualScore.TournamentId, individualScore.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// check if the team is playing in the match or not
	if !teamIsPlayingInAMatchOfATournament(individualScore.TournamentId, individualScore.MatchId, individualScore.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	// // check if player is playing in the match or not
	// if !playerIsPlayingInAMatchOfATournament(individualScore.TournamentId, individualScore.MatchId, individualScore.PlayerRegNo) {
	// set response header as forbidden
	//	w.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(w).Encode("Player is not playing in the match!")
	// 	return
	// }

	// check if player is in the team or not
	if !playerIsInATeamOfATournament(individualScore.TournamentId, individualScore.TeamDeptCode, individualScore.PlayerRegNo) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not in the team!")
		return
	}

	// check if the player is a starting player or not
	if !playerIsAStartingPlayer(individualScore.TournamentId, individualScore.MatchId, individualScore.TeamDeptCode, individualScore.PlayerRegNo) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is neither a starting player nor a substitute player!")
		return
	}

	// check if individual score already exists
	if individualScoreExists(individualScore.TournamentId, individualScore.MatchId, individualScore.PlayerRegNo) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Individual score already exists!")
		return
	}

	// insert new individual score
	insertNewIndividualScore(individualScore)
	json.NewEncoder(w).Encode(individualScore)
}

// check if the player is a starting player or not
func playerIsAStartingPlayer(tournamentId string, matchId string, teamDeptCode int, playerRegNo int) bool {
	var startingEleven models.StartingEleven
	err := db.QueryRow("SELECT * FROM tblplaying11 WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ? AND startingPlayer1RegNo = ? OR startingPlayer2RegNo = ? OR startingPlayer3RegNo = ? OR startingPlayer4RegNo = ? OR startingPlayer5RegNo = ? OR startingPlayer6RegNo = ? OR startingPlayer7RegNo = ? OR startingPlayer8RegNo = ? OR startingPlayer9RegNo = ? OR startingPlayer10RegNo = ? OR startingPlayer11RegNo = ? OR substitutePlayer1RegNo = ? OR substitutePlayer2RegNo = ? OR substitutePlayer3RegNo = ?", tournamentId, matchId, teamDeptCode, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo).Scan(&startingEleven.TournamentId, &startingEleven.MatchId, &startingEleven.TeamDeptCode, &startingEleven.StartingPlayerRegNo[0], &startingEleven.StartingPlayerRegNo[1], &startingEleven.StartingPlayerRegNo[2], &startingEleven.StartingPlayerRegNo[3], &startingEleven.StartingPlayerRegNo[4], &startingEleven.StartingPlayerRegNo[5], &startingEleven.StartingPlayerRegNo[6], &startingEleven.StartingPlayerRegNo[7], &startingEleven.StartingPlayerRegNo[8], &startingEleven.StartingPlayerRegNo[9], &startingEleven.StartingPlayerRegNo[10], &startingEleven.SubstitutePlayerRegNo[0], &startingEleven.SubstitutedPlayerRegNo[0], &startingEleven.SubstitutePlayerRegNo[1], &startingEleven.SubstitutedPlayerRegNo[1], &startingEleven.SubstitutePlayerRegNo[2], &startingEleven.SubstitutedPlayerRegNo[2])

	if err != nil {
		return false
	}

	return true
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

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	var individualPunishment models.IndividualPunishment
	_ = json.NewDecoder(r.Body).Decode(&individualPunishment)

	// null check
	if individualPunishment.TournamentId == "" || individualPunishment.MatchId == "" || individualPunishment.PlayerRegNo == 0 || individualPunishment.TeamDeptCode == 0 || individualPunishment.PunishmentType == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if tournament exists
	if !tournamentExists(individualPunishment.TournamentId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// check if match exists
	if !matchExists(individualPunishment.TournamentId, individualPunishment.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team exists
	if !teamExists(individualPunishment.TournamentId, individualPunishment.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// check if the team is playing in the match or not
	if !teamIsPlayingInAMatchOfATournament(individualPunishment.TournamentId, individualPunishment.MatchId, individualPunishment.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	// // check if player is playing in the match or not
	// if !playerIsPlayingInAMatchOfATournament(individualPunishment.TournamentId, individualPunishment.MatchId, individualPunishment.PlayerRegNo) {
	// 	json.NewEncoder(w).Encode("Player is not playing in the match!")
	// 	return
	// }

	// check if player is in the team or not
	if !playerIsInATeamOfATournament(individualPunishment.TournamentId, individualPunishment.TeamDeptCode, individualPunishment.PlayerRegNo) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not in the team!")
		return
	}

	// check if individual punishment already exists
	if individualPunishmentExists(individualPunishment.TournamentId, individualPunishment.MatchId, individualPunishment.PlayerRegNo) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Individual punishment already exists!")
		return
	}

	// insert new individual punishment
	insertNewIndividualPunishment(individualPunishment)
	json.NewEncoder(w).Encode(individualPunishment)
}








// getting info from database

// get a techer
func getATeacher(email string) models.Teacher {
	var teacher models.Teacher

	err := db.QueryRow("SELECT * FROM tblteacher WHERE email = ?", email).Scan(&teacher.Email, &teacher.Name, &teacher.DeptCode, &teacher.Title)

	if err != nil {
		return models.Teacher{}
	}

	return teacher
}

// controller function to get a teacher
func GetATeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/teacher/{email}", controller.GetATeacher).Methods("GET")
	// get email from url
	params := mux.Vars(r)

	teacher := getATeacher(params["email"])

	json.NewEncoder(w).Encode(teacher)
}





// get an operator
func getAnOperator(email string) models.Operator {
	var operator models.Operator

	err := db.QueryRow("SELECT * FROM tbloperator WHERE email = ?", email).Scan(&operator.Email, &operator.Password, &operator.Name, &operator.Office)

	operator.Password = "" // don't send password

	if err != nil {
		return models.Operator{}
	}

	return operator
}

func GetAnOperator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/operator/{email}", controller.GetAnOperator).Methods("GET")
	// get email from url
	params := mux.Vars(r)

	operator := getAnOperator(params["email"])

	json.NewEncoder(w).Encode(operator)
}



// get a team manager
func getATeamManager(email string) models.TeamManager {
	var teamManager models.TeamManager

	err := db.QueryRow("SELECT * FROM tblteammanager WHERE email = ?", email).Scan(&teamManager.Email, &teamManager.TournamentId)

	if err != nil {
		return models.TeamManager{}
	}

	return teamManager
}

// controller function to get a team manager
func GetATeamManager(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/teammanager/{email}", controller.GetATeamManager).Methods("GET")
	// get email from url
	params := mux.Vars(r)

	teamManager := getATeamManager(params["email"])

	json.NewEncoder(w).Encode(teamManager)
}



// get all depts from database
func getAllDepts() []models.Dept {
	var dept models.Dept
	var depts []models.Dept

	result, err := db.Query("SELECT * FROM tbldept ORDER BY deptCode ASC")

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// sort descending order of tournamentYear
	result, err := db.Query("SELECT * FROM tbltournament ORDER BY endingDate ASC, startingDate ASC")

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.StartingDate, &tournament.EndingDate)

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

	err := db.QueryRow("SELECT * FROM tbltournament WHERE tournamentId = ?", tournamentId).Scan(&tournament.TournamentId, &tournament.TournamentName, &tournament.StartingDate, &tournament.EndingDate)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// ascending order by deptCode
	result, err := db.Query("SELECT * FROM tblteam WHERE tournamentId = ? ORDER BY deptCode ASC", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManagerEmail, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19], &team.IsKnockedOut)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	result, err := db.Query("SELECT * FROM tblplayer WHERE playerDeptCode = ? ORDER BY playerRegNo DESC", deptCode)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode, &player.PlayerJerseyNo)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	err := db.QueryRow("SELECT * FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode).Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManagerEmail, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19], &team.IsKnockedOut)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// order by matchId ascending
	result, err := db.Query("SELECT * FROM tblmatch WHERE tournamentId = ? ORDER BY matchID ASC", tournamentId)

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		err = result.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	err := db.QueryRow("SELECT * FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tournamentId, matchId).Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	var match models.Match
	match = getAMatchOfATournament(tournamentId, matchId)

	json.NewEncoder(w).Encode(match)
}

// get starting eleven of a team of a match
func getStartingElevenOfATeamOfAMatch(tournamentId string, matchId string, deptCode int) models.StartingEleven {
	var startingEleven models.StartingEleven

	err := db.QueryRow("SELECT * FROM tblplaying11 WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ?", tournamentId, matchId, deptCode).Scan(&startingEleven.TournamentId, &startingEleven.MatchId, &startingEleven.TeamDeptCode, &startingEleven.StartingPlayerRegNo[0], &startingEleven.StartingPlayerRegNo[1], &startingEleven.StartingPlayerRegNo[2], &startingEleven.StartingPlayerRegNo[3], &startingEleven.StartingPlayerRegNo[4], &startingEleven.StartingPlayerRegNo[5], &startingEleven.StartingPlayerRegNo[6], &startingEleven.StartingPlayerRegNo[7], &startingEleven.StartingPlayerRegNo[8], &startingEleven.StartingPlayerRegNo[9], &startingEleven.StartingPlayerRegNo[10], &startingEleven.SubstitutePlayerRegNo[0], &startingEleven.SubstitutedPlayerRegNo[0], &startingEleven.SubstitutePlayerRegNo[1], &startingEleven.SubstitutedPlayerRegNo[1], &startingEleven.SubstitutePlayerRegNo[2], &startingEleven.SubstitutedPlayerRegNo[2])

	if err != nil {
		//panic(err.Error())
		return models.StartingEleven{}
	}

	return startingEleven
}

// controller function to get starting eleven of a team of a match
func GetStartingElevenOfATeamOfAMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")

	// router.HandleFunc("/api/match/startingeleven/{tournamentId}/{matchId}/{deptCode}", controller.GetStartingElevenOfATeamOfAMatch).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and deptCode from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	deptCode, _ := strconv.Atoi(params["deptCode"])

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// team exists or not
	if !teamExists(tournamentId, deptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// team is playing in the match or not
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, deptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	var startingEleven models.StartingEleven
	startingEleven = getStartingElevenOfATeamOfAMatch(tournamentId, matchId, deptCode)

	json.NewEncoder(w).Encode(startingEleven)
}

// get a player
func getAPlayer(playerRegNo int) models.Player {
	var player models.Player

	err := db.QueryRow("SELECT * FROM tblplayer WHERE playerRegNo = ?", playerRegNo).Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode, &player.PlayerJerseyNo)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// order by refereeID ascending
	result, err := db.Query("SELECT * FROM tblreferee ORDER BY refereeID ASC")

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// order by matchID ascending
	result, err := db.Query("SELECT * FROM tbltiebreaker WHERE tournamentId = ? ORDER BY matchID ASC", tournamentId)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// order by playerRegNo descending
	result, err := db.Query("SELECT * FROM tblindividualscore WHERE tournamentId = ? ORDER BY playerRegNo DESC", tournamentId)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var individualScores []models.IndividualScore
	individualScores = getAllIndividualScoresOfATournament(id)

	json.NewEncoder(w).Encode(individualScores)
}

// get all individual scores of a match by a team
func getAllIndividualScoresOfAMatchByATeam(tournamentId string, matchId string, teamDeptCode int) []models.IndividualScore {
	var individualScore models.IndividualScore
	var individualScores []models.IndividualScore

	// order by playerRegNo descending
	result, err := db.Query("SELECT * FROM tblindividualscore WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ? ORDER BY playerRegNo DESC", tournamentId, matchId, teamDeptCode)

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

// controller function to get all individual scores of a match by a team
func GetAllIndividualScoresOfAMatchByATeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/match/team/individualscores/{tournamentId}/{matchId}/{teamDeptCode}", controller.GetAllIndividualScoresOfAMatchByATeam).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and teamDeptCode from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	teamDeptCode, _ := params["teamDeptCode"]

	// convert teamDeptCode from string to int
	teamDeptCodeInt, err := strconv.Atoi(teamDeptCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if the team is playing in the match
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, teamDeptCodeInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	var individualScores []models.IndividualScore
	individualScores = getAllIndividualScoresOfAMatchByATeam(tournamentId, matchId, teamDeptCodeInt)

	json.NewEncoder(w).Encode(individualScores)
}

// get all individual scores of a player in a tournament
func getAllIndividualScoresOfAPlayerInATournament(tournamentId string, playerRegNo int) []models.IndividualScore {
	var individualScore models.IndividualScore
	var individualScores []models.IndividualScore

	// order by matchID ascending
	result, err := db.Query("SELECT * FROM tblindividualscore WHERE tournamentId = ? AND playerRegNo = ? ORDER BY matchID ASC", tournamentId, playerRegNo)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// player exists or not
	if !playerExists(playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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

	// order by teamDeptCode ascending and playerRegNo descending
	result, err := db.Query("SELECT * FROM tblindividualpunishment WHERE tournamentId = ? ORDER BY teamDeptCode ASC, playerRegNo DESC", tournamentId)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	var individualPunishments []models.IndividualPunishment
	individualPunishments = getAllIndividualPunishmentsOfATournament(id)

	json.NewEncoder(w).Encode(individualPunishments)
}

// get all individual punishments of a match by a team
func getAllIndividualPunishmentsOfAMatchByATeam(tournamentId string, matchId string, teamDeptCode int) []models.IndividualPunishment {
	var individualPunishment models.IndividualPunishment
	var individualPunishments []models.IndividualPunishment

	// order by playerRegNo descending
	result, err := db.Query("SELECT * FROM tblindividualpunishment WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ? ORDER BY playerRegNo DESC", tournamentId, matchId, teamDeptCode)

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

// controller function to get all individual punishments of a match by a team
func GetAllIndividualPunishmentsOfAMatchByATeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// router.HandleFunc("/api/tournament/match/team/individualpunishments/{tournamentId}/{matchId}/{teamDeptCode}", controller.GetAllIndividualPunishmentsOfAMatchByATeam).Methods("GET")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and teamDeptCode from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	teamDeptCode, _ := params["teamDeptCode"]

	// convert teamDeptCode from string to int
	teamDeptCodeInt, err := strconv.Atoi(teamDeptCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if the team is playing in the match
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, teamDeptCodeInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	var individualPunishments []models.IndividualPunishment
	individualPunishments = getAllIndividualPunishmentsOfAMatchByATeam(tournamentId, matchId, teamDeptCodeInt)

	json.NewEncoder(w).Encode(individualPunishments)
}

// get all individual punishments of a player in a tournament
func getAllIndividualPunishmentsOfAPlayerInATournament(tournamentId string, playerRegNo int) []models.IndividualPunishment {
	var individualPunishment models.IndividualPunishment
	var individualPunishments []models.IndividualPunishment

	// descending order of playerRegNo
	result, err := db.Query("SELECT * FROM tblindividualpunishment WHERE tournamentId = ? AND playerRegNo = ? ORDER BY playerRegNo DESC", tournamentId, playerRegNo)

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// player exists or not
	if !playerExists(playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
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
	query := "UPDATE tbltournament SET tournamentName = ?, startingDate = ?, endingDate = ? WHERE tournamentId = ?"

	_, err := db.Exec(query, tournament.TournamentName, tournament.StartingDate, tournament.EndingDate, tournamentId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a tournament
func UpdateATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/tournament/{tournamentId}", controller.UpdateATournament).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// get tournament from body
	var tournament models.Tournament
	_ = json.NewDecoder(r.Body).Decode(&tournament)

	// null value check
	if tournament.TournamentName == "" || tournament.StartingDate == "" || tournament.EndingDate == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// tournamentId can't be changed
	if id != tournament.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}

	updateATournament(id, tournament)

	json.NewEncoder(w).Encode(tournament)
}

// update a player
func updateAPlayer(playerRegNo int, player models.Player) {
	query := "UPDATE tblplayer SET playerSession = ?, playerSemester = ?, playerName = ?, playerDeptCode = ?, playerJerseyNo = ? WHERE playerRegNo = ?"

	_, err := db.Exec(query, player.PlayerSession, player.PlayerSemester, player.PlayerName, player.PlayerDeptCode, player.PlayerJerseyNo, playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a player
func UpdateAPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// get player from body
	var player models.Player
	_ = json.NewDecoder(r.Body).Decode(&player)

	// null value check
	if player.PlayerSession == "" || player.PlayerSemester == 0 || player.PlayerName == "" || player.PlayerDeptCode == 0 || player.PlayerJerseyNo == 0 {
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the playerRegNo is changed
	if id != player.PlayerRegNo {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("PlayerRegNo can't be changed!")
		return
	}

	// check if dept is changed
	if player.PlayerDeptCode != getAPlayer(id).PlayerDeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("PlayerDeptCode can't be changed!")
		return
	}

	updateAPlayer(id, player)

	json.NewEncoder(w).Encode(player)
}

// update a dept
func updateADept(deptCode int, dept models.Dept) {
	query := "UPDATE tbldept SET deptName = ?, deptShortName = ?, deptHeadName = ? WHERE deptCode = ?"

	_, err := db.Exec(query, dept.DeptName, dept.DeptShortName, dept.DeptHeadName, deptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a dept
func UpdateADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// get dept from body
	var dept models.Dept
	_ = json.NewDecoder(r.Body).Decode(&dept)

	// null value check
	if dept.DeptName == "" || dept.DeptShortName == "" || dept.DeptHeadName == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the deptCode is changed
	if id != dept.DeptCode {
		json.NewEncoder(w).Encode("DeptCode can't be changed!")
		return
	}

	updateADept(id, dept)

	json.NewEncoder(w).Encode(dept)
}

// update a team
func updateATeam(tournamentId string, deptCode int, team models.Team) {
	query := "UPDATE tblteam SET teamSubmissionDate = ?, teamManagerEmail = ?, teamCaptainRegID = ?, player1RegNo = ?, player2RegNo = ?, player3RegNo = ?, player4RegNo = ?, player5RegNo = ?, player6RegNo = ?, player7RegNo = ?, player8RegNo = ?, player9RegNo = ?, player10RegNo = ?, player11RegNo = ?, player12RegNo = ?, player13RegNo = ?, player14RegNo = ?, player15RegNo = ?, player16RegNo = ?, player17RegNo = ?, player18RegNo = ?, player19RegNo = ?, player20RegNo = ?, isKnockedOut = ? WHERE tournamentId = ? AND deptCode = ?"

	_, err := db.Exec(query, team.TeamSubmissionDate, team.TeamManagerEmail, team.TeamCaptainRegID, team.PlayerRegNo[0], team.PlayerRegNo[1], team.PlayerRegNo[2], team.PlayerRegNo[3], team.PlayerRegNo[4], team.PlayerRegNo[5], team.PlayerRegNo[6], team.PlayerRegNo[7], team.PlayerRegNo[8], team.PlayerRegNo[9], team.PlayerRegNo[10], team.PlayerRegNo[11], team.PlayerRegNo[12], team.PlayerRegNo[13], team.PlayerRegNo[14], team.PlayerRegNo[15], team.PlayerRegNo[16], team.PlayerRegNo[17], team.PlayerRegNo[18], team.PlayerRegNo[19], team.IsKnockedOut, tournamentId, deptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a team
func UpdateATeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// get team from body
	var team models.Team
	_ = json.NewDecoder(r.Body).Decode(&team)

	// null value check
	if team.TeamSubmissionDate == "" || team.TeamManagerEmail == "" || team.TeamCaptainRegID == 0 || team.PlayerRegNo[0] == 0 || team.PlayerRegNo[1] == 0 || team.PlayerRegNo[2] == 0 || team.PlayerRegNo[3] == 0 || team.PlayerRegNo[4] == 0 || team.PlayerRegNo[5] == 0 || team.PlayerRegNo[6] == 0 || team.PlayerRegNo[7] == 0 || team.PlayerRegNo[8] == 0 || team.PlayerRegNo[9] == 0 || team.PlayerRegNo[10] == 0 || team.PlayerRegNo[11] == 0 || team.PlayerRegNo[12] == 0 || team.PlayerRegNo[13] == 0 || team.PlayerRegNo[14] == 0 || team.PlayerRegNo[15] == 0 || team.PlayerRegNo[16] == 0 || team.PlayerRegNo[17] == 0 || team.PlayerRegNo[18] == 0 || team.PlayerRegNo[19] == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the tournamentId and deptCode is changed
	if tournamentId != team.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}

	if deptCodeInt != team.DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("DeptCode can't be changed!")
		return
	}

	// check if all players exist
	if !playerExists(team.TeamCaptainRegID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team captain doesn't exist!")
		return
	}
	for i := 0; i < 20; i++ {
		if !playerExists(team.PlayerRegNo[i]) {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " doesn't exist!")
			return
		}
	}

	// check if all players are from same dept.
	if getPlayerDeptCode(team.TeamCaptainRegID) != team.DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team captain is not from this dept!")
		return
	}
	for i := 0; i < 20; i++ {
		if getPlayerDeptCode(team.PlayerRegNo[i]) != team.DeptCode {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team captain is not in player list!")
		return
	}

	// check if player list has duplicate players
	for i := 0; i < 20-1; i++ {
		for j := i + 1; j < 20; j++ {
			if team.PlayerRegNo[i] == team.PlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " and Player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check if team manager exists
	if !teamManagerExists(team.TeamManagerEmail) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team manager doesn't exist!")
		return
	}

	updateATeam(tournamentId, deptCodeInt, team)

	json.NewEncoder(w).Encode(team)
}

// update a match
func updateAMatch(tournamentId string, matchId string, match models.Match) {
	query := "UPDATE tblmatch SET matchDate = ?, team1DeptCode = ?, team2DeptCode = ?, team1Score = ?, team2Score = ?, winnerTeamDeptCode = ?, matchRefereeID = ?, matchLineman1ID = ?, matchLineman2ID = ?, matchFourthRefereeID = ?, venue = ? WHERE tournamentId = ? AND matchID = ?"

	_, err := db.Exec(query, match.MatchDate, match.Team1DeptCode, match.Team2DeptCode, match.Team1Score, match.Team2Score, match.WinnerTeamDeptCode, match.MatchRefereeID, match.MatchLinesman1ID, match.MatchLinesman2ID, match.MatchFourthRefereeID, match.Venue, tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a match
func UpdateAMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.UpdateAMatch).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// get match from body
	var match models.Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	// null value check
	if match.MatchDate == "" || match.Team1DeptCode == 0 || match.Team2DeptCode == 0 || match.MatchRefereeID == 0 || match.MatchLinesman1ID == 0 || match.MatchLinesman2ID == 0 || match.MatchFourthRefereeID == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the tournamentId and matchId is changed
	if tournamentId != match.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}

	if matchId != match.MatchId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("MatchId can't be changed!")
		return
	}

	// check both teams participates in this tournament or not
	if !teamExistsInATournament(tournamentId, match.Team1DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 doesn't participate in this tournament!")
		return
	}
	if !teamExistsInATournament(tournamentId, match.Team2DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team2 doesn't participate in this tournament!")
		return
	}

	// check if all referee exists
	if !refereeExists(match.MatchRefereeID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match referee doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman1ID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match linesman1 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchLinesman2ID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match linesman2 doesn't exist!")
		return
	}
	if !refereeExists(match.MatchFourthRefereeID) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match fourth referee doesn't exist!")
		return
	}

	// check if the winner team is one of the two teams
	if match.WinnerTeamDeptCode != 0 {
		if match.WinnerTeamDeptCode != match.Team1DeptCode && match.WinnerTeamDeptCode != match.Team2DeptCode {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Winner team is not one of the two teams!")
			return
		}
	}

	// check if both teams deptCode is same or not
	if match.Team1DeptCode == match.Team2DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Both teams are from same dept!")
		return
	}

	updateAMatch(tournamentId, matchId, match)

	json.NewEncoder(w).Encode(match)
}

// update a starting eleven
func updateAStartingEleven(tournamentId string, matchId string, teamDeptCode int, startingEleven models.StartingEleven) {
	query := "UPDATE tblplaying11 SET startingPlayer1RegNo = ?, startingPlayer2RegNo = ?, startingPlayer3RegNo = ?, startingPlayer4RegNo = ?, startingPlayer5RegNo = ?, startingPlayer6RegNo = ?, startingPlayer7RegNo = ?, startingPlayer8RegNo = ?, startingPlayer9RegNo = ?, startingPlayer10RegNo = ?, startingPlayer11RegNo = ?, substitutePlayer1RegNo = ?, substitutedPlayer1RegNo = ?, substitutePlayer2RegNo = ?, substitutedPlayer2RegNo = ?, substitutePlayer3RegNo = ?, substitutedPlayer3RegNo = ? WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ?"

	_, err := db.Exec(query, startingEleven.StartingPlayerRegNo[0], startingEleven.StartingPlayerRegNo[1], startingEleven.StartingPlayerRegNo[2], startingEleven.StartingPlayerRegNo[3], startingEleven.StartingPlayerRegNo[4], startingEleven.StartingPlayerRegNo[5], startingEleven.StartingPlayerRegNo[6], startingEleven.StartingPlayerRegNo[7], startingEleven.StartingPlayerRegNo[8], startingEleven.StartingPlayerRegNo[9], startingEleven.StartingPlayerRegNo[10], startingEleven.SubstitutePlayerRegNo[0], startingEleven.SubstitutedPlayerRegNo[0], startingEleven.SubstitutePlayerRegNo[1], startingEleven.SubstitutedPlayerRegNo[1], startingEleven.SubstitutePlayerRegNo[2], startingEleven.SubstitutedPlayerRegNo[2], tournamentId, matchId, teamDeptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a starting eleven
func UpdateAStartingEleven(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}", controller.UpdateAStartingEleven).Methods("PUT")

	params := mux.Vars(r)

	// get tournamentId, matchId and teamDeptCode from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	teamDeptCode, _ := params["teamDeptCode"]

	// convert teamDeptCode from string to int
	teamDeptCodeInt, err := strconv.Atoi(teamDeptCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// get starting eleven from body
	var startingEleven models.StartingEleven
	_ = json.NewDecoder(r.Body).Decode(&startingEleven)

	// tournamentId, matchId and teamDeptCode can't be changed
	if tournamentId != startingEleven.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}
	if matchId != startingEleven.MatchId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("MatchId can't be changed!")
		return
	}
	if teamDeptCodeInt != startingEleven.TeamDeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TeamDeptCode can't be changed!")
		return
	}

	// null check
	if startingEleven.TournamentId == "" || startingEleven.MatchId == "" || startingEleven.TeamDeptCode == 0 || startingEleven.StartingPlayerRegNo[0] == 0 || startingEleven.StartingPlayerRegNo[1] == 0 || startingEleven.StartingPlayerRegNo[2] == 0 || startingEleven.StartingPlayerRegNo[3] == 0 || startingEleven.StartingPlayerRegNo[4] == 0 || startingEleven.StartingPlayerRegNo[5] == 0 || startingEleven.StartingPlayerRegNo[6] == 0 || startingEleven.StartingPlayerRegNo[7] == 0 || startingEleven.StartingPlayerRegNo[8] == 0 || startingEleven.StartingPlayerRegNo[9] == 0 || startingEleven.StartingPlayerRegNo[10] == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields are required!")
		return
	}

	// check if match exists
	if !matchExists(startingEleven.TournamentId, startingEleven.MatchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// check if team is playing in the match or not
	if !teamIsPlayingInAMatchOfATournament(startingEleven.TournamentId, startingEleven.MatchId, startingEleven.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	// check players are from the team or not
	for i := 0; i < 11; i++ {
		if !playerIsInATeamOfATournament(startingEleven.TournamentId, startingEleven.TeamDeptCode, startingEleven.StartingPlayerRegNo[i]) {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " is not from the team!")
			return
		}
	}

	// check if substitute players are from the team or not
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutePlayerRegNo[i] != 0 {
			if !playerIsInATeamOfATournament(startingEleven.TournamentId, startingEleven.TeamDeptCode, startingEleven.SubstitutePlayerRegNo[i]) {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " is not from the team!")
				return
			}
		}
	}

	// check if for all the substitute players, there is a substituted player
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutePlayerRegNo[i] == 0 {
			if startingEleven.SubstitutedPlayerRegNo[i] != 0 {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " needed!")
				return
			}
			continue
		}
		if startingEleven.SubstitutedPlayerRegNo[i] == 0 {
			// set response header as forbidden
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Substituted player" + strconv.Itoa(i+1) + " needed!")
			return
		}
	}

	// check duplicate players
	for i := 0; i < 11; i++ {
		for j := i + 1; j < 11; j++ {
			if startingEleven.StartingPlayerRegNo[i] == startingEleven.StartingPlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Player" + strconv.Itoa(i+1) + " and Player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check duplicate substitute players
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if startingEleven.SubstitutePlayerRegNo[i] != 0 && startingEleven.SubstitutePlayerRegNo[i] == startingEleven.SubstitutePlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " and Substitute player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check duplicate substituted players
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if startingEleven.SubstitutedPlayerRegNo[i] != 0 && startingEleven.SubstitutedPlayerRegNo[i] == startingEleven.SubstitutedPlayerRegNo[j] {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substituted player" + strconv.Itoa(i+1) + " and Substituted player" + strconv.Itoa(j+1) + " are same!")
				return
			}
		}
	}

	// check if the substitute players are from starting eleven or not
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutePlayerRegNo[i] != 0 {
			var found bool = false
			for j := 0; j < 11; j++ {
				if startingEleven.SubstitutePlayerRegNo[i] == startingEleven.StartingPlayerRegNo[j] {
					found = true
					break
				}
			}
			if found {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substitute player" + strconv.Itoa(i+1) + " is from starting eleven!")
				return
			}
		}
	}

	// check if the substitued players are from starting eleven or not
	for i := 0; i < 3; i++ {
		if startingEleven.SubstitutedPlayerRegNo[i] != 0 {
			var found bool = false
			for j := 0; j < 11; j++ {
				if startingEleven.SubstitutedPlayerRegNo[i] == startingEleven.StartingPlayerRegNo[j] {
					found = true
					break
				}
			}
			if !found {
				// set response header as forbidden
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode("Substituted player" + strconv.Itoa(i+1) + " is not from starting eleven!")
				return
			}
		}
	}

	updateAStartingEleven(tournamentId, matchId, teamDeptCodeInt, startingEleven)

	json.NewEncoder(w).Encode(startingEleven)
}

// update a referee
func updateAReferee(refereeId int, referee models.Referee) {
	query := "UPDATE tblreferee SET refereeName = ?, refereeInstitute = ? WHERE refereeID = ?"

	_, err := db.Exec(query, referee.RefereeName, referee.RefereeInstitute, refereeId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a referee
func UpdateAReferee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/referee/{refereeId}", controller.UpdateAReferee).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["refereeId"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// referee exists or not
	if !refereeExists(id) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Referee doesn't exist!")
		return
	}

	// get referee from body
	var referee models.Referee
	_ = json.NewDecoder(r.Body).Decode(&referee)

	// null value check
	if referee.RefereeName == "" || referee.RefereeInstitute == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the refereeId is changed
	if id != referee.RefereeID {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("RefereeId can't be changed!")
		return
	}

	updateAReferee(id, referee)

	json.NewEncoder(w).Encode(referee)
}

// update a tiebreaker
func updateATiebreaker(tournamentId string, matchId string, tiebreaker models.Tiebreaker) {
	query := "UPDATE tbltiebreaker SET team1DeptCode = ?, team2DeptCode = ?, team1TieBreakerScore = ?, team2TieBreakerScore = ? WHERE tournamentId = ? AND matchID = ?"

	_, err := db.Exec(query, tiebreaker.Team1DeptCode, tiebreaker.Team2DeptCode, tiebreaker.Team1TieBreakerScore, tiebreaker.Team2TieBreakerScore, tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update a tiebreaker
func UpdateATiebreaker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/tiebreaker/{tournamentId}/{matchId}", controller.UpdateATiebreaker).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// get tiebreaker from body
	var tiebreaker models.Tiebreaker
	_ = json.NewDecoder(r.Body).Decode(&tiebreaker)

	// null value check
	if tiebreaker.Team1DeptCode == 0 || tiebreaker.Team2DeptCode == 0 || tiebreaker.Team1TieBreakerScore == 0 || tiebreaker.Team2TieBreakerScore == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the tournamentId and matchId is changed
	if tournamentId != tiebreaker.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}

	if matchId != tiebreaker.MatchId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("MatchId can't be changed!")
		return
	}

	// check if both teams are playing in the match
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, tiebreaker.Team1DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 is not playing in the match!")
		return
	}
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, tiebreaker.Team2DeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team2 is not playing in the match!")
		return
	}

	// check if team1 and team2 are different
	if tiebreaker.Team1DeptCode == tiebreaker.Team2DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 and Team2 are same!")
		return
	}

	// check if team1 matches with the team1 of the match
	var team1DeptCode int
	err := db.QueryRow("SELECT team1_deptCode FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tiebreaker.TournamentId, tiebreaker.MatchId).Scan(&team1DeptCode)
	if err != nil || tiebreaker.Team1DeptCode != team1DeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team1 and Team2 are misplaced!")
		return
	}

	// check if tiebreaker valid or not
	if tiebreaker.Team1TieBreakerScore == tiebreaker.Team2TieBreakerScore {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tiebreaker score can not be same!")
		return
	}

	// check tie breaker eligibility
	var team1Score int
	var team2Score int
	err = db.QueryRow("SELECT team1Score, team2Score FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tiebreaker.TournamentId, tiebreaker.MatchId).Scan(&team1Score, &team2Score)
	if err != nil {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Error in getting match score!")
		return
	}
	if team1Score != team2Score {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tie breaker is not eligible for this match!")
		return
	}

	updateATiebreaker(tournamentId, matchId, tiebreaker)

	json.NewEncoder(w).Encode(tiebreaker)
}

// update an individual score
func updateAnIndividualScore(tournamentId string, matchId string, playerRegNo int, individualScore models.IndividualScore) {
	query := "UPDATE tblindividualscore SET teamDeptCode = ?, goals = ? WHERE tournamentId = ? AND matchID = ? AND playerRegNo = ?"

	_, err := db.Exec(query, individualScore.TeamDeptCode, individualScore.Goals, tournamentId, matchId, playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update an individual score
func UpdateAnIndividualScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/individualscore/{tournamentId}/{matchId}/{playerRegNo}", controller.UpdateAnIndividualScore).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and playerRegNo from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	playerRegNo, _ := params["playerRegNo"]

	// convert playerRegNo from string to int
	playerRegNoInt, err := strconv.Atoi(playerRegNo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// player playing in the match or not
	if !playerIsPlayingInAMatchOfATournament(tournamentId, matchId, playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not playing in the match!")
		return
	}

	// get individualScore from body
	var individualScore models.IndividualScore
	_ = json.NewDecoder(r.Body).Decode(&individualScore)

	// null value check
	if individualScore.TeamDeptCode == 0 || individualScore.Goals == 0 {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the tournamentId, matchId and playerRegNo is changed
	if tournamentId != individualScore.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}

	if matchId != individualScore.MatchId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("MatchId can't be changed!")
		return
	}

	if playerRegNoInt != individualScore.PlayerRegNo {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("PlayerRegNo can't be changed!")
		return
	}

	// check if the team is playing in the match
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, individualScore.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	// check if the player is from the team
	if getPlayerDeptCode(playerRegNoInt) != individualScore.TeamDeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not from the team!")
		return
	}

	updateAnIndividualScore(tournamentId, matchId, playerRegNoInt, individualScore)

	json.NewEncoder(w).Encode(individualScore)
}

// update an individual punishment
func updateAnIndividualPunishment(tournamentId string, matchId string, playerRegNo int, individualPunishment models.IndividualPunishment) {
	query := "UPDATE tblindividualpunishment SET teamDeptCode = ?, punishmentType = ? WHERE tournamentId = ? AND matchID = ? AND playerRegNo = ?"

	_, err := db.Exec(query, individualPunishment.TeamDeptCode, individualPunishment.PunishmentType, tournamentId, matchId, playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to update an individual punishment
func UpdateAnIndividualPunishment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}", controller.UpdateAnIndividualPunishment).Methods("PUT")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and playerRegNo from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	playerRegNo, _ := params["playerRegNo"]

	// convert playerRegNo from string to int
	playerRegNoInt, err := strconv.Atoi(playerRegNo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// player playing in the match or not
	if !playerIsPlayingInAMatchOfATournament(tournamentId, matchId, playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not playing in the match!")
		return
	}

	// get individualPunishment from body
	var individualPunishment models.IndividualPunishment
	_ = json.NewDecoder(r.Body).Decode(&individualPunishment)

	// null value check
	if individualPunishment.TeamDeptCode == 0 || individualPunishment.PunishmentType == "" {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("All fields must be filled!")
		return
	}

	// check if the tournamentId, matchId and playerRegNo is changed
	if tournamentId != individualPunishment.TournamentId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("TournamentId can't be changed!")
		return
	}

	if matchId != individualPunishment.MatchId {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("MatchId can't be changed!")
		return
	}

	if playerRegNoInt != individualPunishment.PlayerRegNo {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("PlayerRegNo can't be changed!")
		return
	}

	// check if the team is playing in the match
	if !teamIsPlayingInAMatchOfATournament(tournamentId, matchId, individualPunishment.TeamDeptCode) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team is not playing in the match!")
		return
	}

	// check if the player is from the team
	if getPlayerDeptCode(playerRegNoInt) != individualPunishment.TeamDeptCode {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not from the team!")
		return
	}

	updateAnIndividualPunishment(tournamentId, matchId, playerRegNoInt, individualPunishment)

	json.NewEncoder(w).Encode(individualPunishment)
}

// delete operations

// delete a tournament
func deleteATournament(tournamentId string) {
	_, err := db.Query("DELETE FROM tbltournament WHERE tournamentId = ?", tournamentId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a tournament
func DeleteATournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/tournament/{tournamentId}", controller.DeleteATournament).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	id, _ := params["tournamentId"]
	// id is string type

	// tournament exists or not
	if !tournamentExists(id) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tournament doesn't exist!")
		return
	}

	// prequisite check
	matchExistsInATournament(w, r, id)
	anyTeamExistsInATournament(w, r, id)

	deleteATournament(id)

	json.NewEncoder(w).Encode("Tournament deleted successfully!")
}

// match exists in a tournament or not
func matchExistsInATournament(w http.ResponseWriter, r *http.Request, tournamentId string) bool {
	query := "SELECT * FROM tblmatch WHERE tournamentId = ?"
	rows, err := db.Query(query, tournamentId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// delete all tiebreakers of this match using api call
		var match models.Match
		err = rows.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)
		if err != nil {
			panic(err.Error())
		}
		url := host + "http://localhost:5000/api/tournament/match/" + match.TournamentId + "/" + match.MatchId
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// any team exists in a tournament or not
func anyTeamExistsInATournament(w http.ResponseWriter, r *http.Request, tournamentId string) bool {
	query := "SELECT deptCode FROM tblteam WHERE tournamentId = ?"
	rows, err := db.Query(query, tournamentId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// delete all teams of this tournament using api call
		var deptCode int
		err = rows.Scan(&deptCode)
		if err != nil {
			panic(err.Error())
		}
		url := host + "http://localhost:5000/api/tournament/team/" + tournamentId + "/" + strconv.Itoa(deptCode)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// delete a player
func deleteAPlayer(playerRegNo int) {
	_, err := db.Query("DELETE FROM tblplayer WHERE playerRegNo = ?", playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a player
func DeleteAPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/player/{playerRegNo}", controller.DeleteAPlayer).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["playerRegNo"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// player exists or not
	if !playerExists(id) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player doesn't exist!")
		return
	}

	// prequisite check. check all tables where playerRegNo is used
	playerIsInATeam(w, r, id)
	playerHasIndividualPunishment(w, r, id)
	playerHasIndividualScore(w, r, id)
	playerExistsInAPlayingEleven(w, r, id)

	deleteAPlayer(id)

	json.NewEncoder(w).Encode("Player deleted successfully!")
}

// player is in a playing eleven or not
func playerExistsInAPlayingEleven(w http.ResponseWriter, r *http.Request, playerRegNo int) bool {
	query := "SELECT * FROM tblplaying11 WHERE startingPlayer1RegNo = ? OR startingPlayer2RegNo = ? OR startingPlayer3RegNo = ? OR startingPlayer4RegNo = ? OR startingPlayer5RegNo = ? OR startingPlayer6RegNo = ? OR startingPlayer7RegNo = ? OR startingPlayer8RegNo = ? OR startingPlayer9RegNo = ? OR startingPlayer10RegNo = ? OR startingPlayer11RegNo = ?"
	rows, err := db.Query(query, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the playing eleven
		var playingEleven models.StartingEleven
		err = rows.Scan(&playingEleven.TournamentId, &playingEleven.MatchId, &playingEleven.TeamDeptCode, &playingEleven.StartingPlayerRegNo[0], &playingEleven.StartingPlayerRegNo[1], &playingEleven.StartingPlayerRegNo[2], &playingEleven.StartingPlayerRegNo[3], &playingEleven.StartingPlayerRegNo[4], &playingEleven.StartingPlayerRegNo[5], &playingEleven.StartingPlayerRegNo[6], &playingEleven.StartingPlayerRegNo[7], &playingEleven.StartingPlayerRegNo[8], &playingEleven.StartingPlayerRegNo[9], &playingEleven.StartingPlayerRegNo[10], &playingEleven.SubstitutePlayerRegNo[0], &playingEleven.SubstitutedPlayerRegNo[0], &playingEleven.SubstitutePlayerRegNo[1], &playingEleven.SubstitutedPlayerRegNo[1], &playingEleven.SubstitutePlayerRegNo[2], &playingEleven.SubstitutedPlayerRegNo[2])

		if err != nil {
			panic(err.Error())
		}
		// now call delete playing eleven api
		url := host + "http://localhost:5000/api/match/startingeleven/" + playingEleven.TournamentId + "/" + playingEleven.MatchId + "/" + strconv.Itoa(playingEleven.TeamDeptCode)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// player is in a team or not
func playerIsInATeam(w http.ResponseWriter, r *http.Request, playerRegNo int) bool {
	query := "SELECT * FROM tblteam WHERE player1RegNo = ? OR player2RegNo = ? OR player3RegNo = ? OR player4RegNo = ? OR player5RegNo = ? OR player6RegNo = ? OR player7RegNo = ? OR player8RegNo = ? OR player9RegNo = ? OR player10RegNo = ? OR player11RegNo = ? OR player12RegNo = ? OR player13RegNo = ? OR player14RegNo = ? OR player15RegNo = ? OR player16RegNo = ? OR player17RegNo = ? OR player18RegNo = ? OR player19RegNo = ? OR player20RegNo = ?"
	rows, err := db.Query(query, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo, playerRegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and deptCode and delete the team
		var team models.Team
		err = rows.Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManagerEmail, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19], &team.IsKnockedOut)
		if err != nil {
			panic(err.Error())
		}
		// now call delete team api
		url := host + "http://localhost:5000/api/tournament/team/" + team.TournamentId + "/" + strconv.Itoa(team.DeptCode)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// player has individual punishment or not
func playerHasIndividualPunishment(w http.ResponseWriter, r *http.Request, playerRegNo int) bool {
	query := "SELECT * FROM tblindividualpunishment WHERE playerRegNo = ?"
	rows, err := db.Query(query, playerRegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the individual punishment
		var individualPunishment models.IndividualPunishment
		err = rows.Scan(&individualPunishment.TournamentId, &individualPunishment.MatchId, &individualPunishment.PlayerRegNo, &individualPunishment.TeamDeptCode, &individualPunishment.PunishmentType)
		if err != nil {
			panic(err.Error())
		}
		// now call delete individual punishment api
		url := host + "http://localhost:5000/api/match/individualpunishment/" + individualPunishment.TournamentId + "/" + individualPunishment.MatchId + "/" + strconv.Itoa(individualPunishment.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// player has individual score or not
func playerHasIndividualScore(w http.ResponseWriter, r *http.Request, playerRegNo int) bool {
	query := "SELECT * FROM tblindividualscore WHERE playerRegNo = ?"
	rows, err := db.Query(query, playerRegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the individual score
		var individualScore models.IndividualScore
		err = rows.Scan(&individualScore.TournamentId, &individualScore.MatchId, &individualScore.PlayerRegNo, &individualScore.TeamDeptCode, &individualScore.Goals)
		if err != nil {
			panic(err.Error())
		}
		// now call delete individual score api
		url := host + "http://localhost:5000/api/match/individualscore/" + individualScore.TournamentId + "/" + individualScore.MatchId + "/" + strconv.Itoa(individualScore.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// delete a dept
func deleteADept(deptCode int) {
	_, err := db.Query("DELETE FROM tbldept WHERE deptCode = ?", deptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a dept
func DeleteADept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/dept/{deptCode}", controller.DeleteADept).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["deptCode"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// dept exists or not
	if !deptExists(id) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Dept doesn't exist!")
		return
	}

	// prequisite check. check all tables where deptCode is used
	deptExistsInATeam(w, r, id)
	//deptExistsInAMatch(w, r, id)
	deptExistsInATiebreaker(w, r, id)
	deptExistsInAnIndividualScore(w, r, id)
	deptExistsInAnIndividualPunishment(w, r, id)
	deptExistsInAPlayer(w, r, id)

	deleteADept(id)

	json.NewEncoder(w).Encode("Dept deleted successfully!")
}

// dept exists in a tiebreaker or not
func deptExistsInATiebreaker(w http.ResponseWriter, r *http.Request, deptCode int) bool {
	query := "SELECT * FROM tbltiebreaker WHERE team1DeptCode = ? OR team2DeptCode = ?"
	rows, err := db.Query(query, deptCode, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the tiebreaker
		var tiebreaker models.Tiebreaker
		err = rows.Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)
		if err != nil {
			panic(err.Error())
		}
		// now call delete tiebreaker api
		url := host + "/api/match/tiebreaker/" + tiebreaker.TournamentId + "/" + tiebreaker.MatchId
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// dept exists in an individual score or not
func deptExistsInAnIndividualScore(w http.ResponseWriter, r *http.Request, deptCode int) bool {
	query := "SELECT * FROM tblindividualscore WHERE teamDeptCode = ?"
	rows, err := db.Query(query, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId, matchId and playerRegNo and delete the individual score
		var individualScore models.IndividualScore
		err = rows.Scan(&individualScore.TournamentId, &individualScore.MatchId, &individualScore.PlayerRegNo, &individualScore.TeamDeptCode, &individualScore.Goals)
		if err != nil {
			panic(err.Error())
		}
		// now call delete individual score api
		url := host + "/api/match/individualscore/" + individualScore.TournamentId + "/" + individualScore.MatchId + "/" + strconv.Itoa(individualScore.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// dept exists in an individual punishment or not
func deptExistsInAnIndividualPunishment(w http.ResponseWriter, r *http.Request, deptCode int) bool {
	query := "SELECT * FROM tblindividualpunishment WHERE teamDeptCode = ?"
	rows, err := db.Query(query, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId, matchId and playerRegNo and delete the individual punishment
		var individualPunishment models.IndividualPunishment
		err = rows.Scan(&individualPunishment.TournamentId, &individualPunishment.MatchId, &individualPunishment.PlayerRegNo, &individualPunishment.TeamDeptCode, &individualPunishment.PunishmentType)
		if err != nil {
			panic(err.Error())
		}
		// now call delete individual punishment api
		url := host + "/api/match/individualpunishment/" + individualPunishment.TournamentId + "/" + individualPunishment.MatchId + "/" + strconv.Itoa(individualPunishment.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// dept exists in a player table or not
func deptExistsInAPlayer(w http.ResponseWriter, r *http.Request, deptCode int) bool {
	query := "SELECT * FROM tblplayer WHERE playerDeptCode = ?"
	rows, err := db.Query(query, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the playerRegNo and delete the player
		var player models.Player
		err = rows.Scan(&player.PlayerRegNo, &player.PlayerSession, &player.PlayerSemester, &player.PlayerName, &player.PlayerDeptCode, &player.PlayerJerseyNo)
		if err != nil {
			panic(err.Error())
		}
		// now call delete player api
		url := host + "/api/player/" + strconv.Itoa(player.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// dept exists in a team or not
func deptExistsInATeam(w http.ResponseWriter, r *http.Request, deptCode int) bool {
	query := "SELECT * FROM tblteam WHERE deptCode = ?"
	rows, err := db.Query(query, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and delete the team
		var team models.Team
		err = rows.Scan(&team.TournamentId, &team.TeamSubmissionDate, &team.DeptCode, &team.TeamManagerEmail, &team.TeamCaptainRegID, &team.PlayerRegNo[0], &team.PlayerRegNo[1], &team.PlayerRegNo[2], &team.PlayerRegNo[3], &team.PlayerRegNo[4], &team.PlayerRegNo[5], &team.PlayerRegNo[6], &team.PlayerRegNo[7], &team.PlayerRegNo[8], &team.PlayerRegNo[9], &team.PlayerRegNo[10], &team.PlayerRegNo[11], &team.PlayerRegNo[12], &team.PlayerRegNo[13], &team.PlayerRegNo[14], &team.PlayerRegNo[15], &team.PlayerRegNo[16], &team.PlayerRegNo[17], &team.PlayerRegNo[18], &team.PlayerRegNo[19], &team.IsKnockedOut)
		if err != nil {
			panic(err.Error())
		}
		// now call delete team api
		url := host + "/api/tournament/team/" + team.TournamentId + "/" + strconv.Itoa(team.DeptCode)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// dept exists in a match or not
func deptExistsInAMatch(w http.ResponseWriter, r *http.Request, deptCode int) bool {
	query := "SELECT * FROM tblmatch WHERE team1DeptCode = ? OR team2DeptCode = ?"
	rows, err := db.Query(query, deptCode, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the match
		var match models.Match
		err = rows.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, match.MatchFourthRefereeID, &match.Venue)
		if err != nil {
			panic(err.Error())
		}
		// now call delete match api
		url := host + "/api/tournament/match/" + match.TournamentId + "/" + match.MatchId
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// delete a team
func deleteATeam(tournamentId string, deptCode int) {
	_, err := db.Query("DELETE FROM tblteam WHERE tournamentId = ? AND deptCode = ?", tournamentId, deptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a team
func DeleteATeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/tournament/team/{tournamentId}/{deptCode}", controller.DeleteATeam).Methods("DELETE")
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
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Team doesn't exist!")
		return
	}

	// prequisite check. check all tables where tournamentId and deptCode is used
	teamExistsInAMatch(w, r, tournamentId, deptCodeInt)

	deleteATeam(tournamentId, deptCodeInt)

	json.NewEncoder(w).Encode("Team deleted successfully!")
}

// team exists in a match or not
func teamExistsInAMatch(w http.ResponseWriter, r *http.Request, tournamentId string, deptCode int) bool {
	query := "SELECT * FROM tblmatch WHERE tournamentId = ? AND (team1DeptCode = ? OR team2DeptCode = ?)"
	rows, err := db.Query(query, tournamentId, deptCode, deptCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the match
		var match models.Match
		err = rows.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)
		if err != nil {
			panic(err.Error())
		}
		// now call delete match api
		url := host + "/api/tournament/match/" + match.TournamentId + "/" + match.MatchId
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// delete a match
func deleteAMatch(tournamentId string, matchId string) {
	_, err := db.Query("DELETE FROM tblmatch WHERE tournamentId = ? AND matchID = ?", tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a match
func DeleteAMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/tournament/match/{tournamentId}/{matchId}", controller.DeleteAMatch).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// prequisite check. check all tables where tournamentId and matchId is used
	matchExistsInATiebreaker(w, r, tournamentId, matchId)
	matchExistsInAnIndividualScore(w, r, tournamentId, matchId)
	matchExistsInAnIndividualPunishment(w, r, tournamentId, matchId)
	matchExistsInAStartingEleven(w, r, tournamentId, matchId)

	deleteAMatch(tournamentId, matchId)

	json.NewEncoder(w).Encode("Match deleted successfully!")
}

// match exists in a starting eleven or not
func matchExistsInAStartingEleven(w http.ResponseWriter, r *http.Request, tournamentId string, matchId string) bool {
	query := "SELECT * FROM tblplaying11 WHERE tournamentId = ? AND matchID = ?"
	rows, err := db.Query(query, tournamentId, matchId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		var startingEleven models.StartingEleven
		err = rows.Scan(&startingEleven.TournamentId, &startingEleven.MatchId, &startingEleven.TeamDeptCode, &startingEleven.StartingPlayerRegNo[0], &startingEleven.StartingPlayerRegNo[1], &startingEleven.StartingPlayerRegNo[2], &startingEleven.StartingPlayerRegNo[3], &startingEleven.StartingPlayerRegNo[4], &startingEleven.StartingPlayerRegNo[5], &startingEleven.StartingPlayerRegNo[6], &startingEleven.StartingPlayerRegNo[7], &startingEleven.StartingPlayerRegNo[8], &startingEleven.StartingPlayerRegNo[9], &startingEleven.StartingPlayerRegNo[10], &startingEleven.SubstitutePlayerRegNo[0], &startingEleven.SubstitutedPlayerRegNo[0], &startingEleven.SubstitutePlayerRegNo[1], &startingEleven.SubstitutedPlayerRegNo[1], &startingEleven.SubstitutePlayerRegNo[2], &startingEleven.SubstitutedPlayerRegNo[2])

		if err != nil {
			panic(err.Error())
		}
		// now call delete starting eleven api
		url := host + "/api/match/startingeleven/" + startingEleven.TournamentId + "/" + startingEleven.MatchId + "/" + strconv.Itoa(startingEleven.TeamDeptCode)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// match exists in a tiebreaker or not
func matchExistsInATiebreaker(w http.ResponseWriter, r *http.Request, tournamentId string, matchId string) bool {
	query := "SELECT * FROM tbltiebreaker WHERE tournamentId = ? AND matchID = ?"
	rows, err := db.Query(query, tournamentId, matchId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the tiebreaker
		var tiebreaker models.Tiebreaker
		err = rows.Scan(&tiebreaker.TournamentId, &tiebreaker.MatchId, &tiebreaker.Team1DeptCode, &tiebreaker.Team2DeptCode, &tiebreaker.Team1TieBreakerScore, &tiebreaker.Team2TieBreakerScore)
		if err != nil {
			panic(err.Error())
		}
		// now call delete tiebreaker api
		url := host + "/api/match/tiebreaker/" + tiebreaker.TournamentId + "/" + tiebreaker.MatchId
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// match exists in an individual score or not
func matchExistsInAnIndividualScore(w http.ResponseWriter, r *http.Request, tournamentId string, matchId string) bool {
	query := "SELECT * FROM tblindividualscore WHERE tournamentId = ? AND matchID = ?"
	rows, err := db.Query(query, tournamentId, matchId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId, matchId and playerRegNo and delete the individual score
		var individualScore models.IndividualScore
		err = rows.Scan(&individualScore.TournamentId, &individualScore.MatchId, &individualScore.PlayerRegNo, &individualScore.TeamDeptCode, &individualScore.Goals)
		if err != nil {
			panic(err.Error())
		}
		// now call delete individual score api
		url := host + "/api/match/individualscore/" + individualScore.TournamentId + "/" + individualScore.MatchId + "/" + strconv.Itoa(individualScore.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// match exists in an individual punishment or not
func matchExistsInAnIndividualPunishment(w http.ResponseWriter, r *http.Request, tournamentId string, matchId string) bool {
	query := "SELECT * FROM tblindividualpunishment WHERE tournamentId = ? AND matchID = ?"
	rows, err := db.Query(query, tournamentId, matchId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId, matchId and playerRegNo and delete the individual punishment
		var individualPunishment models.IndividualPunishment
		err = rows.Scan(&individualPunishment.TournamentId, &individualPunishment.MatchId, &individualPunishment.PlayerRegNo, &individualPunishment.TeamDeptCode, &individualPunishment.PunishmentType)
		if err != nil {
			panic(err.Error())
		}
		// now call delete individual punishment api
		url := host + "/api/match/individualpunishment/" + individualPunishment.TournamentId + "/" + individualPunishment.MatchId + "/" + strconv.Itoa(individualPunishment.PlayerRegNo)
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// delete a starting eleven
func deleteAStartingEleven(tournamentId string, matchId string, teamDeptCode int) {
	_, err := db.Query("DELETE FROM tblplaying11 WHERE tournamentId = ? AND matchID = ? AND teamDeptCode = ?", tournamentId, matchId, teamDeptCode)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a starting eleven
func DeleteAStartingEleven(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}", controller.DeleteAStartingEleven).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and teamDeptCode from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	teamDeptCode, _ := params["teamDeptCode"]

	// convert teamDeptCode from string to int
	teamDeptCodeInt, err := strconv.Atoi(teamDeptCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// starting eleven exists or not
	if !startingElevenExists(tournamentId, matchId, teamDeptCodeInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Starting eleven doesn't exist!")
		return
	}

	deleteAStartingEleven(tournamentId, matchId, teamDeptCodeInt)
	json.NewEncoder(w).Encode("Starting eleven deleted successfully!")
}

// delete a referee
func deleteAReferee(refereeId int) {
	_, err := db.Query("DELETE FROM tblreferee WHERE refereeID = ?", refereeId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a referee
func DeleteAReferee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/referee/{refereeId}", controller.DeleteAReferee).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// convert id from string to int
	id, err := strconv.Atoi(params["refereeId"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// referee exists or not
	if !refereeExists(id) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Referee doesn't exist!")
		return
	}

	// prequisite check. check all tables where refereeId is used
	refereeExistsInAMatch(w, r, id)

	deleteAReferee(id)

	json.NewEncoder(w).Encode("Referee deleted successfully!")
}

// referee exists in a match or not
func refereeExistsInAMatch(w http.ResponseWriter, r *http.Request, refereeId int) bool {
	query := "SELECT * FROM tblmatch WHERE matchRefereeID = ? OR matchLineman1ID = ? OR matchLineman2ID = ? OR matchFourthRefereeID = ?"
	rows, err := db.Query(query, refereeId, refereeId, refereeId, refereeId)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	for rows.Next() {
		// get the tournamentId and matchId and delete the match
		var match models.Match
		err = rows.Scan(&match.TournamentId, &match.MatchId, &match.MatchDate, &match.Team1DeptCode, &match.Team2DeptCode, &match.Team1Score, &match.Team2Score, &match.WinnerTeamDeptCode, &match.MatchRefereeID, &match.MatchLinesman1ID, &match.MatchLinesman2ID, &match.MatchFourthRefereeID, &match.Venue)
		if err != nil {
			panic(err.Error())
		}
		// now call delete match api
		url := host + "/api/tournament/match/" + match.TournamentId + "/" + match.MatchId
		// // get cookie and set it in request header
		// cookie, err := r.Cookie("jwtToken")
		// if err != nil {
		// 	panic(err.Error())
		// }
		// // set cookie in request header
		// w.Header().Set("Cookie", cookie.String())
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err.Error())
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		defer resp.Body.Close()
	}

	return true
}

// delete a tiebreaker
func deleteATiebreaker(tournamentId string, matchId string) {
	_, err := db.Query("DELETE FROM tbltiebreaker WHERE tournamentId = ? AND matchID = ?", tournamentId, matchId)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete a tiebreaker
func DeleteATiebreaker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/tiebreaker/{tournamentId}/{matchId}", controller.DeleteATiebreaker).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId and matchId from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// tiebreaker exists or not
	if !tiebreakerExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Tiebreaker doesn't exist!")
		return
	}

	deleteATiebreaker(tournamentId, matchId)

	json.NewEncoder(w).Encode("Tiebreaker deleted successfully!")
}

// delete an individual score
func deleteAnIndividualScore(tournamentId string, matchId string, playerRegNo int) {
	_, err := db.Query("DELETE FROM tblindividualscore WHERE tournamentId = ? AND matchID = ? AND playerRegNo = ?", tournamentId, matchId, playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete an individual score
func DeleteAnIndividualScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/individualscore/{tournamentId}/{matchId}/{playerRegNo}", controller.DeleteAnIndividualScore).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and playerRegNo from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	playerRegNo, _ := params["playerRegNo"]

	// convert playerRegNo from string to int
	playerRegNoInt, err := strconv.Atoi(playerRegNo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// player playing in the match or not
	if !playerIsPlayingInAMatchOfATournament(tournamentId, matchId, playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not playing in the match!")
		return
	}

	// individual score exists or not
	if !individualScoreExists(tournamentId, matchId, playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Individual score doesn't exist!")
		return
	}

	deleteAnIndividualScore(tournamentId, matchId, playerRegNoInt)

	json.NewEncoder(w).Encode("Individual score deleted successfully!")
}

// individual score exists or not
func individualScoreExists(tournamentId string, matchId string, playerRegNo int) bool {
	query := "SELECT tournamentId FROM tblindividualscore WHERE tournamentId = ? AND matchID = ? AND playerRegNo = ?"
	rows, err := db.Query(query, tournamentId, matchId, playerRegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	count := 0
	for rows.Next() {
		count++
		break
	}

	if count > 0 {
		return true
	}

	return false
}

// delete an individual punishment
func deleteAnIndividualPunishment(tournamentId string, matchId string, playerRegNo int) {
	_, err := db.Query("DELETE FROM tblindividualpunishment WHERE tournamentId = ? AND matchID = ? AND playerRegNo = ?", tournamentId, matchId, playerRegNo)

	if err != nil {
		panic(err.Error())
	}
}

// controller function to delete an individual punishment
func DeleteAnIndividualPunishment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // token validation check
	// if isTokenValid(w, r) == false {
	// 	return
	// }

	// router.HandleFunc("/api/match/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}", controller.DeleteAnIndividualPunishment).Methods("DELETE")
	// get id from url
	params := mux.Vars(r)

	// get tournamentId, matchId and playerRegNo from url
	tournamentId, _ := params["tournamentId"]
	matchId, _ := params["matchId"]
	playerRegNo, _ := params["playerRegNo"]

	// convert playerRegNo from string to int
	playerRegNoInt, err := strconv.Atoi(playerRegNo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// match exists or not
	if !matchExists(tournamentId, matchId) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Match doesn't exist!")
		return
	}

	// player playing in the match or not
	if !playerIsPlayingInAMatchOfATournament(tournamentId, matchId, playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Player is not playing in the match!")
		return
	}

	// individual punishment exists or not
	if !individualPunishmentExists(tournamentId, matchId, playerRegNoInt) {
		// set response header as forbidden
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode("Individual punishment doesn't exist!")
		return
	}

	deleteAnIndividualPunishment(tournamentId, matchId, playerRegNoInt)

	json.NewEncoder(w).Encode("Individual punishment deleted successfully!")
}

// individual punishment exists or not
func individualPunishmentExists(tournamentId string, matchId string, playerRegNo int) bool {
	query := "SELECT tournamentId FROM tblindividualpunishment WHERE tournamentId = ? AND matchID = ? AND playerRegNo = ?"
	rows, err := db.Query(query, tournamentId, matchId, playerRegNo)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err.Error())
	}

	count := 0
	for rows.Next() {
		count++
		break
	}

	if count > 0 {
		return true
	}

	return false
}
