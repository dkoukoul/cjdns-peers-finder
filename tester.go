package main

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func checkStatus(peer Peer) string {
	cmd := exec.Command(cjdnsToolsPath + "peerStats")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("Error executing command: %v, output: %s", err, output)
		return ""
	}

	lines := strings.Split(string(output), "\n")
	// Iterate through each line
	for _, line := range lines {
		// Check if the line contains the peer's public key
		if strings.Contains(line, peer.PublicKey) {
			status := strings.Split(line, " ")[2]
			return status
		}
	}

	return ""
}
func connectPeer(peer Peer) error {
	// logger.Info("Connecting to ", peer.Name)
	ipAndPort := peer.IP + ":" + strconv.Itoa(peer.Port)
	cmd := exec.Command(cjdnsToolsPath+"cexec", "UDPInterface_beginConnection(\""+peer.PublicKey+"\",\""+ipAndPort+"\",0,\"\",\""+peer.Password+"\",\""+peer.Login+"\",0)")
	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorf("Error executing command: %v, output: %s", err, output)
		return err
	}
	return nil
}

func peerTest() {
	peers, err := readPeersFromFile()
	if err != nil {
		logger.Error("Error reading peers from file:", err)
		return
	}
	for _, peer := range peers {
		connectPeer(peer)
		status := ""
		for i := 0; i < 10; i++ {
			status = checkStatus(peer)
			if status == "ESTABLISHED" {
				break
			}
			time.Sleep(3 * time.Second)
		}
		peer.Status = status
		logger.Info("Peer ", peer.Name, " status: ", peer.Status)
		savePeerToFile(peer)
	}
}
