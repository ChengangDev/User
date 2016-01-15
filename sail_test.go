package sail

import (
	"io/ioutil"
	"net/http"
	//	"net/http"
	"testing"
)

var TestSeed = Seed{
	FixedFormater: "http://xueqiu.com/friendships/followers.json?uid=%v&pageNo=%v&size=%v",
	ID:            "3037882447",
	PageNo:        1,
	PageSize:      2,
	Host:          "com.xueqiu",
	Depth:         6,
	Thread:        8,
	Interval:      100,
}

var XueQiuRudder = Rudder{

	CountPatterns: []string{
		"\\\"count\\\":[0-9]*",
		"[0-9]+"},
	PageNoPatterns: []string{
		"\\\"page\\\":[0-9]*",
		"[0-9]+"},

	PageCountPatterns: []string{
		"\\\"maxPage\\\":[0-9]*",
		"[0-9]+"},

	IDsPatterns: []string{
		"\\\"id\\\":[0-9]*",
		"[0-9]+"},
	NamesPatterns: []string{
		"\\\"screen_name\\\":\\\"*\\\"",
		""},
	OtherPatterns:     []string{},
	OtherListPatterns: []string{},
}

func TestGetRequest(t *testing.T) {
	cases := []struct {
		url    string
		header *map[string]string
	}{
		{"http://127.0.0.1:8000/index.html", &DefaultHeader},
		{"http://www.baidu.com", &DefaultHeader},
		{"http://xueqiu.com/friendships/followers.json?pageNo=1&uid=3037882447&size=20", &DefaultHeader},
	}

	for _, c := range cases {
		_, err := GetRequest(c.url, c.header)
		if err != nil {
			t.Log(err.Error())
		}
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		url    string
		header *map[string]string
		rudder *Rudder
	}{
		{"http://xueqiu.com/friendships/followers.json?pageNo=1&uid=3037882447&size=20",
			&DefaultHeader, &XueQiuRudder},
	}

	for _, c := range cases {
		resp, err := GetRequest(c.url, c.header)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		Parse(resp, c.rudder)
	}
}

func TestGetCookie(t *testing.T) {
	cj, err := GetCookie("http://xueqiu.com")

	cli := http.Client{Jar: cj}
	resp, err := cli.Get("http://xueqiu.com/friendships/followers.json?pageNo=1&uid=3037882447&size=20")
	s, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%q", s)
}

func TestUsers(t *testing.T) {
	//	_, err := Users(&TestSeed, &XueQiuRudder)
	//	if err != nil {
	//		t.Log(err.Error())
	//	}
}
