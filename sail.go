package sail

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var DefaultHeader = map[string]string{
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/47.0.2526.73 Chrome/47.0.2526.73 Safari/537.36",
	"Cookie":     "s=1iyy12llah; __utma=1.1311703092.1452697027.1452697027.1452697027.1; __utmz=1.1452697027.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); Hm_lvt_1db88642e346389874251b5a1eded6e3=1452697028; xq_a_token=5ca726f28df667ac9a6a0e5ff22f2cd2a6fd7f8d; xq_r_token=79470dac5976ff342dc4accc39cacaa265ef4604",
}

//since regexp does not support lookbehind/ahead, use multi patterns instead
type Rudder struct {

	//pattern to get number of related users
	CountPatterns []string
	//pattern to get current page number
	PageNoPatterns []string
	//pattern to get page size
	//PageSizePattern string
	//pattern to get number of pages
	PageCountPatterns []string

	//pattern to get user id
	IDsPatterns []string
	//pattern to get user name
	NamesPatterns []string
	//other patterns to get user info
	OtherPatterns     []string
	OtherListPatterns []string
}

type Seed struct {
	FixedFormater string
	ID            string
	PageNo        int
	PageSize      int
	Host          string
	Depth         int
	Thread        int
	Interval      float32 //interval between each request
}

func (s *Seed) GetUrl() (url string) {
	url = fmt.Sprintf(s.FixedFormater, s.ID, s.PageNo, s.PageSize)
	return url
}

func Parse(resp string, rudder *Rudder) (pageNo int, pageCount int, ids []string, err error) {
	if rudder == nil {
		return pageNo, pageCount, ids, errors.New("Rudder is nil.")
	}

	count, err := parseSingle(resp, rudder.CountPatterns)
	if err != nil {
		return
	}
	fmt.Println("Parse Count:" + count)

	sPageNo, err := parseSingle(resp, rudder.PageNoPatterns)
	if err != nil {
		return
	}
	fmt.Println("Parse PageNo:", pageNo)
	pageNo, err = strconv.Atoi(sPageNo)

	sPageCount, err := parseSingle(resp, rudder.PageCountPatterns)
	if err != nil {
		return
	}
	fmt.Println("Parse PageCount:", pageCount)
	pageCount, err = strconv.Atoi(sPageCount)

	ids, err = parseList(resp, rudder.IDsPatterns)
	if err != nil {
		return
	}
	fmt.Println("Parse IDs:", ids)
	//skip names

	//skip other

	return
}

func parseSingle(in string, patterns []string) (out string, err error) {
	out = in
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}

		out = re.FindString(out)
		fmt.Println("ParseSingle:", pattern, out)
	}
	return out, nil
}

func parseList(in string, patterns []string) (out []string, err error) {
	if len(patterns) < 2 {
		return out, errors.New("List Patterns need 2+ patterns.")
	}
	re, err := regexp.Compile(patterns[0])
	if err != nil {
		return out, err
	}
	all := re.FindAllString(in, -1)

	for _, one := range all {
		for i := 1; i < len(patterns); i++ {
			pattern := patterns[i]
			re, err := regexp.Compile(pattern)
			if err != nil {
				return out, nil
			}

			one = re.FindString(one)
		}
		out = append(out, one)
	}
	return out, nil
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

func Users(s *Seed, r *Rudder) (l []string, err error) {
	if s.Depth < 0 {
		return l, errors.New(
			fmt.Sprintf("Stop at max depth:%v", s.Depth))
	}

	for {
		url := s.GetUrl()
		resp, err := GetRequest(url, &DefaultHeader)
		if err != nil {
			return l, err
		}
		pageNo, pageCount, ids, err := Parse(resp, r)

		fmt.Printf("PageNo:%v, PageCount:%v, IDs:%v\n", pageNo, pageCount, ids)
		s.PageNo++
		if s.PageNo > pageCount {
			break
		}
		time.Sleep(time.Second)
	}

	return l, err
}
