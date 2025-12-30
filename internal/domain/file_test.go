package domain

import (
	"testing"
	"time"
)

func TestNewFile_SetsCreatedAt(t *testing.T) {
	before := time.Now()

	file, err := NewFile(
		"1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		VisibilityPrivate,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	after := time.Now()

	if file.CreatedAt().Before(before) || file.CreatedAt().After(after) {
		t.Errorf("createdAt should be set to now")
	}
}

func TestNewFile_StatusIsPending(t *testing.T) {
	file, err := NewFile(
		"1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		VisibilityPrivate,
	)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if file.Status() != StatusPending {
		t.Errorf("expected status pending")
	}
	if file.StorageKey() != "" {
		t.Errorf("expected empty storageKey")
	}
	if file.ID() == "" {
		t.Errorf("expected non-empty ID")
	}
	if file.OwnerID() != "user-1" {
		t.Errorf("ownerID mismatch")
	}
	if file.FileName() != "file.txt" {
		t.Errorf("fileName mismatch")
	}
	if file.MimeType() != "text/plain" {
		t.Errorf("mimeType mismatch")
	}
	if file.Size() != 100 {
		t.Errorf("size mismatch")
	}
	if file.Visibility() != VisibilityPrivate {
		t.Errorf("visibility mismatch")
	}
	if file.ProjectID() != nil {
		t.Errorf("expected nil projectID")
	}
	if time.Since(file.CreatedAt()) < 0 {
		t.Errorf("createdAt mismatch")
	}
}

func TestRehydrateFile_RequiresCreatedAt(t *testing.T) {
	_, err := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		VisibilityPrivate,
		StatusAvailable,
		time.Time{},
	)

	if err == nil {
		t.Fatalf("expected error for zero createdAt")
	}
}

func TestRehydrateFile_InvalidVisibility(t *testing.T) {
	_, err := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		"INVALID",
		StatusAvailable,
		time.Now(),
	)
	if err == nil {
		t.Fatalf("expected error for invalid visibility")
	}
}

func TestRehydrateFile_InvalidStatus(t *testing.T) {
	_, err := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		VisibilityPrivate,
		"INVALID",
		time.Now(),
	)
	if err == nil {
		t.Fatalf("expected error for invalid status")
	}
}

func TestRehydrateFile_InvalidId(t *testing.T) {
	_, err := RehydrateFile(
		"",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		VisibilityPrivate,
		StatusAvailable,
		time.Time{},
	)

	if err == nil {
		t.Fatalf("expected error for zero createdAt")
	}
}

func TestRehydrateFile_InvalidFileName(t *testing.T) {
	_, err := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"",
		"text/plain",
		100,
		"s3/key",
		VisibilityPrivate,
		StatusAvailable,
		time.Time{},
	)
	if err == nil {
		t.Fatalf("expected error for zero createdAt")
	}
}

func TestRehydrateFile_InvalidSize(t *testing.T) {
	_, err := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		-10,
		"s3/key",
		VisibilityPrivate,
		StatusAvailable,
		time.Time{},
	)
	if err == nil {
		t.Fatalf("expected error for zero createdAt")
	}
}

func TestRehydrateFile_InvalidOwnerId(t *testing.T) {
	_, err := RehydrateFile(
		"file-1",
		"",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		VisibilityPrivate,
		StatusAvailable,
		time.Time{},
	)

	if err == nil {
		t.Fatalf("expected error for zero createdAt")
	}
}

func TestRehydrateFile_Success(t *testing.T) {
	now := time.Now()

	file, err := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		VisibilityPublic,
		StatusAvailable,
		now,
	)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if !file.CreatedAt().Equal(now) {
		t.Errorf("createdAt mismatch")
	}
}

func TestMarkAsAvailable_Success(t *testing.T) {
	file, _ := NewFile(
		"1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		VisibilityPrivate,
	)

	err := file.MarkAsAvailable("s3/key")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if file.Status() != StatusAvailable {
		t.Errorf("expected status available")
	}

	if file.StorageKey() != "s3/key" {
		t.Errorf("storageKey not set")
	}
}

func TestMarkAsAvailable_InvalidTransition(t *testing.T) {
	file, _ := RehydrateFile(
		"file-1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		"s3/key",
		VisibilityPrivate,
		StatusAvailable,
		time.Now(),
	)

	err := file.MarkAsAvailable("another/key")
	if err == nil {
		t.Fatalf("expected error for invalid status transition")
	}
}

func TestMarkAsAvailable_EmptyStorageKey(t *testing.T) {
	file, _ := NewFile(
		"1",
		"user-1",
		nil,
		"file.txt",
		"text/plain",
		100,
		VisibilityPrivate,
	)

	err := file.MarkAsAvailable("")
	if err == nil {
		t.Fatalf("expected error for empty storageKey")
	}
}
