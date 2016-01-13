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

func TestGetRequest(t *testing.T) {
	cases := []struct {
		inUrl string
		want  string
	}{
		{"https://www.baidu.com", ""},
	}

	for _, c := range cases {
		ret, err := GetRequest(c.inUrl, nil)
		if err == nil {
			t.Log(err.Error())
		}
		t.Log(ret)
	}

}
