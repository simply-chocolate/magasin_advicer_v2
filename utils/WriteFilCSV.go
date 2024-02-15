package utils

import (
	"fmt"
	"os"
)

// SaveDataAsCSV saves the provided data to a CSV file locally.
// The fileName parameter is used as the name of the CSV file, and data is the content to write.
// The brandName parameter can be used to organize files into different directories or for naming conventions.
func SaveDataAsCSV(fileName string, data string, brandName string) error {
	// Define the local directory path based on the brandName.
	// You might want to adjust the directory structure or naming as per your requirements.
	var localDirectory string
	if brandName == "simply" || brandName == "SIMPLY" {
		localDirectory = "simply"
	} else if brandName == "magasin" || brandName == "MAGASIN" {
		localDirectory = "magasin"
	} else {
		return fmt.Errorf("unknown brand, unable to determine local directory for file: %s", fileName)
	}

	// Ensure the directory exists.
	err := os.MkdirAll(localDirectory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating local directory: %v", err)
	}

	// Create or open the file within the specified directory.
	filePath := fmt.Sprintf("./%s/%s", localDirectory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating the file on local disk: %s, %v", filePath, err)
	}
	defer file.Close()

	// Write the data to the file.
	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("error writing data to the file: %s, %v", filePath, err)
	}

	return nil
}
