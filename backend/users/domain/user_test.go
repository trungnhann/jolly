package domain

import (
	"testing"
)

func TestUser_UpdatePassword(t *testing.T) {
	id := UserUUID{}
	_ = id.UnmarshalText([]byte("019efe38-7f21-7265-a3da-abcb5ac2c1c9"))

	role := RoleCustomer()
	user, err := NewUser(id, "test@email.com", "Test Name", "initial_hash", role)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Run("successful password update", func(t *testing.T) {
		newHash := "new_secure_hash"
		err := user.UpdatePassword(newHash)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if user.PasswordHash() != newHash {
			t.Errorf("expected hash %q, got %q", newHash, user.PasswordHash())
		}
	})

	t.Run("empty password hash rejection", func(t *testing.T) {
		err := user.UpdatePassword("")
		if err == nil {
			t.Error("expected error for empty password hash, got nil")
		}
		if err != ErrPasswordHashEmpty {
			t.Errorf("expected ErrPasswordHashEmpty, got %v", err)
		}
	})
}
