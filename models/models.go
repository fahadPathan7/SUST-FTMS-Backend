package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Operator struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Office   string `json:"office"`
}

type TeamManager struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	DeptCode int    `json:"deptCode"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type Dept struct {
	DeptCode      int    `json:"deptCode"`
	DeptName      string `json:"deptName"`
	DeptHeadName  string `json:"deptHeadName"`
	DeptShortName string `json:"deptShortName"`
}

type Player struct {
	PlayerRegNo    int    `json:"playerRegNo"`
	PlayerSession  string `json:"playerSession"`
	PlayerSemester int    `json:"playerSemester"`
	PlayerName     string `json:"playerName"`
	PlayerDeptCode int    `json:"playerDeptCode"`
	PlayerJerseyNo int    `json:"playerJerseyNo"`
}

type Team struct {
	TournamentId       string  `json:"tournamentId"`
	TeamSubmissionDate string  `json:"teamSubmissionDate"`
	DeptCode           int     `json:"deptCode"`
	TeamManagerEmail   string  `json:"teamManagerEmail"`
	TeamCaptainRegID   int     `json:"teamCaptainRegID"`
	PlayerRegNo        [20]int `json:"playerRegNo"`
	IsKnockedOut       bool    `json:"isKnockedOut"`
}

type IndividualPunishment struct {
	TournamentId   string `json:"tournamentId"`
	MatchId        string `json:"matchId"`
	PlayerRegNo    int    `json:"playerRegNo"`
	TeamDeptCode   int    `json:"teamDeptCode"`
	PunishmentType string `json:"punishmentType"`
}

type IndividualScore struct {
	TournamentId string `json:"tournamentId"`
	MatchId      string `json:"matchId"`
	PlayerRegNo  int    `json:"playerRegNo"`
	TeamDeptCode int    `json:"teamDeptCode"`
	Goals        int    `json:"goals"`
}

type Tournament struct {
	TournamentId   string `json:"tournamentId"`
	TournamentName string `json:"tournamentName"`
	StartingDate   string `json:"startingDate"`
	EndingDate     string `json:"endingDate"`
}

type Tiebreaker struct {
	TournamentId         string `json:"tournamentId"`
	MatchId              string `json:"matchId"`
	Team1DeptCode        int    `json:"team1DeptCode"`
	Team2DeptCode        int    `json:"team2DeptCode"`
	Team1TieBreakerScore int    `json:"team1TieBreakerScore"`
	Team2TieBreakerScore int    `json:"team2TieBreakerScore"`
}

type Match struct {
	TournamentId         string `json:"tournamentId"`
	MatchId              string `json:"matchId"`
	MatchDate            string `json:"matchDate"`
	Team1DeptCode        int    `json:"team1DeptCode"`
	Team2DeptCode        int    `json:"team2DeptCode"`
	Team1Score           int    `json:"team1Score"`
	Team2Score           int    `json:"team2Score"`
	WinnerTeamDeptCode   int    `json:"winnerTeamDeptCode"`
	MatchRefereeID       int    `json:"matchRefereeID"`
	MatchLinesman1ID     int    `json:"matchLinesman1ID"`
	MatchLinesman2ID     int    `json:"matchLinesman2ID"`
	MatchFourthRefereeID int    `json:"matchFourthRefereeID"`
	Venue                string `json:"venue"`
}

type StartingEleven struct {
	TournamentId           string  `json:"tournamentId"`
	MatchId                string  `json:"matchId"`
	TeamDeptCode           int     `json:"teamDeptCode"`
	StartingPlayerRegNo    [11]int `json:"startingPlayerRegNo"`
	SubstitutePlayerRegNo  [3]int  `json:"substitutePlayerRegNo"`
	SubstitutedPlayerRegNo [3]int  `json:"substituedPlayerRegNo"`
}

type Referee struct {
	RefereeID        int    `json:"refereeID"`
	RefereeName      string `json:"refereeName"`
	RefereeInstitute string `json:"refereeInstitute"`
}
