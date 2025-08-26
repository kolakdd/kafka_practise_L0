package cache

import (
	"bytes"
	"testing"
)

func TestLRUCacheCapacity(t *testing.T) {
	cacheRepo := NewCacheRepo(3)

	cacheRepo.Set(1, []byte("a"))
	cacheRepo.Set(2, []byte("b"))
	cacheRepo.Set(3, []byte("c"))
	cacheRepo.Set(4, []byte("d"))
	cacheRepo.Set(5, []byte("e"))
	cacheRepo.Set(6, []byte("f"))
	cacheRepo.Set(7, []byte("g"))
	cacheRepo.Set(8, []byte("h"))
	cacheRepo.Set(9, []byte("i"))
	cacheRepo.Set(10, []byte("j"))

	expected := "j | i | h"
	msg := cacheRepo.Debug()
	if msg != expected {
		t.Errorf(`Test %q, want %q, error`, msg, expected)
	}
}

func TestLRUCacheGet(t *testing.T) {
	cacheRepo := NewCacheRepo(3)

	cacheRepo.Set(1, []byte("a"))
	cacheRepo.Set(2, []byte("b"))

	expected := []byte("a")
	msg, _ := cacheRepo.Get(1, true)
	if !bytes.Equal(msg, expected) {
		t.Errorf(`Test %q, want %q, error`, msg, expected)
	}
}

func TestLRUCacheGetWrong(t *testing.T) {
	cacheRepo := NewCacheRepo(3)

	cacheRepo.Set(1, []byte("a"))
	cacheRepo.Set(2, []byte("b"))

	expected := []byte("a")
	msg, _ := cacheRepo.Get(2, true)
	if bytes.Equal(msg, expected) {
		t.Errorf(`Test %q, want %q, error`, msg, expected)
	}
}
