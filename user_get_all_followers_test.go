package User

import "testing"
import "github.com/ChengangDev/User/sea"
import "github.com/ChengangDev/User/sail"

var TestSeed = sail.Seed{
	FixedFormater: "http://xueqiu.com/friendships/followers.json?uid=%v&pageNo=%v&size=%v",
	ID:            "3037882447",
	PageNo:        1,
	PageSize:      1000,
	Interval:      100,
}

func TestGetAndSaveFollowers(t *testing.T) {

	sc, err := sea.NewSeaClient()
	if err != nil {
		t.Error(err)
		return
	}
	defer sc.Close()

	ch := make(chan []int)
	GetAndSaveFollowers(&TestSeed, sc, ch)

	ad := <-ch
	t.Log(ad)
}
