package sea

import (
	"errors"
	"fmt"
	"log"
	//	"log"
	"strconv"
)

type SeedOp interface {
	//add id as preparation
	AddPreparation(id string, override bool, more ...string) error
	//delete all seed data including preparation id
	ClearAllSeeds()

	//check if user is a seed
	UserIsSeed(id string) (bool, error)
	//mark user as a seed. if user does not exist, add it
	MarkSeed(id string) error
	//unmark user
	UnmarkSeed(id string) error
	//unmark all seeds
	//UnmarkAllSeeds()
	//pick a user as seed, user will be mark as seed
	PickSeed() (id string, err error)
}

var NAMESPACE_XUEQIU_FOLLOWER = "com.xueqiu:follower:"
var OFFSET = int64(400)
var MIN_COUNT = int64(10000)

func getZName() string {
	return NAMESPACE_XUEQIU_FOLLOWER + "count"
}

func (sc *SeaClient) clearSortedUsers() {
	sc.cli.Cmd("DEL", getZName())
}

func (sc *SeaClient) sortUser(id string, cnt int64) {
	sc.cli.Cmd("ZADD", getZName(), cnt, id)
}

func (sc *SeaClient) deleteSortedUser(id string) {
	sc.cli.Cmd("ZREM", getZName(), id)
}

func (sc *SeaClient) userInSortedSet(id string) bool {
	resp := sc.cli.Cmd("ZSCORE", getZName(), id)
	if resp.String() == "Resp(Nil)" {
		return false
	} else {
		return true
	}
}

//start from 0
func (sc *SeaClient) getUserRank(id string) (rank int64, err error) {
	rank, err = sc.cli.Cmd("ZRANK", getZName(), id).Int64()
	if err != nil {
		return -1, err
	}
	return rank, nil
}

//get score of memeber
func (sc *SeaClient) getUserScore(id string) (score int64, err error) {
	resp := sc.cli.Cmd("ZSCORE", getZName(), id)
	if resp.String() == "Resp(Nil)" {
		return 0, errors.New("User not exists.")
	}

	score, err = resp.Int64()
	if err != nil {
		return -1, err
	}
	return score, nil
}

func (sc *SeaClient) AddPreparation(id string, override bool, more ...string) error {
	if len(more) != 1 {
		return errors.New(fmt.Sprint("AddPreparation needs 3 arguments, but get", 2+len(more)))
	}
	if !override {
		bSeed := sc.UserIsSeed(id)
		if bSeed {
			//skip
			return nil
		}
	}

	count, err := strconv.Atoi(more[0])
	if err != nil {
		return err
	}
	sc.sortUser(id, int64(count))
	return nil
}

func (sc *SeaClient) ClearAllSeeds() {
	sc.clearSortedUsers()
	return
}

//used to be seed or not
func (sc *SeaClient) UserIsSeed(id string) (b bool) {
	if !sc.userInSortedSet(id) {
		return false
	}
	score, err := sc.getUserScore(id)
	if err != nil {
		log.Fatalln(err)
		return false
	}

	if score < 0 {
		return true
	}
	return false
}

func (sc *SeaClient) MarkSeed(id string) (err error) {
	count, err := sc.getUserScore(id)
	if err != nil {
		return err
	}
	if count > 0 {
		sc.sortUser(id, -count)
	}
	return nil
}

func (sc *SeaClient) UnmarkSeed(id string) (err error) {
	count, err := sc.getUserScore(id)
	if err != nil {
		return err
	}
	if count < 0 {
		sc.sortUser(id, -count)
	}
	return nil
}

//generate seed from users by scan hmap skip users which are seed
//already or have less than
func (sc *SeaClient) PickSeed() (id string, err error) {
	total, err := sc.cli.Cmd("ZCARD", getZName()).Int64()
	if err != nil {
		return "", err
	} else if total == 0 {
		return "", errors.New("No more seeds.")
	}

	ids, err := sc.cli.Cmd("ZRANGE", getZName(), 0, 1).List()
	if err != nil {
		return "", err
	}

	return ids[0], nil
}
