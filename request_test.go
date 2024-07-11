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

type getResponse struct {
	Args    map[string]string `json:"args"`
	Headers struct {
		Accept          string `json:"Accept"`
		AcceptEncoding  string `json:"Accept-Encoding"`
		AcceptLanguage  string `json:"Accept-Language"`
		Dnt             string `json:"Dnt"`
		Host            string `json:"Host"`
		Priority        string `json:"Priority"`
		Referer         string `json:"Referer"`
		SecChUa         string `json:"Sec-Ch-Ua"`
		SecChUaMobile   string `json:"Sec-Ch-Ua-Mobile"`
		SecChUaPlatform string `json:"Sec-Ch-Ua-Platform"`
		SecFetchDest    string `json:"Sec-Fetch-Dest"`
		SecFetchMode    string `json:"Sec-Fetch-Mode"`
		SecFetchSite    string `json:"Sec-Fetch-Site"`
		UserAgent       string `json:"User-Agent"`
		XAmznTraceId    string `json:"X-Amzn-Trace-Id"`
	} `json:"headers"`
	Origin string `json:"origin"`
	Url    string `json:"url"`
}

func TestGet(t *testing.T) {
	response, err := Get[getResponse]("https://httpbin.org/get", WithParameter("key", "value"))
	if err != nil {
		t.Errorf("failed to get JSON: %v", err)
	}

	if response.Args["key"] != "value" {
		t.Errorf("expected key to be 'value', got %s", response.Args["key"])
	}
}

func TestGet_JSON(t *testing.T) {
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
