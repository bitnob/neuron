package utils

import (
	"testing"
)

func TestCrypto_Hash(t *testing.T) {
	crypto := NewCrypto([]byte("test-key"))
	tests := []struct {
		name     string
		input    string
		wantHash bool // true if we expect a non-empty hash
	}{
		{
			name:     "valid input",
			input:    "test-password",
			wantHash: true,
		},
		{
			name:     "empty input",
			input:    "",
			wantHash: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := crypto.Hash(tt.input)
			if (got != "") != tt.wantHash {
				t.Errorf("Hash() = %v, want non-empty = %v", got, tt.wantHash)
			}
		})
	}
}

func TestCrypto_HashPassword(t *testing.T) {
	crypto := NewCrypto([]byte("test-key"))
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid password",
			input:   "test-password",
			wantErr: false,
		},
		{
			name:    "empty password",
			input:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := crypto.HashPassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Error("HashPassword() returned empty hash")
			}
		})
	}
}

func BenchmarkCrypto_Hash(b *testing.B) {
	crypto := NewCrypto([]byte("test-key"))
	input := "benchmark-password"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		crypto.Hash(input)
	}
}
