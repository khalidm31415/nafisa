package elasticsarch_helper

import "strings"

// Utility function to convert last_education to terms array based on the provided order
func getLastEducationTerms(minEducation string) []string {
	educations := []string{"SD", "SMP", "SMA", "D3", "D4", "S1", "S2", "S3"}
	index := -1
	for i, education := range educations {
		if strings.EqualFold(education, minEducation) {
			index = i
			break
		}
	}
	if index < 0 {
		return educations
	}
	return educations[index:]
}
