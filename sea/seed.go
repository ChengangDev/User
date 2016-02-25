package sea

import (
	"errors"
	"log"
)

type SeedOp interface {
	//clear all seeds
	ClearSeeds()
	//check if user is a seed
	UserIsSeed(id string) (bool, error)
	//mark user as a seed. if user does not exist, add it
	AddSeed(id string) error
	//unmark user
	DeleteSeed(id string) error
	//generate new seed
	GenerateNewSeed() (id string, err error)
}

var NAMESPACE_XUEQIU_SEED = "com.xueqiu:seed:"
var NAMESPACE_XUEQIU_FOLLOWER = "com.xueqiu:follower:"
var FLAG_SEED = "FLAG_SEED"
var VALUE_INIT_SEED = int64(0)

func getSeedID(id string) string {
	return NAMESPACE_XUEQIU_SEED + id
}

func getZName() string {
	return NAMESPACE_XUEQIU_FOLLOWER + "count"
}

//get current seed flag
func (sc *SeaClient) getFlagSeedValue() int64 {
	seed, err := sc.cli.Cmd("GET", getSeedID(FLAG_SEED)).Int64()
	if err != nil {
		log.Println(err)
		log.Println("GetSeed failed, use VALUE_INIT_SEED")
		return VALUE_INIT_SEED
	}

	if seed < 0 {
		log.Println("GetSeed smaller than 0, ", seed)
		return VALUE_INIT_SEED
	}

	return seed
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

//start from 0
func (sc *SeaClient) getUserRank(id string) (rank int64, err error) {
	rank, err = sc.cli.Cmd("ZRANK", getZName(), id).Int64()
	if err != nil {
		return -1, err
	}
	return rank, nil
}

//get score of memeber
func (sc *SeaClient) getUserScore(id string) (rank int64, err error) {
	rank, err = sc.cli.Cmd("ZSCORE", getZName(), id).Int64()
	if err != nil {
		return -1, err
	}
	return rank, nil
}

//make all old seeds be not seeds by increasing the seed flag by one
//since there are too many old seeds
func (sc *SeaClient) ClearSeeds() {
	seed := sc.getFlagSeedValue()
	seed++
	sc.cli.Cmd("SET", getSeedID(FLAG_SEED), seed)
	return
}

//used to be seed or not
func (sc *SeaClient) UserIsSeed(id string) (b bool, err error) {
	v, err := sc.cli.Cmd("GET", getSeedID(id)).Int64()
	if err != nil {
		return false, err
	}

	if v == sc.getFlagSeedValue() {
		return true, nil
	}
	return false, nil
}

func (sc *SeaClient) AddSeed(id string) (err error) {
	err = sc.cli.Cmd("SET", getSeedID(id), sc.getFlagSeedValue()).Err
	return
}

func (sc *SeaClient) DeleteSeed(id string) (err error) {
	ok, err := sc.cli.Cmd("SET", getSeedID(id), -1).Str()
	if err != nil {
		return err
	}

	if ok != "OK" {
		return errors.New("Failed to set KEY_SEED -1")
	}
	return nil
}

//generate seed from users by scan hmap skip users which are seed
//already or have less than
func (sc *SeaClient) GenerateNewSeed() (id string, err error) {
	count, err := sc.cli.Cmd("ZCARD", getZName()).Int64()
	if err != nil {
		return "", err
	} else if count == 0 {
		return "", errors.New("No more seeds.")
	}

	ids, err := sc.cli.Cmd("ZRANGE", getZName(), 0, 1).List()
	if err != nil {
		return "", err
	}
	id = ids[0]
	return
}
