package wordfreq

import (
	"bufio"
	"os"
	"strings"
	"unicode"
)

// WordFreqResult represents a word and its frequency count
type WordFreqResult struct {
	Word  string
	Count int
}

// WordFreq reads a text file and returns the 10 most frequent words
func WordFreq(filename string) ([]WordFreqResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	wordCounts := make(map[string]int)
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.FieldsFunc(line, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		
		for _, word := range words {
			// Convert to lowercase
			lowerWord := strings.ToLower(word)
			// Strip punctuation (already handled by FieldsFunc)
			wordCounts[lowerWord]++
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice for sorting
	var results []WordFreqResult
	for word, count := range wordCounts {
		results = append(results, WordFreqResult{Word: word, Count: count})
	}

	// Sort by count (descending) and then by word (ascending) for ties
	sortedResults := mergeSort(results)
	
	// Return top 10
	if len(sortedResults) > 10 {
		return sortedResults[:10], nil
	}
	return sortedResults, nil
}

// mergeSort sorts the results by count (descending) and then by word (ascending)
func mergeSort(results []WordFreqResult) []WordFreqResult {
	if len(results) <= 1 {
		return results
	}
	
	mid := len(results) / 2
	left := mergeSort(results[:mid])
	right := mergeSort(results[mid:])
	
	return merge(left, right)
}

// merge merges two sorted slices
func merge(left, right []WordFreqResult) []WordFreqResult {
	result := make([]WordFreqResult, 0, len(left)+len(right))
	i, j := 0, 0
	
	for i < len(left) && j < len(right) {
		// Sort by count descending, then by word ascending
		if left[i].Count > right[j].Count || 
		   (left[i].Count == right[j].Count && left[i].Word < right[j].Word) {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	
	// Add remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	
	return result
}