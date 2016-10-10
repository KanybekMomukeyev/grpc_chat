package database
//package main
//
//import (
//	"log"
//	"fmt"
//	"github.com/boltdb/bolt"
//	"time"
//	"strconv"
//	"runtime"
//)
//
//var world = []byte("world")
//
//func main() {
//
//	db, err := bolt.Open("bolt.db", 0644, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	t := []string{"g", "h", "i"}
//	fmt.Print(t)
//
//	key := []byte("hello")
//	value := []byte("value is Koke")
//
//	// store some data
//	err = db.Update(func(tx *bolt.Tx) error {
//		bucket, err := tx.CreateBucketIfNotExists(world)
//		if err != nil {
//			return err
//		}
//
//		err = bucket.Put(key, value)
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// retrieve the data
//	err = db.View(func(tx *bolt.Tx) error {
//		bucket := tx.Bucket(world)
//		if bucket == nil {
//			return fmt.Errorf("Bucket %q not found!", world)
//		}
//
//		val := bucket.Get(key)
//		fmt.Println(string(val))
//
//		return nil
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//
//const BUCKET = "widgets"
//
//func someMethod() {
//	runtime.GOMAXPROCS(runtime.NumCPU())
//	db, err := bolt.Open("testbolt.db", 0666, nil)
//	defer db.Close()
//	if err != nil{
//		fmt.Println(err)
//		return
//	}
//
//	tx, err := db.Begin(true)
//	if err != nil{
//		fmt.Println(err)
//		return
//	}
//
//	commitSize := 100000
//	infoSize := 100000
//
//	b, err := tx.CreateBucket([]byte(BUCKET))
//	if err != nil{
//		fmt.Println(err)
//		return
//	}
//	var startTime = time.Now()
//
//	for i:=0; i<10000; i++{
//		key := "foo" + strconv.Itoa(i) + "M"
//		b.Put([]byte(key), []byte("bar"))
//		key = "baz" + strconv.Itoa(i) + "Z"
//		b.Put([]byte(key), []byte("bat"))
//		//b.Delete([]byte("foo"))
//		if i%commitSize== 0 && i != 0{
//			err = tx.Commit()
//			if err != nil{
//				fmt.Println(err)
//				return
//			}
//			tx, err = db.Begin(true)
//			b = tx.Bucket([]byte(BUCKET))
//			if err != nil{
//				fmt.Println(err)
//				return
//			}
//		}
//		if i%infoSize == 0 && i != 0{
//			fmt.Println(i, time.Since(startTime))
//			startTime = time.Now()
//		}
//	}
//
//
//
//}