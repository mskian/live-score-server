package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// Define the data structure for Batsman and Bowler.
type Batsman struct {
	Name       string `json:"name"`
	Runs       string `json:"runs"`
	Balls      string `json:"balls"`
	StrikeRate string `json:"strike_rate"`
}

type Bowler struct {
	Name    string `json:"name"`
	Overs   string `json:"overs"`
	Runs    string `json:"runs"`
	Wickets string `json:"wickets"`
}

type LiveScore struct {
	Title          string    `json:"title"`
	Update         string    `json:"update"`
	LiveScore      string    `json:"livescore"`
	MatchDate      string    `json:"match_date"`
	RunRate        string    `json:"runrate"`
	CurrentBatsmen []Batsman `json:"current_batsmen"`
	CurrentBowler  []Bowler  `json:"current_bowler"`
}

type Config struct {
	APIURL string `yaml:"api_url"`
}

const (
	timeout         = 10 * time.Second
	port            = 6053
	configFilename  = "config.yaml"
	escapeTextUsage = "Implement escaping if needed"
	maxMatchIDLen   = 10
)

var httpClient = &http.Client{
	Timeout: timeout,
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.APIURL == "" {
		return nil, fmt.Errorf("API URL is missing in config")
	}

	return &config, nil
}

func fetchScore(matchID string, apiURL string) (*LiveScore, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID cannot be empty")
	}

	url := fmt.Sprintf("%s%s", apiURL, matchID)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch score: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var score LiveScore
	if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if err := validateScore(score); err != nil {
		return nil, fmt.Errorf("invalid score data: %w", err)
	}

	return &score, nil
}

func validateScore(score LiveScore) error {
	if score.Title == "" || score.LiveScore == "" || score.MatchDate == "" || score.RunRate == "" {
		return fmt.Errorf("required fields are missing")
	}

	if len(score.CurrentBatsmen) == 0 || len(score.CurrentBowler) == 0 {
		return fmt.Errorf("batsmen or bowler data is missing")
	}

	for _, batsman := range score.CurrentBatsmen {
		if _, err := strconv.ParseFloat(batsman.StrikeRate, 64); err != nil {
			return fmt.Errorf("invalid strike rate for batsman %s: %w", batsman.Name, err)
		}
	}

	for _, bowler := range score.CurrentBowler {
		if _, err := strconv.ParseFloat(bowler.Overs, 64); err != nil {
			return fmt.Errorf("invalid overs for bowler %s: %w", bowler.Name, err)
		}
	}

	return nil
}

func formatScore(score *LiveScore) string {
	result := fmt.Sprintf(
		"\n\nMatch Details:\n\n  Title: %s\n  Update: %s\n  Live Score: %s\n  Match Date: %s\n  Run Rate: %s\n\n",
		escapeText(score.Title),
		escapeText(score.Update),
		escapeText(score.LiveScore),
		escapeText(score.MatchDate),
		escapeText(score.RunRate),
	)

	result += "Current Batsmen:\n\n"
	for _, batsman := range score.CurrentBatsmen {
		result += fmt.Sprintf(
			"  - Name: %s\n    Runs: %s\n    Balls: %s\n    Strike Rate: %s\n\n",
			escapeText(batsman.Name),
			escapeText(batsman.Runs),
			escapeText(batsman.Balls),
			escapeText(batsman.StrikeRate),
		)
	}

	result += "Current Bowlers:\n\n"
	for _, bowler := range score.CurrentBowler {
		result += fmt.Sprintf(
			"  - Name: %s\n    Overs: %s\n    Runs: %s\n    Wickets: %s\n\n",
			escapeText(bowler.Name),
			escapeText(bowler.Overs),
			escapeText(bowler.Runs),
			escapeText(bowler.Wickets),
		)
	}

	return result
}

func escapeText(text string) string {
	return text // Implement escaping if needed
}

func liveScoreHandler(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig(configFilename)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")

	matchID := r.URL.Query().Get("id")

	if matchID == "" {
		http.Error(w, "match ID is required", http.StatusBadRequest)
		return
	}
	if len(matchID) > maxMatchIDLen {
		http.Error(w, "match ID is too long", http.StatusBadRequest)
		return
	}

	score, err := fetchScore(matchID, config.APIURL)
	if err != nil {
		log.Printf("Error fetching score: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, formatScore(score))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 Page Not Found")
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "500 Internal Server Error")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/livescore", liveScoreHandler)
	mux.HandleFunc("/404", notFoundHandler)
	mux.HandleFunc("/500", internalServerErrorHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFoundHandler().ServeHTTP(w, r)
	})

	log.Printf("Server starting on port %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
