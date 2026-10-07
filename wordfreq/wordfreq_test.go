package wordfreq

import (
	"os"
	"strings"
	"testing"
)

func TestEmptyInput(t *testing.T) {
	content := ""
	tmpfile, err := createTempFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	results, err := WordFreq(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(results) != 0 {
		t.Errorf("Expected empty results, got %d items", len(results))
	}
}

func TestPunctuationStripping(t *testing.T) {
	content := "Hello, world! Hello."
	tmpfile, err := createTempFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	results, err := WordFreq(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should have "hello" and "world" with counts 2 and 1 respectively
	expected := map[string]int{
		"hello": 2,
		"world": 1,
	}
	
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	
	for _, result := range results {
		if count, ok := expected[result.Word]; ok {
			if count != result.Count {
				t.Errorf("Expected count %d for word %s, got %d", count, result.Word, result.Count)
			}
			delete(expected, result.Word)
		} else {
			t.Errorf("Unexpected word %s with count %d", result.Word, result.Count)
		}
	}
	
	if len(expected) != 0 {
		t.Errorf("Missing expected words: %v", expected)
	}
}

func TestCaseInsensitivity(t *testing.T) {
	content := "Hello HELLO hello World WORLD"
	tmpfile, err := createTempFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	results, err := WordFreq(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// All should be counted as "hello" and "world"
	expected := map[string]int{
		"hello": 3,
		"world": 2,
	}
	
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	
	for _, result := range results {
		if count, ok := expected[result.Word]; ok {
			if count != result.Count {
				t.Errorf("Expected count %d for word %s, got %d", count, result.Word, result.Count)
			}
			delete(expected, result.Word)
		} else {
			t.Errorf("Unexpected word %s with count %d", result.Word, result.Count)
		}
	}
	
	if len(expected) != 0 {
		t.Errorf("Missing expected words: %v", expected)
	}
}

func TestAlphabeticalTieBreaking(t *testing.T) {
	content := "zebra apple banana apple zebra banana"
	tmpfile, err := createTempFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	results, err := WordFreq(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Both "apple" and "banana" and "zebra" have count 2
	// Should be sorted alphabetically: apple, banana, zebra
	if len(results) < 3 {
		t.Fatalf("Expected at least 3 results, got %d", len(results))
	}
	
	expectedOrder := []string{"apple", "banana", "zebra"}
	for i, expectedWord := range expectedOrder {
		if results[i].Word != expectedWord {
			t.Errorf("Expected word %s at position %d, got %s", expectedWord, i, results[i].Word)
		}
	}
}

func TestTopTenLimit(t *testing.T) {
	// Create content with more than 10 unique words
	var content strings.Builder
	for i := 0; i < 15; i++ {
		content.WriteString("word")
		content.WriteString(string(rune('a'+i)))
		content.WriteString(" ")
	}
	
	tmpfile, err := createTempFile(content.String())
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	results, err := WordFreq(tmpfile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Should only return top 10
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}
	
	// All should have count 1
	for _, result := range results {
		if result.Count != 1 {
			t.Errorf("Expected count 1 for all words, got %d for %s", result.Count, result.Word)
		}
	}
}

func createTempFile(content string) (*os.File, error) {
	tmpfile, err := os.CreateTemp("", "wordfreq_test")
	if err != nil {
		return nil, err
	}
	
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		return nil, err
	}
	
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}
	
	return tmpfile, nil
}