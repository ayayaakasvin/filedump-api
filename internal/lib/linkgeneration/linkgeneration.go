package linkgeneration

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	InvalidInput = "invalid input"
	WayTooLong   = "expiration time is longer than a week"
)

func GenerateRandomLink(fileuuid string, expDur time.Duration) (string, error) {
	if err := uuid.Validate(fileuuid); err != nil {
		return "", errors.New(InvalidInput)
	}

	if expDur <= 0 {
		return "", errors.New(InvalidInput)
	}

	if expDur > 7*24*time.Hour {
		return "", errors.New(WayTooLong)
	}

	var rndByte [24]byte
	if _, err := rand.Read(rndByte[:]); err != nil {
		return "", err
	}

	// Encode to base64 URL encoding for safe links
	return base64.RawURLEncoding.EncodeToString(rndByte[:]), nil
}
