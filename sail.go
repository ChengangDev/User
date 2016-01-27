package sail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	//	"reflect"
	"regexp"
	"strconv"
	"time"
)

var DefaultHeader = map[string]string{
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/47.0.2526.73 Chrome/47.0.2526.73 Safari/537.36",
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
	//pattern to get followers count
	FollowersCountPatterns []string
	//other patterns to get user info
	OtherPatterns     []string
	OtherListPatterns []string
}

type UserInfo map[string]string

type Seed struct {
	FixedFormater string
	ID            string
	PageNo        int
	PageSize      int
	Depth         int
	Thread        int
	Interval      float32 //interval between each request

}

func (s *Seed) GetUrl() (url string) {
	url = fmt.Sprintf(s.FixedFormater, s.ID, s.PageNo, s.PageSize)
	return url
}

//old parse method
func Parse(resp string, rudder *Rudder) (pageNo int, pageCount int, users map[string][]string, err error) {
	if rudder == nil {
		return pageNo, pageCount, users, errors.New("Rudder is nil.")
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

	users = make(map[string][]string)
	ids, err := parseList(resp, rudder.IDsPatterns)
	if err != nil {
		return
	}
	users["ids"] = ids
	fmt.Println("Parse IDs:", ids)

	//parse names
	names, err := parseList(resp, rudder.NamesPatterns)
	if err != nil {
		return
	}
	users["names"] = names
	fmt.Println("Parse Names:", names)

	//parse followers count
	counts, err := parseList(resp, rudder.FollowersCountPatterns)
	if err != nil {
		return
	}
	users["counts"] = counts
	//skip other

	return
}

//parse single text
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

//parse list text and return list
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

//raw header
func GetRequest(url string, header map[string]string) (string, error) {
	fmt.Println("Get url: ", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	if header != nil {
		fmt.Println("Header:")
		for key, value := range header {
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

//get cookie from main page
func GetCookie(url string) (cj *cookiejar.Jar, err error) {
	cj, err = cookiejar.New(nil)
	cli := http.Client{Jar: cj}
	resp, err := cli.Get(url)
	if err != nil {
		log.Fatalln("GetCookie():", err.Error())
		return cj, nil
	}

	sc := resp.Header.Get("Set-Cookie")
	log.Println("GetCookie():", sc)
	return
}

//
func GetRequestByCookie(url string, cj *cookiejar.Jar) (string, error) {
	log.Println("Get url: ", url)

	client := &http.Client{Jar: cj}

	resp, err := client.Get(url)
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
	//fmt.Printf("sRet: %q\n", sRet)
	return string(sRet), err
}

func ParseJson(resp string, out *map[string]interface{}) (err error) {
	b := []byte(resp)
	json.Unmarshal(b, out)

	//log.Println(out)
	return
}

//fetch all followers of the seed
func FetchFollowers(s *Seed, ch chan UserInfo) (err error) {
	log.Println("FetchFollowers:Start Fetch Followers.")
	if s.Depth < 0 {
		close(ch)
		fmt.Println("Depth is less zero. Fetch over.")
	}

	cj, err := GetCookie("http://xueqiu.com")
	for {
		url := s.GetUrl()
		resp, err := GetRequestByCookie(url, cj)
		if err != nil {
			fmt.Print(err.Error())
			return err
		}

		m := map[string]interface{}{}
		ParseJson(resp, &m)

		followers, ok := m["followers"].([]interface{})
		if !ok {
			log.Fatalln("Parse followers as []interface{} failed")
			break
		}
		for _, v := range followers {
			//log.Println(reflect.TypeOf(u))
			mv := v.(map[string]interface{})
			//log.Println(mv)
			//log.Println(reflect.TypeOf(mv["id"]))
			u, _ := ValueToString(&mv)
			ch <- UserInfo(*u)
		}

		pageCount := m["count"].(float64)
		log.Println(pageCount)
		if !ok {
			log.Fatalln("Parse page count failed.")
			break
		}

		s.PageNo++
		if s.PageNo > int(pageCount) {
			log.Println("Parse end of followers.")
			break
		}
		time.Sleep(time.Second)
	}
	return
}

func ValueToString(in *map[string]interface{}) (out *map[string]string, err error) {
	out = &map[string]string{}

	for k, v := range *in {
		switch vv := v.(type) {
		case int:
		case int64:
			(*out)[k] = strconv.Itoa(int(vv))
		case float32:
		case float64:
			(*out)[k] = strconv.FormatFloat(vv, 'f', -1, 64)
		case string:
			(*out)[k] = vv
		}
	}
	//log.Println(*out)
	return
}

func SaveUsers(ch chan UserInfo) (err error) {
	log.Println("Save Users")
	u := <-ch
	fmt.Println(u)
	return
}

func GetAllUsers(s *Seed) (cnt int, err error) {
	fmt.Println("Start GetAllUsers")

	sc, err := NewSeaClient()
	bExist, err := sc.UserExisted(s.ID)
	if !bExist {
		//sc.AddUser()
	}

	bSeed, err := sc.UserIsSeed(s.ID)
	if !bSeed {
		sc.AddSeed(s.ID)
	}

	ch := make(chan UserInfo)
	go FetchFollowers(s, ch)
	//go SaveUsers(ch)
	for i := 1; ; i++ {
		u, ok := <-ch
		if !ok {
			break
		}
		if i%100 == 0 {
			fmt.Println(i)
		}
		m := map[string]string(u)
		sc.AddUser(u["id"], &m)
	}
	return
}
