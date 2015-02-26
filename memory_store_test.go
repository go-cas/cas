package cas

import (
	"testing"
)

func TestMemoryStore(t *testing.T) {
	user1 := &AuthenticationResponse{User: "user1"}
	user2 := &AuthenticationResponse{User: "user2"}
	store := &MemoryStore{}

	if err := store.Write("user1", user1); err != nil {
		t.Errorf("Expected store.Write(user1) to succeed, got error: %v", err)
	}

	if err := store.Write("user2", user2); err != nil {
		t.Errorf("Expected store.Write(user2) to succeed, got error: %v", err)
	}

	ar, err := store.Read("user2")
	if err != nil {
		t.Errorf("Expected store.Read(user2) to succeed, got error: %v", err)
	}

	if ar != user2 {
		t.Errorf("Expected retrieved AuthenticationResponse to be %v, got %v", user2, ar)
	}

	if err := store.Clear(); err != nil {
		t.Errorf("Expected store.Clear() to succeed, got error: %v", err)
	}

	_, err = store.Read("user1")
	if err == nil {
		t.Errorf("Expected an error from store.Read(user1), got nil")
	}

	if err != ErrInvalidTicket {
		t.Errorf("Expected ErrInvalidTicket from store.Read(user1), got %v", err)
	}
}
