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

	t.Log("Start GetSeed", sc.getFlagSeedValue())
	sc.ClearSeeds()
	t.Log("After ClearSeeds", sc.getFlagSeedValue())
}

func TestUserSeed(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	id := "123test"
	t.Log("Start GetFlagSeedValue", sc.getFlagSeedValue())

	err = sc.DeleteSeed(id)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("DeleteSeed in success.")
	}

	b, err := sc.UserIsSeed(id)
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

	t.Log("After ClearSeed", sc.getFlagSeedValue())
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

func TestSort(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	sc.clearSortedUsers()
	rank, err := sc.getUserRank("2")
	if err != nil {
		t.Log(err)
		t.Log("ClearSortedUsers is OK")
	} else {
		t.Error("There should have error, but gets rank", rank)
	}

	score, err := sc.getUserScore("2")
	if err != nil {
		t.Log(err)
		t.Log("ClearSortedUsers is OK")
	} else {
		t.Error("There should have error, but gets score", score)
	}

	sc.sortUser("0", 0)
	sc.sortUser("1", 1)

	rank, err = sc.getUserRank("1")
	if err != nil {
		t.Error(err)
	} else {
		if rank == 1 {
			t.Log("SortUser 1 OK")
		} else {
			t.Error("SortUser is failed:", rank)
		}
	}

	rank, err = sc.getUserRank("2")
	if err != nil {
		t.Log(err)
	} else {
		t.Error("There should have error")
	}

	sc.sortUser("3", 3)
	rank, err = sc.getUserRank("3")
	if err != nil {
		t.Error(err)
	} else {
		if rank != 2 {
			t.Error("3 needs rank 2, but gets", rank)
		}
	}

	sc.sortUser("2", 2)
	rank, err = sc.getUserRank("2")
	if err != nil {
		t.Error(err)
	} else {
		if rank != 2 {
			t.Error("2 needs rank 2, but gets", rank)
		}
	}

	score, err = sc.getUserScore("2")
	if err != nil {
		t.Error(err)
	} else {
		if score != 2 {
			t.Error("2 needs score 2, but gets", score)
		}
	}

	rank, err = sc.getUserRank("3")
	if err != nil {
		t.Error(err)
	} else {
		if rank != 3 {
			t.Error("3 needs rank 3, but gets", rank)
		}
	}

	sc.sortUser("2", 4)
	rank, err = sc.getUserRank("2")
	if err != nil {
		t.Error(err)
	} else {
		if rank != 3 {
			t.Error("2 needs rank 3, but gets", rank)
		}
	}

	sc.deleteSortedUser("3")
	rank, err = sc.getUserRank("3")
	if err != nil {
		t.Log(err)
		t.Log("Delete 3 OK")
	} else {
		t.Log("DeleteSortedUser failed")
	}

	rank, err = sc.getUserRank("2")
	if err != nil {
		t.Error(err)
	} else {
		if rank != 2 {
			t.Error("2 needs rank 2, but gets", rank)
		}
	}

	score, err = sc.getUserScore("2")
	if err != nil {
		t.Error(err)
	} else {
		if score != 4 {
			t.Error("2 needs score 4, but gets", score)
		}
	}

	sc.clearSortedUsers()
	rank, err = sc.getUserRank("2")
	if err != nil {
		t.Log(err)
		t.Log("ClearSortedUsers ok")
	} else {
		t.Error("There should have error")
	}
}

func TestGenerateNewSeed(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	sc.clearSortedUsers()
	sc.sortUser("1", 1)
	sc.sortUser("2", 2)
	sc.sortUser("3", 3)
	sc.sortUser("4", 4)

	id, err := sc.GenerateNewSeed()
	if err != nil {
		t.Error(err)
	} else if id != "1" {
		t.Error("GenerateNewSeed needs 1, but gets", id)
	}

	sc.deleteSortedUser("1")
	id, err = sc.GenerateNewSeed()
	if err != nil {
		t.Error(err)
	} else if id != "2" {
		t.Error("GenerateNewSeed needs 2, but gets", id)
	}

	sc.deleteSortedUser("2")
	sc.deleteSortedUser("3")
	sc.deleteSortedUser("4")
	id, err = sc.GenerateNewSeed()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("Shoud get error, but gets", id)
	}

	sc.clearSortedUsers()
}
