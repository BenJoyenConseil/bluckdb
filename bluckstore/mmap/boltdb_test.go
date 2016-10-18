package mmap

import (
	"github.com/boltdb/bolt"
	"testing"
	"strconv"
	"log"
	"fmt"
	"os"
	"math/rand"
)

func BenchmarkBoltDBPut(b *testing.B) {
	tableName := "table_test"

	os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte(tableName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.Update(func(tx *bolt.Tx) error {
			table := tx.Bucket([]byte("table_test"))
			err := table.Put([]byte(strconv.Itoa(n)), []byte("mec, elle est où ma caisse ??"))
			return err
		})
	}

}



func setupBolt() {
	tableName := "table_test"

	os.Remove("my.db")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var table *bolt.Bucket
	db.Update(func(tx *bolt.Tx) error {
		table, err = tx.CreateBucket([]byte(tableName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	size := 300
	for i := 0; i < size; i++ {
		db.Update(func(tx *bolt.Tx) error {
			table := tx.Bucket([]byte("table_test"))
			err := table.Put([]byte(strconv.Itoa(i)), []byte("mec, elle est où ma caisse ??"))
			return err
		})
	}
	db.Sync()
}

func BenchmarkBoltDBGet(b *testing.B) {
	setupBolt()
	db, _ := bolt.Open("my.db", 0600, nil)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.View(func(tx *bolt.Tx) error {
			table := tx.Bucket([]byte("table_test"))
			table.Get([]byte(strconv.Itoa(rand.Intn(30))))
			return nil
		})
	}
}