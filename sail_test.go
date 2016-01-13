package sail

import (
	//	"net/http"
	"testing"
)

func TestPlay(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello", "Hello"},
		{"Thank", "Thank"},
	}

	for _, c := range cases {
		got := Play(c.in)
		if got == c.want {
			t.Logf("Success: %s", c.in)
		} else {
			t.Errorf("Failure: in:%s got:%s want:%s", c.in, got, c.want)
		}
	}
}

var XueQiuRudder = Rudder{
	CountPattern:  "\\\"count\\\":[0-9]*",
	PageNoPattern: "\\\"page\\\":[0-9]*",
	//PageSizePattern:  "",
	PageCountPattern: "\\\"maxPage\\\":[0-9]*",

	IDPattern:    "\\\"id\\\":[0-9]*",
	NamePattern:  "\\\"screen_name\\\":\\\"*\\\"",
	OtherPattern: map[string]string{},
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
