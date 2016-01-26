package sail

import (
	"testing"
)

func TestGetClearSeed(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	t.Log("Start GetSeed", sc.GetSeed())
	sc.ClearSeeds()
	t.Log("After ClearSeeds", sc.GetSeed())
}

func TestUserSeed(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	id := "123test"
	t.Log("Start GetSeed", sc.GetSeed())
	sc.DeleteSeed(id)
	b, err := sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error(id, "is seed after delete seed")
	}

	sc.AddSeed(id)
	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error(id, "is not seed after add seed")
	}

	sc.ClearSeeds()
	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error(id, "is seed after clear seed")
	}

	sc.AddSeed(id)
	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error(id, "is not seed after add seed")
	}

	t.Log("After ClearSeed", sc.GetSeed())
}

func TestAddDelChkUser(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	id := "123456test"
	err = sc.DeleteSeed(id)
	if err != nil {
		t.Log(err)
	}

	b, err := sc.UserExisted(id)
	if err != nil {
		t.Error(err)
	}
	if b {
		t.Error(id, "exists after delete")
	} else {
		t.Log(id, "delete successfully")
	}

	m := map[string]string{
		"count": "233",
		"name":  "2233",
		"home":  "bilibili",
	}
	err = sc.AddUser(id, &m)
	if err != nil {
		t.Error(err)
	}

	b, err = sc.UserExisted(id)
	if err != nil {
		t.Error(err)
	}
	if !b {
		t.Error(id, "do not exist after add")
	} else {
		t.Log(id, "add successfully")
	}

	mm, err := sc.GetUser(id)
	if err != nil {
		t.Error(err)
	}
	for k, v := range m {
		if v != (*mm)[k] {
			t.Error(k, v, "is not added or getted")
		} else {
			t.Log("Get", k, v, "in success")
		}
	}

	//	err = sc.DeleteUser(id)
	//	if err != nil {
	//		t.Error(err)
	//	}

	b, err = sc.UserExisted(id)
	if err != nil {
		t.Error(err)
	}
	if b {
		t.Error(id, "not deleted")
	} else {
		t.Log(id, "delete in success")
	}

	_, err = sc.GetUser(id)
	if err != nil {
		t.Error(err)
	}

}
