package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
    logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logrus.Fatalf("Failed to open log file: %s", err)
    }

    multiWriter := io.MultiWriter(os.Stdout, logFile)
    logger.SetOutput(multiWriter)
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })
    logger.SetLevel(logrus.InfoLevel)
}

// Handler function for the API endpoint
func peersHandler(w http.ResponseWriter, r *http.Request) {
	var requestPeer Peer
	if r.Method == http.MethodPost {
        remoteAddr := r.RemoteAddr
        logger.Info("Request received from:", remoteAddr)

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
    }
	
	peers, err := findGoodPeers(requestPeer)
	if err != nil {
		logger.Error("Error finding good peers:", err)
		http.Error(w, "Error finding peers", http.StatusInternalServerError)
		return
	}
    if len(peers) == 0 {
        logger.Info("No good peers found")
        http.Error(w, "No good peers found", http.StatusNotFound)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func main() {

	if cjdnsToolsPath == "/tools/" {
        // DEFAULT for pkt-server check if peerstats exists in /server/cjdns/tools location
        if _, err := os.Stat("/server/cjdns/tools/peerStats"); os.IsNotExist(err) {
            logger.Fatal("peerStats not found in /server/cjdns/tools, try setting CJDNS_PATH env var. Exiting...")
            os.Exit(1)
        } else {
            logger.Info("cjdns tools found in /server/cjdns/tools")
        }
    }

    logger.Info("Starting server on port ", serverPort)
	http.HandleFunc("/api/peers", peersHandler)
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    logger.Info("cjdns tools path fount at ", cjdnsToolsPath)
    // Start a goroutine that calls peerTest every hour
    go func() {
        for {
            select {
            case <-ticker.C:
                logger.Info("Running peer test at ", time.Now())
                peerTest()
            }
        }
    }()

	logger.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverPort), nil))
}