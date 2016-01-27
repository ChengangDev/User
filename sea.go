package sail

import (
	"errors"
	"log"
	//	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type DbOp interface {
	//clear all seeds
	ClearSeeds()
	//get seed current flag
	GetSeedValue() int64
	//check if user is a seed
	UserIsSeed(id string) (bool, error)
	//mark user as a seed. if user does not exist, add it
	AddSeed(id string) error
	//unmark user
	DeleteSeed(id string) error
	//generate new seed
	GenerateNewSeed() (id string, err error)

	//clear all users
	//ClearUsers() error

	//check if user is existed
	UserExisted(id string) (bool, error)
	//add a user
	AddUser(id string, v *map[string]string) error
	//get user all info
	GetUserAll(id string) (v *map[string]string, err error)
	//get user specific info
	GetUserInfo(id string) (info string, err error)
	//delete a user
	DeleteUser(id string) error
}

type SeaClient struct {
	p   *pool.Pool
	cli *redis.Client
}

var NAMESPACE_XUEQIU_USER = "com.xueqiu:user:"
var NAMESPACE_XUEQIU_SEED = "com.xueqiu:seed:"
var KEY_SEED = "KEY_SEED"
var VALUE_INIT_SEED = int64(0)

func GetID(id string) string {
	return NAMESPACE_XUEQIU_USER + id
}

//methods

//create a new client
func NewSeaClient() (sea *SeaClient, err error) {
	sea = &SeaClient{}
	sea.p, err = pool.New("tcp", "localhost:6379", 8)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	sea.cli, err = sea.p.Get()
	if err != nil {
		//log.
		return nil, err
	}

	return sea, nil
}

//close the connection
func (sc *SeaClient) Close() {
	sc.p.Put(sc.cli)
}

//interfaces

//get current seed flag
func (sc *SeaClient) GetSeedValue() int64 {
	seed, err := sc.cli.Cmd("GET", NAMESPACE_XUEQIU_SEED+KEY_SEED).Int64()
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

//make all not seed by increase seed flag by one
//since there are too many old seeds
func (sc *SeaClient) ClearSeeds() {
	seed := sc.GetSeedValue()
	seed++
	sc.cli.Cmd("SET", NAMESPACE_XUEQIU_SEED+KEY_SEED, seed)
	return
}

//used to be seed or not
func (sc *SeaClient) UserIsSeed(id string) (b bool, err error) {
	v, err := sc.cli.Cmd("HGET", GetID(id), KEY_SEED).Int64()
	if err != nil {
		return false, err
	}

	if v == sc.GetSeedValue() {
		return true, nil
	}
	return false, nil
}

func (sc *SeaClient) AddSeed(id string) (err error) {
	err = sc.cli.Cmd("HMSET", GetID(id), KEY_SEED, sc.GetSeedValue()).Err
	return
}

func (sc *SeaClient) DeleteSeed(id string) (err error) {
	ok, err := sc.cli.Cmd("HMSET", GetID(id), KEY_SEED, -1).Str()
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

	return
}

func (sc *SeaClient) UserExisted(id string) (b bool, err error) {
	ret, err := sc.cli.Cmd("EXISTS", GetID(id)).Int()
	if err != nil {
		log.Fatalln(err)
		return
	}

	if ret == 1 {
		return true, nil
	}
	return false, nil
}

func (sc *SeaClient) AddUser(id string, v *map[string]string) (err error) {
	ar := []string{GetID(id)}
	for kk, vv := range *v {
		ar = append(ar, kk, vv)
	}
	sc.cli.Cmd("HMSET", ar)
	return
}

func (sc *SeaClient) GetUserAll(id string) (v *map[string]string, err error) {
	m, err := sc.cli.Cmd("HGETALL", GetID(id)).Map()
	return &m, err
}

func (sc *SeaClient) DeleteUser(id string) (err error) {
	sc.cli.Cmd("DEL", GetID(id))
	return
}
