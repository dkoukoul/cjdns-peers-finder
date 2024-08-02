package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
    logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logrus.Fatalf("Failed to open log file: %s", err)
    }

    multiWriter := io.MultiWriter(os.Stdout, logFile)
    logger.SetOutput(multiWriter)
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    logger.SetLevel(logrus.InfoLevel)
}

// Handler function for the API endpoint
func peersHandler(w http.ResponseWriter, r *http.Request) {
	var requestPeer Peer
	if r.Method == http.MethodPost {
        err := json.NewDecoder(r.Body).Decode(&requestPeer)
        if err != nil {
			logger.Error("Invalid request payload:", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        err = savePeerToFile(requestPeer)
        if err != nil {
			logger.Error("Error saving peer:", err)
            http.Error(w, "Error saving peer", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        return
    }
	
	peers, err := findGoodPeers(requestPeer)
	if err != nil {
		logger.Error("Error finding good peers:", err)
		http.Error(w, "Error finding peers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func main() {
	logger.Info("Starting server on port 8080")
	http.HandleFunc("/api/peers", peersHandler)
	logger.Fatal(http.ListenAndServe(":8080", nil))
}