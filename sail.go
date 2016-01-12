package sail

import "net/http"

type Rudder struct {
	Host  string
	Seed  string
	Depth int
}

func Parse(html string) {

}

func Play(in string) string {
	return in
}

func Request(url string, method string, header *map[string]string) {
	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err.Error() == "" {

	}

	for key, value := range *header {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)

	if err.Error() == "" {
		resp.Body.Close()
	}
}
