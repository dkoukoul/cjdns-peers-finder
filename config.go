package main

import (
	"github.com/sirupsen/logrus"
)

const routeServerUrl = "https://routeserver.cjd.li/api/"
const nodesInfoEndpoint = "ni"
const peersFilePath = "peers.json"

var logger = logrus.New()

//RESPONSE Peer represents a CJDNS peer 
type Peer struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	PublicKey string `json:"publicKey"`
}

//RouteServer /ni RESPONSE 
type Node struct {
    Announcements int    `json:"announcements"`
    IP6           string `json:"ip6"`
    Rst           bool   `json:"rst"`
}

type PeerInfo struct {
    AnnByHashLen  int   `json:"annByHashLen"`
    Announcements int   `json:"announcements"`
    Peers         []int `json:"peers"`
}

type NodeDataResponseData struct {
    Nodes             []Node   `json:"nodes"`
    PeerInfo          PeerInfo `json:"peerInfo"`
    TotalAnnouncements int     `json:"totalAnnouncements"`
    TotalNodes        int      `json:"totalNodes"`
    TotalWithRsts     int      `json:"totalWithRsts"`
}

//RouteServer /ni/cjdns-ipv6 RESPONSE 
type NodeInfoResponseData struct {
	NodeInfo NodeInfo `json:"node"`
}

type NodeInfo struct {
    EncodingScheme   []interface{}       `json:"encodingScheme"`
    InwardLinksByIp  map[string][]Link   `json:"inwardLinksByIp"`
    Ipv6             string              `json:"ipv6"`
    Key              string              `json:"key"`
    Type             string              `json:"type"`
	Mut              interface{}         `json:"mut"`
    Version          int                 `json:"version"`
}

type Link struct {
    // Define fields if necessary
}