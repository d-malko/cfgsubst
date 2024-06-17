package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	var filename, outputFilename string
	flag.StringVar(&filename, "f", "", "file to read")                       //map -f flag to filename variable
	flag.StringVar(&filename, "file", "", "file to read")                    //map --file flag to filename variable
	flag.StringVar(&outputFilename, "o", "", "file to write output to")      //map -o flag to outputFilename variable
	flag.StringVar(&outputFilename, "output", "", "file to write output to") //map --output flag to outputFilename variable

	// Parse the command-line arguments
	flag.Parse()

	if filename == "" {
		log.Fatalf("Please provide a filename via -f or --file")
	}

	if outputFilename == "" {
		log.Fatalf("Please provide an output filename via -o or --output")
	}

	log.Printf("Templating file %s to %s", filename, outputFilename)
	// Read configuration file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var lines []string

	// Create a new Scanner for the file
	scanner := bufio.NewScanner(file)

	// Define regex pattern
	pattern := "\\w+=.*"
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatalf("Failed to Compile regexp %v", err)
	}

	// Iterate through every line in file and substitute values if necessary
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			split := strings.SplitN(line, "=", 2) // Split the line into key-value pair
			if len(split) < 2 {
				log.Fatalf("Invalid line %s", line)
			}
			key, value := split[0], split[1]

			// If environment variable exists, substitute. Otherwise, leave original value.
			if envValue, exists := os.LookupEnv(key); exists {
				lines = append(lines, key+"="+envValue)
			} else {
				lines = append(lines, key+"="+value)
			}
		} else {
			// If line does not match the pattern, write it as-is.
			lines = append(lines, line)
		}
	}

	// Check for Scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error %v", err)
	}

	// Write to the output file
	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(outputFilename, []byte(newContent), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write to file %v", err)
	}

	log.Printf("Configuration file updated and written to %s successfully.", outputFilename)
}
