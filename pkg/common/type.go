package common

import (
	"time"
)

// IcmpReturn ICMP Response time details
type IcmpReturn struct {
	Success bool
	Addr    string
	Elapsed time.Duration
}

// IcmpHop ICMP HOP Response time details
type IcmpHop struct {
	Success  bool          `json:"success"`
	Address  string        `json:"address"`
	Host     string        `json:"host"`
	N        int           `json:"n"`
	TTL      int           `json:"ttl"`
	Snt      int           `json:"snt"`
	LastTime time.Duration `json:"last"`
	AvgTime  time.Duration `json:"avg"`
	BestTime time.Duration `json:"best"`
	WrstTime time.Duration `json:"worst"`
	Loss     float32       `json:"loss"`
}