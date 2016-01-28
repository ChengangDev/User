package sail

import "testing"

func TestGetAllUsers(t *testing.T) {
	//GetAllUsers(&TestSeed)
	sc, err := NewSeaClient()
	if err != nil {
		t.Error(err)
		return
	}
	defer sc.Close()

	//GetAndSaveFollowers(&TestSeed, sc)
}
