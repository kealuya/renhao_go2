package badger

import (
	badger "github.com/dgraph-io/badger/v4"
	"log"
	"time"
)

// go install github.com/dgraph-io/badger/v4/badger@latest
func Badger() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("./badger/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	const Processed = 0b00000001 // 00000001
	// Start a writable transaction.
	err = db.Update(func(txn *badger.Txn) error {
		// Create a new entry with a key, value, user meta, and TTL.
		entry := badger.NewEntry([]byte("answer1"), []byte("11123")).
			WithMeta(Processed).      // Set user meta to 1.
			WithTTL(20 * time.Second) // Set TTL to 1 hour.

		// Set the entry in the transaction.
		err := txn.SetEntry(entry)
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	// Start a read-only transaction.
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer2"))
		if err != nil {
			return err
		}

		// Fetch the user meta from the item.
		userMeta := item.UserMeta()
		log.Printf("User Meta: %d", userMeta) // Should print "User Meta: 1"

		// Fetch the value from the item.
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		log.Printf("Value: %s", val) // Should print "Value: 42"
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
