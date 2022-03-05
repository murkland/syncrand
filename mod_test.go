package syncrand_test

import (
	"bytes"
	"testing"

	"github.com/murkland/syncrand"
)

func TestCommit(t *testing.T) {
	commitment := syncrand.Commit([]byte("hello"))
	if expected := []byte{171, 228, 27, 72, 163, 255, 25, 92, 68, 116, 250, 75, 192, 215, 100, 130, 254, 230, 109, 124, 53, 74, 223, 85, 206, 179, 238, 95, 98, 236, 215, 215}; !bytes.Equal(expected, commitment) {
		t.Errorf("syncrand.Commit(): expected %v, got %v", expected, commitment)
	}
}

func TestVerify(t *testing.T) {
	if !syncrand.Verify(syncrand.Commit([]byte("goodbye")), syncrand.Commit([]byte("hello")), []byte("hello")) {
		t.Errorf("syncrand.Verify(): expected ok, got not ok")
	}
}

func TestVerifySameCommitment(t *testing.T) {
	if syncrand.Verify(syncrand.Commit([]byte("hello")), syncrand.Commit([]byte("hello")), []byte("hello")) {
		t.Errorf("syncrand.Verify(): expected not ok, got ok")
	}
}

func TestVerifyBadCommitment(t *testing.T) {
	if syncrand.Verify(syncrand.Commit([]byte("goodbye")), syncrand.Commit([]byte("hello")), []byte("helloe")) {
		t.Errorf("syncrand.Verify(): expected not ok, got ok")
	}
}

func TestMakeSeed(t *testing.T) {
	seed := syncrand.MakeSeed([]byte("hello"), []byte("henlo"))
	if expected := []byte{0, 0, 2, 0, 0}; !bytes.Equal(expected, seed) {
		t.Errorf("syncrand.MakeSeed(): expected %v, got %v", expected, seed)
	}
}

func TestMakeSeedBadLengths(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("did not panic")
		}
	}()
	syncrand.MakeSeed([]byte("hello"), []byte("henloe"))
}

func TestSource(t *testing.T) {
	s := syncrand.NewSource([]byte("hello"))
	n := s.Int63()
	if expected := int64(6164488906953303115); n != expected {
		t.Errorf("syncrand.Source.Int63(): expected %v, got %v", expected, n)
	}
}

func TestSourceClone(t *testing.T) {
	s := syncrand.NewSource([]byte("hello"))
	s2 := s.Clone()
	s.Int63()
	n := s2.Int63()
	if expected := int64(6164488906953303115); n != expected {
		t.Errorf("syncrand.Source.Int63(): expected %v, got %v", expected, n)
	}
}

func TestSourceOffset(t *testing.T) {
	s := syncrand.NewSource([]byte("hello"))
	s.Int63()
	s.Int63()
	s.Int63()
	n := s.SeedOffset()
	if expected := uint(3); n != expected {
		t.Errorf("syncrand.Source.SeedOffset(): expected %v, got %v", expected, n)
	}
}

func TestSourceSeed(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("did not panic")
		}
	}()
	s := syncrand.NewSource([]byte("hello"))
	s.Seed(1)
}

func TestEndToEnd(t *testing.T) {
	nonce1 := []byte("hello")
	nonce2 := []byte("aloha")

	commitment1 := syncrand.Commit(nonce1)
	commitment2 := syncrand.Commit(nonce2)

	if !syncrand.Verify(commitment1, commitment2, nonce2) {
		t.Error("failed to verify commitment2")
	}

	if !syncrand.Verify(commitment2, commitment1, nonce1) {
		t.Error("failed to verify commitment1")
	}

	seed := syncrand.MakeSeed(nonce1, nonce2)
	s := syncrand.NewSource(seed)
	n := s.Int63()
	if expected := int64(4154673063121253460); n != expected {
		t.Errorf("syncrand.Source.Int63(): expected %v, got %v", expected, n)
	}
}
