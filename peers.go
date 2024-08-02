package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// Function to fetch and parse data
func fetchNodeData(url string) (*NodeDataResponseData, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Parse JSON response
    var data NodeDataResponseData
    err = json.Unmarshal(body, &data)
    if err != nil {
        return nil, err
    }

    return &data, nil
}

func fetchNodeInfo(url string) (*NodeInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var data NodeInfoResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
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
            return peers, nil // Return an empty slice if the file does not exist
        }
        return nil, err
    }
    defer file.Close()
    err = json.NewDecoder(file).Decode(&peers)
    if err != nil {
        return nil, err
    }
    return peers, nil
}


// savePeerToFile saves a new peer to the local JSON file if it doesn't already exist.
func savePeerToFile(newPeer Peer) error {
    peers, err := readPeersFromFile()
    if err != nil {
        return err
    }

    // Check for duplicates
    for _, peer := range peers {
        if peer.IP == newPeer.IP {
            return nil // Peer already exists, do nothing
        }
    }

    // Add the new peer
    peers = append(peers, newPeer)

    // Write the updated peers list back to the file
    file, err := os.Create(peersFilePath)
    if err != nil {
        return err
    }
    defer file.Close()
    err = json.NewEncoder(file).Encode(peers)
    if err != nil {
        return err
    }
    return nil
}

// Function to find good peers for a given node
func findGoodPeers(forNode Peer) ([]Peer, error) {
	url := routeServerUrl + nodesInfoEndpoint
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
	// Get nodes and find 3 that do not have >100 peers
	for _,node := range nodeData.Nodes {
		for _,savedPeer := range savedPeers {
			// Check if the node is in our saved peers list, which means we have the information needed
			if (node.IP6 == savedPeer.IP) && (node.IP6 != forNode.IP) {
				nodeInfoUrl := routeServerUrl + nodesInfoEndpoint + "/" + node.IP6
				nodeInfo, err := fetchNodeInfo(nodeInfoUrl)
				if err != nil {
					logger.Error("Error fetching data:", err)
					return nil, err
				}
				
				// exclude requester and nodes with >=100 peers
				if (nodeInfo.Ipv6 != forNode.IP) && (len(nodeInfo.InwardLinksByIp) < 100) {
					peers = append(peers, Peer{
						Name:      nodeInfo.Ipv6,
						Login:    "",
						Password:  "",
						IP:        nodeInfo.Ipv6,
						Port:      0,
						PublicKey: nodeInfo.Key,
					})
				}
				if len(peers) == 3 {
					break
				}
			}
		}
	}
	return peers, nil
}