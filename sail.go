package sail

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var DefaultHeader = map[string]string{
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/47.0.2526.73 Chrome/47.0.2526.73 Safari/537.36",
	"Cookie":     "s=1iyy12llah; xq_a_token=ea8f61e35ef1ad1c8fda1747b148e209e8de179c; xq_r_token=4f20dec99d9204c7ec4f26708feafd600c9cd802; __utma=1.1311703092.1452697027.1452697027.1452697027.1; __utmb=1.1.10.1452697027; __utmc=1; __utmz=1.1452697027.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmt=1; Hm_lvt_1db88642e346389874251b5a1eded6e3=1452697028; Hm_lpvt_1db88642e346389874251b5a1eded6e3=145269702",
}

type Rudder struct {

	//pattern to get number of related users
	CountPattern string
	//pattern to get current page number
	PageNoPattern string
	//pattern to get page size
	//PageSizePattern string
	//pattern to get number of pages
	PageCountPattern string

	//pattern to get user id
	IDPattern string
	//pattern to get user name
	NamePattern string
	//other patterns to get user info
	OtherPattern map[string]string
}

func Parse(resp string, rudder *Rudder) (err error) {
	if rudder == nil {
		return errors.New("Rudder is nil.")
	}

	if rudder.CountPattern == "" {
		return errors.New("Rudder.CountPattern is empty.")
	}
	return
}

func Play(in string) string {
	return in
}

func GetRequest(url string, header *map[string]string) (string, error) {
	fmt.Println("Get url: ", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	if header != nil {
		fmt.Println("Header:")
		for key, value := range *header {
			req.Header.Add(key, value)
			fmt.Println(" ", key, ":", value)
		}
	} else {
		fmt.Println("No header for request.")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	sRet, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}
	fmt.Printf("sRet: %q\n", sRet)
	return string(sRet), err
}
