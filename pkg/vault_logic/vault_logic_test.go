package vault_logic

import (
	"crypto/cipher"
	"testing"
)

var (
	rawKey      = []byte("12345678901234567890123456789012")
	longRawKey  = []byte("123456789012345678901234567890123")
	shortRawKey = []byte("1234567890123456789012345678901")
)

func TestCreateBlock(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "Create working block",
			key:     rawKey,
			wantErr: false,
		},
		{
			name:    "Create non working block (too long)",
			key:     longRawKey,
			wantErr: true,
		},
		{
			name:    "Create non working block (too short)",
			key:     shortRawKey,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, err := CreateBlock(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && block == nil {
				t.Fatalf("CreateBlock() returned nil block for key")
			}
		})
	}
}

func TestCreateCipher(t *testing.T) {
	workingBlock, err := CreateBlock(rawKey)
	if err != nil {
		t.Fatalf("Failed to create a working block: %v", err)
	}

	tests := []struct {
		name    string
		block   cipher.Block
		wantErr bool
	}{
		{
			name:    "Create working cipher",
			block:   workingBlock,
			wantErr: false,
		},
		{
			name:    "Create non working cipher",
			block:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if tt.block == nil {
			// 	defer func() {
			// 		if r := recover(); r == nil {
			// 			t.Error("expected panic when passing nil block, but got none")
			// 		}
			// 	}()
			// }

			gcm, err := CreateCipher(tt.block)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if (gcm == nil) != tt.wantErr {
				t.Fatalf("CreateCipher() returned nil AEAD")
			}
		})
	}
}

func TestGenerateNonceAndRandomByteKey(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "Generate Nonce and ByteKey (1)",
			size:    1,
			wantErr: false,
		},
		{
			name:    "Generate Nonce and ByteKey (10)",
			size:    10,
			wantErr: false,
		},
		{
			name:    "Generate Nonce and ByteKey (100)",
			size:    100,
			wantErr: false,
		},
		{
			name:    "Generate Nonce and ByteKey (0)",
			size:    0,
			wantErr: true,
		},
		{
			name:    "Generate Nonce and ByteKey (-1)",
			size:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, err := GenerateNonce(tt.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if (nonce == nil) != tt.wantErr {
				t.Fatalf("GenerateNonce() returned nil value")
			}

			key, err := GenerateRandomByteKey(tt.size)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if (key == nil) != tt.wantErr {
				t.Fatalf("GenerateRandomByteKey() returned nil value")
			}
		})
	}
}
