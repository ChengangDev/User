package sea

import (
	"testing"
)

func TestGetClearSeed(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	t.Log("Start GetSeed", sc.GetSeedValue())
	sc.ClearSeeds()
	t.Log("After ClearSeeds", sc.GetSeedValue())
}

func TestUserSeed(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	id := "123test"
	t.Log("Start GetSeed", sc.GetSeedValue())
	b, err := sc.UserExisted(id)
	if err != nil {
		t.Error(err)
	} else {
		if !b {
			sc.AddUser(id, &map[string]string{"test": "test"})
			t.Log("Add User firstly:")
		} else {
			t.Log("User exists:")
		}
	}
	v, err := sc.GetUserAll(id)
	t.Log("OldUser:", *v)

	err = sc.DeleteSeed(id)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("DeleteSeed in success.")
	}
	v, err = sc.GetUserAll(id)
	t.Log("After DeleteSeed", *v)

	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Error(err)
	} else {
		if b {
			t.Error(id, "is seed after delete seed")
		} else {
			t.Log("Check UserIsSeed in success.")
		}
	}

	err = sc.AddSeed(id)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("AddSeed no error")
	}
	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error(id, "is not seed after add seed")
	} else {
		t.Log("UserIsSeed after AddSeed")
	}

	sc.ClearSeeds()
	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error(id, "is seed after clear seed")
	} else {
		t.Log("User is not Seed after clearing")
	}

	sc.AddSeed(id)
	b, err = sc.UserIsSeed(id)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error(id, "is not seed after add seed")
	}

	t.Log("After ClearSeed", sc.GetSeedValue())
	sc.DeleteUser(id)
}

func TestAddDelChkUser(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	id := "123test"
	err = sc.DeleteUser(id)
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

	mm, err := sc.GetUserAll(id)
	if err != nil {
		t.Error(err)
	}
	t.Log("RawMap:", *mm)
	for k, v := range m {
		if v != (*mm)[k] {
			t.Error(k, v, "is not added or getted")
		} else {
			t.Log("Get", k, v, "in success")
		}
	}

	err = sc.DeleteUser(id)
	if err != nil {
		t.Error(err)
	}

	b, err = sc.UserExisted(id)
	if err != nil {
		t.Error(err)
	}
	if b {
		t.Error(id, "not deleted")
	} else {
		t.Log(id, "delete in success")
	}

	_, err = sc.GetUserAll(id)
	if err != nil {
		t.Error(err)
	}

}
