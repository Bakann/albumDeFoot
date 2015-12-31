package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type player struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

const dataFile = "./players.json"

var playerMutex = new(sync.Mutex)

// Handle players
func handlePlayers(w http.ResponseWriter, r *http.Request) {
	// Since multiple requests could come in at once, ensure we have a lock
	// around all file operations
	playerMutex.Lock()
	defer playerMutex.Unlock()

	// Stat the file, so we can find its current permissions
	fi, err := os.Stat(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to stat the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	// Read the players from the file.
	playerData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "POST":
		// Decode the JSON data
		var players []player
		if err := json.Unmarshal(playerData, &players); err != nil {
			http.Error(w, fmt.Sprintf("Unable to Unmarshal players from data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		// Add a new player to the in memory slice of players
		players = append(players, player{ID: time.Now().UnixNano() / 1000000, Author: r.FormValue("author"), Text: r.FormValue("text")})

		// Marshal the players to indented json.
		playerData, err = json.MarshalIndent(players, "", "    ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to marshal players to json: %s", err), http.StatusInternalServerError)
			return
		}

		// Write out the players to the file, preserving permissions
		err := ioutil.WriteFile(dataFile, playerData, fi.Mode())
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to write players to data file (%s): %s", dataFile, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		io.Copy(w, bytes.NewReader(playerData))

	case "GET":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache")
		// stream the contents of the file to the response
		io.Copy(w, bytes.NewReader(playerData))

	default:
		// Don't know the method, so error
		http.Error(w, fmt.Sprintf("Unsupported method: %s", r.Method), http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.HandleFunc("/api/players", handlePlayers)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Server started: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
