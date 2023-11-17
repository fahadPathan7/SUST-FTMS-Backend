# My RESTful API

This comprehensive RESTful API facilitates the management of tournaments, players, teams, referees, and other essential components.

**Note:** For PUT and POST operations, ensure that all required information is provided.

## Table of Contents
- [JSON Web Token](#jwt)
- [Operator](#operator)
- [Department](#department)
- [Player](#player)
- [Teacher](#teacher)
- [Team Manager](#team-manager)
- [Team](#team)
- [Tournament](#tournament)
- [Referee](#referee)
- [Match](#match)
- [Starting Eleven](#starting-eleven)
- [Individual Punishment](#individual-punishment)
- [Individual Score](#individual-score)
- [Tiebreaker](#tiebreaker)
- [Setting up the backend](#setting-up-the-backend)


## Endpoints

### <a name="jwt"></a> JSON Web Token

- `GET /api/token/generate/{userEmail}` - Generate token
- `GET /api/token/validate` - Validate token
<br><br>
### <a name="operator"></a> Operator

- `POST /api/operator/login` - Operator login
- `GET /api/operator/{email}` - Get operator info (password will be null)

**A JSON sample for operator**
```json
{
  "email": "YourEmail",
  "password": "YourPassword",
  "name": "YourName",
  "office": "YourOffice"
}
```
<br><br>
### <a name="department"></a> Department

- `POST /api/dept` - Insert a new department
- `PUT /api/dept/{deptCode}` - Update a department
- `GET /api/depts` - Get all departments
- `GET /api/dept/{deptCode}` - Get a specific department
- `DELETE /api/dept/{deptCode}` - Delete a department

**A JSON sample for department**
```json
{
  "deptCode": 1,
  "deptName": "Department Name",
  "deptHeadName": "Department Head Name",
  "deptShortName": "DeptShortName"
}
```
<br><br>
### <a name="player"></a> Player

- `POST /api/player` - Insert a new player
- `PUT /api/player/{playerRegNo}` - Update a player
- `GET /api/player/{playerRegNo}` - Get a specific player
- `GET /api/dept/players/{deptCode}` - Get all players of a department
- `DELETE /api/player/{playerRegNo}` - Delete a player

**A JSON sample for player**
```json
{
  "playerRegNo": 1,
  "playerName": "Player Name",
  "playerDeptCode": 1,
  "playerEmail": "abc@student.sust.edu",
  "playerPassword": "***",
  "playerImage": "local loc"
}
```
<br><br>
### <a name="teacher"></a> Teacher

- `POST /api/techer` - Insert a new teacher
- `GET /api/teacher/{email}` - Get a teacher
- `GET /api/teachers/{deptCode}` - Get all the teachers of a dept

**A JSON sample for teacher**
```json
{
  "email": "teacher email",
  "name": "teacher name",
  "deptCode": 1,
  "title": "lecturer"
}
```
<br><br>
### <a name="team-manager"></a> Team Manager

- `POST /api/teammanager` - Insert a new team manager in a tournament
- `GET /api/teammanager/{tournamentId}/{email}` - Get a team manager of a tournament
- `GET /api/teammanagers/{tournamentId}` - Get all team managers of a tournament
- `DELETE /api/tournament/teammanager/{tournamentId}/{teamManagerEmail}` - Delete a team manager of a tournament

**A JSON sample for team manager**
```json
{
  "email": "Manager email",
  "TournamentId": "TournamentID"
}
```
<br><br>
### <a name="team"></a> Team

- `POST /api/team` - Insert a new team
- `PUT /api/tournament/team/{tournamentId}/{deptCode}` - Update a team
- `GET /api/tournament/team/{tournamentId}/{deptCode}` - Get a team
- `GET /api/tournament/teams/{tournamentId}` - Get all teams of a tournament
- `DELETE /api/tournament/team/{tournamentId}/{deptCode}` - Delete a team

**A JSON sample for team**
```json
{
  "tournamentId": "Tournament ID",
  "teamSubmissionDate": "Submission Date",
  "deptCode": 1,
  "teamManagerEmail": "Team Manager",
  "teamCaptainRegID": 1,
  "playerRegNo": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
  "isKnockedOut": false
}
```
<br><br>
### <a name="tournament"></a> Tournament

- `POST /api/tournament` - Insert a new tournament
- `PUT /api/tournament/{tournamentId}` - Update a tournament
- `GET /api/tournaments` - Get all tournaments
- `GET /api/tournament/{tournamentId}` - Get a specific tournament
- `DELETE /api/tournament/{tournamentId}` - Delete a tournament

**A JSON sample for tournament**
```json
{
  "tournamentId": "Tournament ID",
  "tournamentName": "Tournament Name",
  "startingDate": "Starting Date",
  "endingDate": "Ending Date"
}
```
<br><br>
### <a name="referee"></a> Referee

- `POST /api/referee` - Insert a new referee
- `PUT /api/referee/{refereeId}` - Update a referee
- `GET /api/referees` - Get all referees
- `GET /api/referee/{refereeId}` - Get a specific referee
- `DELETE /api/referee/{refereeId}` - Delete a specific referee

**A JSON sample for referee**
```json
{
  "refereeID": 1,
  "refereeName": "Referee Name",
  "refereeInstitute": "Referee Institute"
}
```
<br><br>
### <a name="match"></a> Match

- `POST /api/match` - Insert a new match
- `PUT /api/match/{tournamentId}/{matchId}` - Update a match
- `GET /api/tournament/matches/{tournamentId}` - Get all matches of a tournament
- `GET /api/tournament/match/{tournamentId}/{matchId}` - Get a specific match in a tournament
- `DELETE /api/match/{tournamentId}/{matchId}` - Delete a match

**Note:** For POST, everything can be nill except tournamentId and matchId.

**A JSON sample for match**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "matchDate": "Match Date",
  "team1DeptCode": 1,
  "team2DeptCode": 2,
  "team1Score": 3,
  "team2Score": 4,
  "winnerTeamDeptCode": 1,
  "matchRefereeID": 5,
  "matchLinesman1ID": 6,
  "matchLinesman2ID": 7,
  "matchFourthRefereeID": 8,
  "venue": "SUST Central Field"
}
```
<br><br>
### <a name="starting-eleven"></a> Starting Eleven

- `POST /api/match/startingeleven` - Insert a new starting eleven
- `PUT /api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}` - Update a starting eleven
- `GET /api/match/startingeleven/{tournamentId}/{matchId}/{deptCode}` - Get the starting eleven of a team in a match
- `DELETE /api/match/startingeleven/{tournamentId}/{matchId}/{teamDeptCode}` - Delete a starting eleven

**A JSON sample for starting eleven**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "teamDeptCode": 1,
  "startingPlayerRegNo": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11],
  "substitutePlayerRegNo": [12, 13, 14],
  "substituedPlayerRegNo": [1, 2, 3]
}
```
<br><br>
### <a name="individual-punishment"></a> Individual Punishment

- `POST /api/individualpunishment` - Insert a new individual punishment
- `PUT /api/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}` - Update an individual punishment
- `GET /api/tournament/individualpunishments/{tournamentId}` - Get all individual punishments (all players) of a tournament
- `GET /api/tournament/match/team/individualpunishments/{tournamentId}/{matchId}/{teamDeptCode}` - Get all individual punishments of a match by a team
- `GET /api/tournament/player/individualpunishments/{tournamentId}/{playerRegNo}` - Get all individual punishments of a player in a tournament
- `DELETE /api/individualpunishment/{tournamentId}/{matchId}/{playerRegNo}` - Delete an individual punishment

**A JSON sample for individual punishment**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "playerRegNo": 1,
  "teamDeptCode": 1,
  "punishmentType": "Punishment Type"
}
```
<br><br>
### <a name="individual-score"></a> Individual Score

- `POST /api/individualscore` - Insert a new individual score
- `PUT /api/individualscore/{tournamentId}/{matchId}/{playerRegNo}` - Update an individual score
- `GET /api/tournament/individualscores/{tournamentId}` - Get all individual scores (all players) of a tournament
- `GET /api/tournament/player/individualscores/{tournamentId}/{playerRegNo}` - Get all individual scores of a player in a tournament
- `GET /api/tournament/match/team/individualscores/{tournamentId}/{matchId}/{teamDeptCode}` - Get all individual scores of a match by a team
- `DELETE /api/individualscore/{tournamentId}/{matchId}/{playerRegNo}` - Delete an individual score

**A JSON sample for individual score**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "playerRegNo": 1,
  "teamDeptCode": 1,
  "goals": 2
}
```
<br><br>
### <a name="tiebreaker"></a> Tiebreaker

- `POST /api/tiebreaker` - Insert a new tiebreaker
- `PUT /api/tiebreaker/{tournamentId}/{matchId}` - Update a tiebreaker
- `GET /api/tournament/tiebreakers/{tournamentId}` - Get all tiebreakers of a tournament
- `GET /api/tournament/tiebreaker/{tournamentId}/{matchId}` - Get a tiebreaker of a match of a tournament
- `DELETE /api/tiebreaker/{tournamentId}/{matchId}` - Delete a tiebreaker

**A JSON sample for tiebreaker**
```json
{
  "tournamentId": "Tournament ID",
  "matchId": "Match ID",
  "team1DeptCode": 1,
  "team2DeptCode": 2,
  "team1TieBreakerScore": 3,
  "team2TieBreakerScore": 4
}
```
<br><br>
## <a name="setting-up-the-backend"></a> Setting Up the Backend

**1. Database Setup**

- Create a database named `ftms` in MySQL (XAAMP).

**2. Table Creation**

- Copy the codes inside `allTables.txt` file from the `Database` folder.
- Paste the contents of `allTables.txt` into your MySQL client (XAAMP) and execute the script to create the necessary tables.

**3. Running the Backend**

- Open the project in VS Code.
- Run the `main.go` file to start the backend server.

**Congratulations! You're now ready to work with the backend.**

---

Thank you for reviewing the documentation for My RESTful API. This API is designed to provide a comprehensive and user-friendly interface for managing tournaments, players, teams, referees, and other essential components for the SUST Football Tournament Management System.
