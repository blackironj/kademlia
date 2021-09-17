package kademlia

import (
	"io/ioutil"
	"math/bits"
	"net/http"

	"github.com/google/uuid"
	"github.com/minio/sha256-simd"
)

const (
	_externalIPcheckerAddress = "https://myexternalip.com/raw"
)

func XOR(a, b []byte) []byte {
	c := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func CommonPrefixLen(a, b []byte) int {
	hashedID := XOR(a, b)
	for i, b := range hashedID {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(hashedID) * 8
}

// ConvertPeerID creates a DHT ID by hashing a Peer ID (Multihash)
func ConvertPeerID(id string) []byte {
	hash := sha256.Sum256([]byte(id))
	return hash[:]
}

func GetMyIP() (string, error) {
	resp, err := http.Get(_externalIPcheckerAddress)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)

	return string(content), nil
}

func NewUUIDv4() string {
	return uuid.NewString()
}
