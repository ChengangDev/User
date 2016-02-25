package sea

import (
	"log"
	//	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type DbOp interface {

	//clear all users
	//ClearUsers() error

	//check if user is existed
	UserExisted(id string) (bool, error)
	//add a user
	AddUser(id string, v *map[string]string) error
	//get user all info
	GetUserAll(id string) (v *map[string]string, err error)
	//get user specific info
	//GetUserInfo(id string) (info string, err error)
	//delete a user
	DeleteUser(id string) error
}

type SeaClient struct {
	p   *pool.Pool
	cli *redis.Client
}

var NAMESPACE_XUEQIU_USER = "com.xueqiu:user:"

func getID(id string) string {
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
func (sc *SeaClient) UserExisted(id string) (b bool, err error) {
	ret, err := sc.cli.Cmd("EXISTS", getID(id)).Int()
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
	ar := []string{getID(id)}
	for kk, vv := range *v {
		ar = append(ar, kk, vv)
	}
	sc.cli.Cmd("HMSET", ar)
	return
}

func (sc *SeaClient) GetUserAll(id string) (v *map[string]string, err error) {
	m, err := sc.cli.Cmd("HGETALL", getID(id)).Map()
	return &m, err
}

func (sc *SeaClient) DeleteUser(id string) (err error) {
	sc.cli.Cmd("DEL", getID(id))
	return
}
