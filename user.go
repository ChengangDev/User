package User

import "github.com/ChengangDev/User/sail"
import "github.com/ChengangDev/User/sea"
import "log"

//get and save followers of user
func GetAndSaveFollowers(s *sail.Seed, db sea.DbOp, fin chan []int) (cnt int, err error) {
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

func GetAllUsers(s *sail.Seed, db sea.DbOp) {

}
