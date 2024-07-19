package utils

import "github.com/google/uuid"

// Contains checks if a slice contains a specific item
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Remove removes a specific item from a slice
func Remove(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// GenerateID generates a new UUID as a string
func GenerateID() string {
	return uuid.NewString()
}
