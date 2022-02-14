package syncrand

import (
	"crypto/hmac"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/sha3"
)

// Commit computes a 128-bit commitment of a nonce to send to the remote.
func Commit(nonce []byte) []byte {
	sh := sha3.NewShake128()
	sh.Write([]byte("syncrand:nonce:"))
	sh.Write(nonce[:])
	var buf [32]byte
	if _, err := io.ReadFull(sh, buf[:]); err != nil {
		panic("syncrand: read error: " + err.Error())
	}
	return buf[:]
}

// Verify verifies that their commitment matches their nonce, and also their commitment is not the same as our commitment.
func Verify(ourCommitment []byte, theirCommitment []byte, theirNonce []byte) bool {
	return !hmac.Equal(ourCommitment, theirCommitment) && hmac.Equal(theirCommitment, Commit(theirNonce))
}

// MakeSeed combines nonces into a seed appropriate to pass to NewRand. All nonces must be of the same length.
func MakeSeed(nonces ...[]byte) []byte {
	seed := make([]byte, len(nonces[0]))
	for _, nonce := range nonces {
		if len(nonce) != len(seed) {
			panic("syncrand: MakeSeed passed nonces of different lengths")
		}
		for i := range seed {
			seed[i] ^= nonce[i]
		}
	}
	return seed
}

// Source is a rand.Source compatible random number generator.
type Source struct {
	sh     sha3.ShakeHash
	offset uint
}

// NewRand creates a new random number generator with the given seed.
func NewSource(seed []byte) *Source {
	r := &Source{sha3.NewShake128(), 0}
	r.sh.Write([]byte("syncrand:seed:"))
	r.sh.Write(seed[:])
	return r
}

// Clone clones this random number generator.
func (r *Source) Clone() *Source {
	return &Source{r.sh.Clone(), r.offset}
}

// Int63 implements Source.Int63.
func (r *Source) Int63() int64 {
	var buf [8]byte
	if _, err := io.ReadFull(r.sh, buf[:]); err != nil {
		panic("syncrand: read error: " + err.Error())
	}
	r.offset++
	return int64(binary.LittleEndian.Uint64(buf[:]) & (1<<63 - 1))
}

// Seed is not supported: you must create a new Source to reseed.
func (r *Source) Seed(seed int64) {
	panic("cannot reseed this rng")
}

// SeedOffset returns how many times Int63 has been called on this random number generator.
func (r *Source) SeedOffset() uint {
	return r.offset
}
