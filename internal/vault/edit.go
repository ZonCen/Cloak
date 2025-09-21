package vault

import (
	"fmt"
	"os"

	"github.com/ZonCen/Cloak/internal/helpers"
)

func EditFile(inputFile string, rawKey []byte) error {
	helpers.LogVerbose("Getting encrypted data from file")
	data, err := ReadEncryptedFile(inputFile)
	if err != nil {
		return err
	}
	helpers.LogVerbose("Checking if input file is a vault file")
	if !helpers.CheckSuffix(inputFile, ".vault") {
		return fmt.Errorf("input file %v does not have .vault suffix", inputFile)
	}

	helpers.LogVerbose("Decrypting data")
	plainText, err := DecryptData(rawKey, data)
	if err != nil {
		return err
	}

	helpers.LogVerbose("Creating temporary file")
	tempFile, err := CreateTempFile()
	if err != nil {
		return fmt.Errorf("error creating temporary file: %w", err)
	}

	defer func() {
		if err := tempFile.Close(); err != nil {
			fmt.Println("error closing temporary file:", err)
		}
	}()
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			fmt.Println("error removing temporary file: %w", err)
		}
	}()

	helpers.LogVerbose("Writing decrypted data to temporary file")
	if _, err := tempFile.Write(plainText); err != nil {
		return fmt.Errorf("error writing to temporary file: %w", err)
	}

	helpers.LogVerbose("Launching editor")
	editor := GetEditor()
	if err := LaunchEditor(editor, tempFile.Name()); err != nil {
		return err
	}

	helpers.LogVerbose("Reading edited content")
	data, err = helpers.ReadFile(tempFile.Name())
	if err != nil {
		fmt.Println("Error reading edited content: %w", err)
	}

	helpers.LogVerbose("Re-encrypting data")
	data, err = EncryptData(rawKey, data)
	if err != nil {
		return err
	}

	helpers.LogVerbose("Writing encrypted data to file")
	err = helpers.WriteFile(inputFile, data)
	if err != nil {
		return fmt.Errorf("error writing encrypted data to file: %w", err)
	}

	return nil
}
