package User

import (
	//	"strconv"
	"time"

	"github.com/ChengangDev/User/sail"
)
import "github.com/ChengangDev/User/sea"
import "log"

//fetch all followers of the seed
func FetchFollowers(s *sail.Seed, sck *sail.SharedCookie, out chan map[string]string) (err error) {
	log.Println("FetchFollowers:Start Fetch Followers of", s.ID)
	defer close(out)

	for {
		//interval between each request
		sch := make(chan int)
		tik := func(sig chan int) {
			time.Sleep(time.Millisecond * time.Duration(s.Interval))
			close(sig)
		}
		tik(sch)

		url := s.GetUrl()
		var resp string
		if sck == nil {
			//get without cookie
			resp, err = sail.GetRequestByCookie(url, nil)
			if err != nil {
				log.Fatal(err)
				continue
			}
		} else {
			//use cookie
			cj, err := sck.GetSharedCookie()
			if err != nil {
				log.Fatal(err)
				return err
			}
			resp, err = sail.GetRequestByCookie(url, cj)
			if err != nil {
				//in case cookie out of date, update cookie for another time
				log.Println("Try update cookie...")
				cj, err := sck.UpdateSharedCookie()
				if err != nil {
					log.Fatal(err)
					return err
				}
				//another get using new cookie
				resp, err = sail.GetRequestByCookie(url, cj)
				if err != nil {
					log.Fatal(err)
					continue
				}
				log.Println("Use new cookie in success.")
			}
		}
		m := map[string]interface{}{}
		sail.ParseJson(resp, &m)

		followers, ok := m["followers"].([]interface{})
		if !ok {
			log.Fatalln("Parse followers as []interface{} failed")
			continue
		}
		for _, v := range followers {
			//log.Println(reflect.TypeOf(u))
			mv := v.(map[string]interface{})
			//log.Println(mv)
			//log.Println(reflect.TypeOf(mv["id"]))
			u, _ := sail.ValueToString(&mv)
			out <- *u
		}

		pageCount := m["maxPage"].(float64)
		//log.Println(pageCount)
		if !ok {
			log.Fatalln("Parse page count failed.")
			break
		}

		s.PageNo++
		if s.PageNo > int(pageCount) {
			log.Println("Parse end of followers.")
			break
		}

		//interval ends
		for _ = range sch {
			//
		}
	}

	return
}

func SaveOneFollower(db sea.DbOp, sop sea.SeedOp, u *map[string]string) (saved bool, err error) {
	//m := map[string]string(*u)

	//add to sorted map
	count := (*u)["followers_count"]
	if err != nil {
		log.Println("SaveOneFollower Error:", err)
	}
	err = sop.AddPreparation((*u)["id"], false, count)
	if err != nil {
		log.Println("SaveOneFollower Error:", err)
	}

	//add to user database
	if db.UserExisted((*u)["id"]) {
		return false, nil
	}
	err = db.AddUser((*u)["id"], u)
	if err != nil {
		return false, err
	}
	return true, nil
}

//get and save followers of user
func FetchAndSaveFollowers(s *sail.Seed, sck *sail.SharedCookie, db sea.DbOp, sop sea.SeedOp, fin chan int) (cnt int, err error) {
	log.Println("FetchAndSaveFollowers of:", s.ID)
	defer close(fin)

	ch := make(chan map[string]string)
	go FetchFollowers(s, sck, ch)

	nAdd, nScan := 0, 0
	for {
		u, ok := <-ch
		if !ok {
			break
		}

		nScan++
		//skip added user

		saved, err := SaveOneFollower(db, sop, &u)
		if err != nil {
			log.Println(err)
			continue
		} else {
			if saved {
				nAdd++
			}
		}

		//more than 10000
		if nScan%10000 == 0 {
			log.Println(s.ID, ":", nAdd, "/", nScan)
		}
	}

	log.Println("FetchAndSaveFollowers of", s.ID, "finished.", nAdd, "/", nScan)

	return
}

//fetch all followers of the seed
func FetchWorker(name int, in chan sail.Seed, sck *sail.SharedCookie, out chan map[string]string) (err error) {
	log.Println("FetchWorker", name, "Start.")
	defer close(out)
	for s := range in {
		log.Println("FetchWorker", name, ":Get New Seed", s.ID)
		for {
			//interval between each request
			sch := make(chan int)
			tik := func(sig chan int) {
				time.Sleep(time.Millisecond * time.Duration(s.Interval))
				close(sig)
			}
			tik(sch)

			url := s.GetUrl()
			var resp string
			if sck == nil {
				//get without cookie
				resp, err = sail.GetRequestByCookie(url, nil)
				if err != nil {
					log.Fatal(err)
					continue
				}
			} else {
				//use cookie
				cj, err := sck.GetSharedCookie()
				if err != nil {
					log.Fatal(err)
					return err
				}
				resp, err = sail.GetRequestByCookie(url, cj)
				if err != nil {
					//in case cookie out of date, update cookie for another time
					log.Println("FetchWorker", name, "Try update cookie...")
					cj, err := sck.UpdateSharedCookie()
					if err != nil {
						log.Fatal(err)
						return err
					}
					//another get using new cookie
					resp, err = sail.GetRequestByCookie(url, cj)
					if err != nil {
						log.Fatal(err)
						continue
					}
					log.Println("FetchWorker", name, "Use new cookie in success.")
				}
			}
			m := map[string]interface{}{}
			sail.ParseJson(resp, &m)

			followers, ok := m["followers"].([]interface{})
			if !ok {
				log.Fatalln("FetchWorker", name, "Parse followers as []interface{} failed")
				continue
			}
			for _, v := range followers {
				//log.Println(reflect.TypeOf(u))
				mv := v.(map[string]interface{})
				//log.Println(mv)
				//log.Println(reflect.TypeOf(mv["id"]))
				u, _ := sail.ValueToString(&mv)
				out <- *u
			}

			pageCount := m["maxPage"].(float64)
			//log.Println(pageCount)
			if !ok {
				log.Fatalln("FetchWorker", name, "Parse page count failed.")
				break
			}

			s.PageNo++
			if s.PageNo > int(pageCount) {
				log.Println("FetchWorker", name, "Parse end of followers.")
				break
			}

			//interval ends
			for _ = range sch {
				//
			}
		}

	}

	return
}

func SaveWorker(name int, in chan map[string]string, db sea.DbOp, sop sea.SeedOp) {
	log.Println("SaveWorker", name, "start.")
	nAdd, nScan := 0, 0
	for u := range in {
		nScan++
		//skip added user

		saved, err := SaveOneFollower(db, sop, &u)
		if err != nil {
			log.Println(err)
			continue
		} else {
			if saved {
				nAdd++
			}
		}

		//more than 10000
		if nScan%10000 == 0 {
			log.Println("SaveWorker", name, nAdd, "/", nScan)
		}
	}
	log.Println("SaveWorker", name, "finish", nAdd, "/", nScan)
}

func Manager(s *sail.Seed, sck *sail.SharedCookie, db sea.DbOp, sop sea.SeedOp, n int) {
	if n < 1 {
		log.Fatal("FetchManager Error: number of workers must > 0")
		return
	}

	in := make(chan sail.Seed)
	out := make(chan map[string]string)

	for i := 0; i < n; i++ {
		go FetchWorker(i, in, sck, out)

		if i == 0 {
			if s != nil {
				in <- *s
			}
		}
	}

	//only one save worker
	go SaveWorker(0, out, db, sop)

	pick := func(score interface{}) bool {
		if score.(int64) >= 1000 {
			return true
		} else {
			return false
		}
	}

	for {
		id, err := sop.GetSeed(pick)
		if err != nil {
			log.Println("GetSeed Error:", err)
			log.Println("Try GetSeed After 1 second.")
			time.Sleep(time.Second * 1)
			continue
		} else {
			if id == "" {
				log.Println("No Seeds. Try GetSeed After 1 second.")
				time.Sleep(time.Second * 1)
				continue
			}

			s.ID = id
			log.Println("GetSeed New:", s.ID)
			in <- *s
		}
	}
}

var DefSeed = sail.Seed{
	FixedFormater: "http://xueqiu.com/friendships/followers.json?uid=%v&pageNo=%v&size=%v",
	ID:            "1234461197",
	PageNo:        1,
	PageSize:      1000,
	Interval:      100,
}

var DefCookie = sail.NewSharedCookie("http://xueqiu.com")

func GetAllUsers(s *sail.Seed) {
	sc, err := sea.NewSeaClient()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sc.Close()

	fin := make(chan int)
	go FetchAndSaveFollowers(s, DefCookie, sc, sc, fin)

	for _ = range fin {

	}
	//log.Println(ad)
}
