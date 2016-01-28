package User

import "github.com/ChengangDev/User/sail"
import "github.com/ChengangDev/User/sea"
import "log"

//
func GetAndSaveFollowers(s *sail.Seed, db sea.DbOp) (cnt int, err error) {
	log.Println("GetAndSaveFollowers of:", s.ID)

	ch := make(chan sail.UserInfo)
	go sail.FetchFollowers(s, ch)

	i := 0
	for {
		u, ok := <-ch
		if !ok {
			break
		}
		m := map[string]string(u)
		db.AddUser(u["id"], &m)
		i++

		//more than 10000
		if i%10000 == 0 {
			log.Println(s.ID, ":", i)
		}
	}

	db.AddSeed(s.ID)
	log.Println(s.ID, "has", i, "followers.")
	return
}
