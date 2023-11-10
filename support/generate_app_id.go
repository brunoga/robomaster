package support

import (
	"crypto/rand"
	"math/big"
)

// GenerateAppID generates a random app ID suitable to be used for Robomaster
// applications.
func GenerateAppID() (uint64, error) {
	// Probably overkill to use crypto/rand here but it also does not hurt.
	n, err := rand.Int(rand.Reader, new(big.Int).SetUint64(^uint64(0)))
	if err != nil {
		return 0, err
	}

	// Create an app ID out of the first 8 bytes of the UUID.
	return n.Uint64(), nil
}
