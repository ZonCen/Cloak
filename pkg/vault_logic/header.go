package vault_logic

import "fmt"

func BuildHeader() []byte {
	var header FileHeader
	copy(header.Magic[:], []byte("CLOK"))
	header.Version = 1

	buf := make([]byte, 0, 5)
	buf = append(buf, header.Magic[:]...)
	buf = append(buf, header.Version)

	return buf
}

func ValidateHeader(header []byte) error {
	if string(header[:4]) != "CLOK" {
		return fmt.Errorf("invalid file format: missing CLOK magic bytes")
	}
	version := header[4]
	if version != 1 {
		return fmt.Errorf("unsupported file version: %d", version)
	}
	return nil
}
