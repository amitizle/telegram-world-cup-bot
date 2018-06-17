package world_cup_http_client

import (
	"fmt"
	"testing"
)

func TestPrepareURL(t *testing.T) {
	testBaseURL := "http://www.wut.com"
	testTable := []struct {
		Name        string
		Path        string
		QueryParams map[string]string
		Expected    string
	}{
		{"simple path no query string", "a/b/c", map[string]string{}, fmt.Sprintf("%s/%s", testBaseURL, "a/b/c")},
		{"simple path and query string",
			"a/b", map[string]string{"encrypted": "true"},
			fmt.Sprintf("%s/%s?%s=%s", testBaseURL, "a/b", "encrypted", "true"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			httpClient, _ := NewClient(testBaseURL)
			url := httpClient.prepareURL(testCase.Path, testCase.QueryParams)
			if url != testCase.Expected {
				t.Fatalf("expected %v, got %v", testCase.Expected, url)
			}
		})
	}
}
