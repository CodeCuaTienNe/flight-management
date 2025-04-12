package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// ValidationService implements validation utilities
type ValidationService struct {
	reader *bufio.Reader
}

// NewValidationService creates a new ValidationService
func NewValidationService() *ValidationService {
	return &ValidationService{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ValidateFlightNumber checks if a flight number matches the required format (Fxxxx)
func (v *ValidationService) ValidateFlightNumber(flightNumber string) bool {
	pattern := regexp.MustCompile(`^F\d{4}$`)
	return pattern.MatchString(flightNumber)
}

// ValidateDates checks if the departure and arrival dates are valid
func (v *ValidationService) ValidateDates(departureTime, arrivalTime time.Time) bool {
	currentTime := time.Now()

	// Departure must be at least 3 hours in the future
	if departureTime.Before(currentTime.Add(3 * time.Hour)) {
		fmt.Println("Departure time must be at least 3 hours from now.")
		return false
	}

	// Arrival must be after departure
	if arrivalTime.Before(departureTime) || arrivalTime.Equal(departureTime) {
		fmt.Println("Arrival time must be after departure time.")
		return false
	}

	// Calculate duration between departure and arrival
	duration := arrivalTime.Sub(departureTime)

	// Duration must be between 30 minutes and 24 hours for commercial flights
	minDuration := 30 * time.Minute
	maxDuration := 24 * time.Hour

	if duration < minDuration || duration > maxDuration {
		fmt.Printf("Flight duration must be between 30 minutes and 24 hours. Current: %v\n", duration)
		return false
	}

	return true
}

// CheckYesOrNo prompts the user with a yes/no question and returns the result
func (v *ValidationService) CheckYesOrNo(prompt string) bool {
	for {
		fmt.Print(prompt)
		input, _ := v.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)

		if input == "Y" || input == "YES" {
			return true
		} else if input == "N" || input == "NO" {
			return false
		}

		fmt.Println("Please enter 'Y' for yes or 'N' for no.")
	}
}

// GetInteger prompts for an integer input within a specified range
func (v *ValidationService) GetInteger(prompt string, errorMsg string, min, max int) int {
	for {
		fmt.Print(prompt)
		var input int
		_, err := fmt.Scanf("%d\n", &input)

		if err != nil {
			fmt.Println(errorMsg)
			v.reader.ReadString('\n') // Clear the input buffer
			continue
		}

		if input < min || input > max {
			fmt.Printf("Input must be between %d and %d.\n", min, max)
			continue
		}

		return input
	}
}

// GetString prompts for a string input
func (v *ValidationService) GetString(prompt string, errorMsg string, allowEmpty bool) string {
	for {
		fmt.Print(prompt)
		input, _ := v.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && !allowEmpty {
			fmt.Println(errorMsg)
			continue
		}

		return input
	}
}

// GetDate prompts for a date input in the specified format
func (v *ValidationService) GetDate(prompt string, errorMsg string, format string, allowEmpty bool) time.Time {
	for {
		fmt.Print(prompt)
		input, _ := v.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && allowEmpty {
			return time.Time{}
		}

		date, err := time.Parse(format, input)
		if err != nil {
			fmt.Println(errorMsg)
			continue
		}

		return date
	}
}

// GetLong prompts for a long integer input
func (v *ValidationService) GetLong(prompt string, errorMsg string, allowEmpty bool) int64 {
	for {
		fmt.Print(prompt)
		var input int64
		_, err := fmt.Scanf("%d\n", &input)

		if err != nil {
			if allowEmpty {
				return 0
			}
			fmt.Println(errorMsg)
			v.reader.ReadString('\n') // Clear the input buffer
			continue
		}

		return input
	}
}