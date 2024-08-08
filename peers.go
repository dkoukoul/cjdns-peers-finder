package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/rand"
)

// Function to fetch and parse data
func fetchNodeData(url string) (*NodeDataResponseData, error) {
    resp, err := http.Get(url)
    if err != nil {
		logger.Error("Error fetching data:", err)
        return nil, err
    }
    defer resp.Body.Close()

    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		logger.Error("Error reading response body:", err)
        return nil, err
    }

    // Parse JSON response
    var data NodeDataResponseData
    err = json.Unmarshal(body, &data)
    if err != nil {
		logger.Error("Error parsing JSON response:", err)
        return nil, err
    }

    return &data, nil
}

func fetchNodeInfo(url string) (*NodeInfo, error) {
	resp, err := http.Get(url)
	logger.Info("Fetching data from:", url)
	if err != nil {
		logger.Error("Error fetching node info data:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return nil, err
	}

	// Parse JSON response
	var data NodeInfoResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error("Error parsing JSON response:", err)
		return nil, err
	}

	return &data.NodeInfo, nil
}

// readPeersFromFile reads the peers from the local JSON file and returns them as a slice of Peer.
func readPeersFromFile() ([]Peer, error) {
    var peers []Peer
    file, err := os.Open(peersFilePath)
    if err != nil {
        if os.IsNotExist(err) {
			logger.Info("Peers file does not exist")
            return peers, nil // Return an empty slice if the file does not exist
        }
        return nil, err
    }
    defer file.Close()
    err = json.NewDecoder(file).Decode(&peers)
    if err != nil {
		logger.Error("Error decoding peers from file:", err)
        return nil, err
    }
    return peers, nil
}


// savePeerToFile saves a new peer to the local JSON file if it doesn't already exist.
func savePeerToFile(newPeer Peer) error {
    peers, err := readPeersFromFile()
    if err != nil {
        logger.Error("Error reading peers from file:", err)
        return err
    }

    // Check for duplicates and update if exists
    peerExists := false
    for i, peer := range peers {
        if peer.IP == newPeer.IP {
            peers[i] = newPeer
            peerExists = true
            break
        }
    }

    // If the peer does not exist, add the new peer
    if !peerExists {
        peers = append(peers, newPeer)
    }

    // Write the updated peers list back to the file
    file, err := os.Create(peersFilePath)
    if err != nil {
        logger.Error("Error creating peers file:", err)
        return err
    }
    defer file.Close()
    err = json.NewEncoder(file).Encode(peers)
    if err != nil {
        logger.Error("Error encoding peers to file:", err)
        return err
    }
    return nil
}

func shuffleNodes(nodes *[]Node) {
    rand.Seed(uint64(time.Now().UnixNano()))
    rand.Shuffle(len(*nodes), func(i, j int) {
        (*nodes)[i], (*nodes)[j] = (*nodes)[j], (*nodes)[i]
    })
}

// Function to find good peers for a given node
func findGoodPeers(forNode Peer) ([]Peer, error) {
	url := routeServerUrl + nodesInfoEndpoint
	logger.Info("Fetching data from:", url)
	nodeData, err := fetchNodeData(url)
    if err != nil {
		logger.Error("Error fetching data:", err)
        return nil, err
    }
	var peers []Peer
	savedPeers, err := readPeersFromFile()
	if err != nil {
		logger.Error("Error reading peers from file:", err)
		return nil, err
	}

	// Check if there are enough nodes to find good peers
	if len(savedPeers) <= MAX_RETURNING_PEERS {
		logger.Info("Not enough saved peers to find good peers, return all saved ones.")
		return savedPeers, nil
	}
	// Shuffle nodes to randomize the order of nodes and the resulting peers
	shuffleNodes(&nodeData.Nodes)
	// Get nodes and find 3 that do not have >100 peers
	for _,node := range nodeData.Nodes {
		for _,savedPeer := range savedPeers {
			// Check if the node is in our saved peers list, which means we have the information needed
			if (node.IP6 == savedPeer.IP6) && (savedPeer.Status == "ESTABLISHED") {
				nodeInfoUrl := routeServerUrl + nodesInfoEndpoint + "/" + node.IP6
				nodeInfo, err := fetchNodeInfo(nodeInfoUrl)
				if err != nil {
					logger.Error("Error fetching data:", err)
					return nil, err
				}
				// exclude requester and nodes with >=100 peers
				if (savedPeer.IP6 != forNode.IP6) && (len(nodeInfo.InwardLinksByIp) < 100) {
					logger.Info("Adding peer:", savedPeer.IP6)
					peers = append(peers, savedPeer)
				}
				if len(peers) == MAX_RETURNING_PEERS {
					logger.Info("Found",MAX_RETURNING_PEERS," good peers")
					return peers, nil
				}
			} 
		} 
	}
	logger.Info("Not enough good peers found")
	return peers, nil
}