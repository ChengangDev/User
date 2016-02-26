package sea

import (
	"testing"
)

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

func TestSeedOp(t *testing.T) {
	sc, err := NewSeaClient()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.ClearAllSeeds()
	defer sc.Close()

	err = sc.AddPreparation("1", false)
	if err == nil {
		t.Error("Need err here.")
	} else {
		t.Log("OK,", err.Error())
	}

	ids := []string{"1", "2", "3", "4"}
	for _, id := range ids {
		err = sc.AddPreparation(id, false, id)
		if err != nil {
			t.Error("AddPreparation Error:", err.Error())
		} else {
			t.Log("AddPreparation OK:", id)
		}

		bSeed := sc.UserIsSeed(id)
		if bSeed {
			t.Error("UserIsSeed Error:", id, "is not seed")
		} else {
			t.Log("UserIsSeed OK,", id)
		}
	}

	for _, id := range ids {
		err := sc.MarkSeed(id)
		if err != nil {
			t.Error("MarkSeed Error:", id, err.Error())
		} else {
			bSeed := sc.UserIsSeed(id)
			if bSeed {
				t.Log("MarkSeed OK:", id, "is seed")
			} else {
				t.Error("MarkSeed Error:", id, "is not seed")
			}
		}

		err = sc.UnmarkSeed(id)
		if err != nil {
			t.Error("UnmarkSeed Error:", id, err.Error())
		} else {
			bSeed := sc.UserIsSeed(id)
			if !bSeed {
				t.Log("UnmarkSeed OK:", id, "is not seed")
			} else {
				t.Error("UnmarkSeed Error:", id, "is  seed")
			}
		}
	}

}
