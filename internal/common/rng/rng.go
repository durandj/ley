package rng

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	// RNG gives a shared random number generator which can be used for
	// test data generation.
	//
	// The seed can be set ahead of time by giving an integer value for
	// the RANDOM_SEED environment variable.
	// nolint: gosec
	RNG = rand.New(rand.NewSource(getRandomSeed()))
)

func getRandomSeed() int64 {
	userSeed := os.Getenv("RANDOM_SEED")
	if userSeed != "" {
		parsedSeed, err := strconv.Atoi(userSeed)
		if err != nil {
			return time.Now().UnixMilli()
		}

		return int64(parsedSeed)
	}

	return time.Now().UnixMilli()
}

// RandomChoice picks a random element from the given values list.
func RandomChoice[Type any](values []Type) Type {
	return values[RNG.Intn(len(values))]
}
