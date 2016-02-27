package User

import "testing"

import "github.com/ChengangDev/User/sea"
import "github.com/ChengangDev/User/sail"

var TestSeed = sail.Seed{
	FixedFormater: "http://xueqiu.com/friendships/followers.json?uid=%v&pageNo=%v&size=%v",
	ID:            "3037882447",
	PageNo:        1,
	PageSize:      1000,
	Interval:      1,
}

func TestFetchFollowers(t *testing.T) {

	//	host_url := "http://xueqiu.com"
	//	ids := []string{
	//		"6346418304",
	//		"9905072371",
	//		"4484481018",
	//		"1855822841",
	//	}
	//	sck := sail.NewSharedCookie(host_url)
	//	seed := &sail.Seed{
	//		FixedFormater: "http://xueqiu.com/friendships/followers.json?uid=%v&pageNo=%v&size=%v",
	//		ID:            "1234461197",
	//		PageNo:        1,
	//		PageSize:      1000,
	//		Interval:      1,
	//	}

	//	seed.ID = ids[0]
	//	ch1 := make(chan map[string]string)
	//go FetchFollowers(seed, sck, ch1)
	//	for u := range ch1 {

	//		t.Log(u)

	//	}
}

func TestGetAllUsers(t *testing.T) {

	//GetAllUsers(&TestSeed)
}

func TestManager(t *testing.T) {
	host_url := "http://xueqiu.com"
	sck := sail.NewSharedCookie(host_url)
	seed := &sail.Seed{
		FixedFormater: "http://xueqiu.com/friendships/followers.json?uid=%v&pageNo=%v&size=%v",
		ID:            "1234461197",
		PageNo:        1,
		PageSize:      1000,
		Interval:      1,
	}
	sc, err := sea.NewSeaClient()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer sc.Close()

	Manager(seed, sck, sc, sc, 4)
}
