package main

import (
	examples "github.com/SpaceTent/GoStructDB/examples"
)

func main() {

	examples.BasicInsert()
	examples.BasicInsertWithCounter()
	examples.ConcurrentInserts()
	examples.InsertDateDefault()
	examples.MissingColumns()
	examples.Lite3()

}
