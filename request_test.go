package request

import "testing"

type getJsonResponse struct {
	Slideshow struct {
		Author string `json:"author"`
		Date   string `json:"date"`
		Slides []struct {
			Title string   `json:"title"`
			Type  string   `json:"type"`
			Items []string `json:"items,omitempty"`
		} `json:"slides"`
		Title string `json:"title"`
	} `json:"slideshow"`
}

func TestGet(t *testing.T) {
	response, err := Get[getJsonResponse]("https://httpbin.org/json")
	if err != nil {
		t.Errorf("failed to get JSON: %v", err)
	}

	if response.Slideshow.Author != "Yours Truly" {
		t.Errorf("expected author to be 'Yours Truly', got %s", response.Slideshow.Author)
	}
}

func TestPost(t *testing.T) {
	_, err := Post[struct{}]("https://httpbin.org/status/404")
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err.Code != 404 {
		t.Errorf("expected status code 404, got %d", err.Code)
	}
}
