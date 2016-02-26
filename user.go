package User

import (
	"strconv"

	"github.com/ChengangDev/User/sail"
)
import "github.com/ChengangDev/User/sea"
import "log"

func SaveOneFollower(db sea.DbOp, sop sea.SeedOp, u *sail.UserInfo, overide bool) (err error) {
	//m := map[string]string(*u)

	//add to sorted map
	count, err := strconv.Atoi((*u)["followers_count"])
	if err != nil {
		return err
	}
	sop.AddSeed((*u)["id"], int64(count))
	return nil
}

//get and save followers of user
func FetchAndSaveFollowers(s *sail.Seed, db sea.DbOp, fin chan []int) (cnt int, err error) {
	log.Println("GetAndSaveFollowers of:", s.ID)

	ch := make(chan sail.UserInfo)
	go sail.FetchFollowers(s, ch)

	nAdd, nScan := 0, 0
	for {
		u, ok := <-ch
		if !ok {
			break
		}
		m := map[string]string(u)

		nScan++
		//skip added user
		b, err := db.UserExisted(u["id"])
		if err != nil {
			log.Fatalln(err)
			continue
		}
		if !b {
			db.AddUser(u["id"], &m)
			nAdd++
		}

		//more than 10000
		if nScan%10000 == 0 {
			log.Println(s.ID, ":", nAdd, "/", nScan)
		}
	}

	log.Println("GetAndSaveFollowers of", s.ID, "finished.", nAdd, "/", nScan)
	fin <- []int{nAdd, nScan}
	return
}

func GetFollowers(s *sail.Seed) (count int, err error) {

	return
}

func GetAllUsers(s *sail.Seed, db sea.DbOp) {
	sc, err := sea.NewSeaClient()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sc.Close()

	ch := make(chan []int)
	GetAndSaveFollowers(s, sc, ch)

	ad := <-ch
	log.Println(ad)
}
