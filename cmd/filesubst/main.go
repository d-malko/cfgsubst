package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	var filename string
	flag.StringVar(&filename, "f", "", "file to read")    //map -f flag to filename variable
	flag.StringVar(&filename, "file", "", "file to read") //map --file flag to filename variable

	// Parse the command-line arguments
	flag.Parse()

	if filename == "" {
		log.Fatalf("Please provide a filename via -f or --file")
	}

	// Read configuration file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var lines []string

	// Create a new Scanner for the file
	scanner := bufio.NewScanner(file)

	// Iterate through every line in file and substitute values if necessary
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2) // Split the line into key-value pair
		if len(split) < 2 {
			log.Fatalf("Invalid line: %s", line)
		}
		key, value := split[0], split[1]

		// If environment variable exists, substitute. Otherwise, leave original value.
		if envValue, exists := os.LookupEnv(key); exists {
			lines = append(lines, key+"="+envValue)
		} else {
			lines = append(lines, key+"="+value)
		}
	}

	// Check for Scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	// Write back to file
	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(filename, []byte(newContent), 0)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	log.Println("Configuration file updated successfully.")
}
