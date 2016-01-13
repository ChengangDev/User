package sail

import (
	"io/ioutil"
	"net/http"
)

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

func GetRequest(url string, header *map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if header != nil {
		for key, value := range *header {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	sRet, err := ioutil.ReadAll(resp.Body)
	return string(sRet), err
}
