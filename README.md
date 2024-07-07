# Live Cricket Score Server

Go-based HTTP server that fetches and display live cricket score.It retrieves data from a Cricket Score JSON API and Display the Live Cricket Score Data in `text/plain` format in Terminal and Browsers.  

## Features

- Fetch live cricket scores using a configurable API.
- Display match details including title, update, live score, match date, and run rate.
- List current batsmen and their performance.
- List current bowlers and their performance.
- Error Handling and Validations

## Setup

- Clone or Download this Repo

```sh
git clone https://github.com/mskian/live-cricket-score-server.git
cd live-cricket-score-server
```

- Create a `config.yaml` file in the project root with the Cricket Score API URL

```yaml

api_url: "https://your-cricket-score-api-url.com/score?id="

```

- Run the Server

```sh
go run main.go
```

- Check Live Cricket Score on Terminal or Web Browser

```sh
curl http://localhost:6053/livescore?id=YOUR_MATCH_ID
```

## API Endpoints

- GET /livescore?id={matchID}:

Fetches and displays live cricket scores for the specified match ID.

- GET /404:

Custom 404 Page Not Found handler.

- GET /500:

Custom 500 Internal Server Error handler.

- Example Response in Terminal

```yaml
Match Details:

  Title: Example Match
  Update: Example Update
  Live Score: 123/4
  Match Date: 2024-07-07
  Run Rate: 6.5

Current Batsmen:

  - Name: Batsman 1
    Runs: 56
    Balls: 34
    Strike Rate: 164.70

  - Name: Batsman 2
    Runs: 23
    Balls: 15
    Strike Rate: 153.33

Current Bowlers:

  - Name: Bowler 1
    Overs: 4
    Runs: 25
    Wickets: 2

  - Name: Bowler 2
    Overs: 3
    Runs: 18
    Wickets: 1
```

## Build Package

- Run Make file to build a package for your Systems

```sh
make build
```

## Packges Build for  

Linux, Apple, Windows and Android - `/makefile`  

- Linux-386
- Linux-arm-7
- Linux-amd64
- Linux-arm64
- Andriod-arm64
- windows-386
- windows-amd64
- darwin-amd64
- darwin-arm64

```sh
chmod +x score
./score
```

## API

Live Cricket Score JSON API - **<https://github.com/sanwebinfo/cricket-score/tree/main/data>**

## Disclaimer 🗃

- This is not an Offical API from Cricbuzz - it's an Unofficial API
- This is for Education Purpose only - use at your own risk on Production Site

All Credits Goes to <https://www.cricbuzz.com/>  

## LICENSE

MIT
