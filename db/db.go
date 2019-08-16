package db

import "github.com/boltdb/bolt"

const myDbFile = "my.db"

func save() {
	db, err := bolt.Open(myDbFile, 0600, nil)

	if err != nil {
		panic(err)
	}

	defer db.Close()

}
