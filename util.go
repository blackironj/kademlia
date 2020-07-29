package kademlia

import (
	"errors"
	"io/ioutil"
	"math/bits"
	"net/http"

	"crypto/rand"

	"github.com/minio/sha256-simd"
)

// Returned if a routing table query returns no results. This is NOT expected
// behaviour
var ErrLookupFailure = errors.New("failed to find any peer in table")

func XOR(a, b []byte) []byte {
	c := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func CommonPrefixLen(a, b []byte) int {
	return zeroPrefixLen(XOR(a, b))
}

// ConvertPeerID creates a DHT ID by hashing a Peer ID (Multihash)
func ConvertPeerID(id string) []byte {
	hash := sha256.Sum256([]byte(id))
	return hash[:]
}

func zeroPrefixLen(hashedID []byte) int {
	for i, b := range hashedID {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(hashedID) * 8
}

func GenerateRandBytes(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetMyIP() (string, error) {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)

	return string(content), nil
}
