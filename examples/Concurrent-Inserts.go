package examples

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	MySQL "github.com/SpaceTent/GoStructDB/gsdb"
)

func ConcurrentInserts() {

	godotenv.Load(".env")
	DSN := os.Getenv("DSN")
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	l := slog.New(textHandler)

	MySQL.New(DSN, l, context.Background())

	type InsertPerson struct {
		Id      int       `db:"column=id primarykey=yes table=Test"`
		Name    string    `db:"column=name"`
		Dtadded time.Time `db:"column=dtadded"`
		Status  int       `db:"column=status"`
		Ignored int       `db:"column=ignored omit=yes"`
	}

	var wg sync.WaitGroup
	numInserts := 1000

	for i := 1; i <= numInserts; i++ {
		wg.Add(1)

		// Launch a goroutine for each insert operation
		go func(insertId int) {
			defer wg.Done() // Decrement the counter when the goroutine completes

			entry := InsertPerson{
				Name:    fmt.Sprintf("Test %d", insertId),
				Dtadded: time.Now(),
				Status:  1,
			}

			// Construct SQL for the insert
			sqlQuery, err := MySQL.DB.Insert(entry)
			if err != nil {
				l.Error("Insert error: " + err.Error())
				return
			}

			// Execute the insert query
			insertedID, _, err := MySQL.DB.Execute(sqlQuery)
			if err != nil {
				l.Error("Execution error: " + err.Error())
			}

			test, err2 := MySQL.QuerySingleStruct[InsertPerson]("SELECT * FROM test WHERE id = ?", insertedID)
			if err2 != nil {
				l.Error("Query error: " + err2.Error())
			} else {
				if test.Id != int(insertedID) {
					l.Error("ID mismatch")
				}
			}

			_, _, err3 := MySQL.DB.Execute("DELETE FROM test WHERE id = ?", insertedID)
			if err3 != nil {
				l.Error("Query error: " + err3.Error())
			}

		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	l.Info("All inserts completed")
}
