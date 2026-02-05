//go:build integration

package calendar

import (
	"net/http"
	"testing"
	"time"
)

type TestClient struct{}

func (c *TestClient) MakeRequest(url string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	return client.Do(req)
}

func (c *TestClient) MakeRequestWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}

func TestParseCalendarWithTeamIDs(t *testing.T) {
	parser := NewParser(&TestClient{})

	url := "https://pfo.fhr.ru/tournaments/pervenstvo-pfo-18171615-let-16735091/"
	matches, err := parser.Parse(url)
	if err != nil {
		t.Fatalf("Error parsing calendar: %v", err)
	}

	t.Logf("Parsed %d matches", len(matches))

	// Count matches with team IDs
	withHomeID := 0
	withAwayID := 0

	for i, m := range matches {
		if m.HomeTeam.ID != "" {
			withHomeID++
		}
		if m.AwayTeam.ID != "" {
			withAwayID++
		}

		// Log first 5 matches
		if i < 5 {
			t.Logf("Match %d:", i+1)
			t.Logf("  Home: ID=%s, Name=%s", m.HomeTeam.ID, m.HomeTeam.Name)
			t.Logf("  Away: ID=%s, Name=%s", m.AwayTeam.ID, m.AwayTeam.Name)
		}
	}

	t.Logf("\nSummary:")
	t.Logf("  Matches with home team ID: %d/%d (%.1f%%)", withHomeID, len(matches), 100.0*float64(withHomeID)/float64(len(matches)))
	t.Logf("  Matches with away team ID: %d/%d (%.1f%%)", withAwayID, len(matches), 100.0*float64(withAwayID)/float64(len(matches)))

	// Verify we got team IDs
	if withHomeID == 0 {
		t.Error("No matches have home team ID extracted")
	}
	if withAwayID == 0 {
		t.Error("No matches have away team ID extracted")
	}
}

func TestExtractTeamIDFromLogo(t *testing.T) {
	tests := []struct {
		logoURL  string
		expected string
	}{
		{"/upload/team_logo/658725.png", "658725"},
		{"/upload/upload-webp/upload/team_logo/658721-70.webp", "658721"},
		{"/upload/team_logo/123456-70.webp", "123456"},
		{"/images/no-team-photo.png", ""},
	}

	for _, tt := range tests {
		t.Run(tt.logoURL, func(t *testing.T) {
			result := extractTeamIDFromLogo(tt.logoURL)
			if result != tt.expected {
				t.Errorf("extractTeamIDFromLogo(%q) = %q, want %q", tt.logoURL, result, tt.expected)
			}
		})
	}
}

func TestExtractTeamIDFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"/teams/lada_651237/", "651237"},
		{"/tournaments/champ-2024/ak-bars_658725/", "658725"},
		{"/teams/spartak-moscow_123456", "123456"},
		{"/teams/without-id/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := extractTeamIDFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("extractTeamIDFromURL(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}
