package db

import (
	"cryptom/model"
	"github.com/boltdb/bolt"
)

const myDbFile = "my.db"

func save() {
	db, err := bolt.Open(myDbFile, 0600, nil)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err := db.Update(update)


}

func update(tx *bolt.Tx) error {
	b := tx.Bucket([]byte(myDbFile))

	if b == nil {
		genesis := model.NewGenesis()
		b, err := tx.CreateBucket([]byte(myDbFile))

		if err != nil {
			panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		err = b.Put([]byte("l"), genesis.Hash)
		tip = genesis.Hash
	} else {
		tip = b.Get([]byte("l"))
	}

	return nil
}




