package main

import (
	"os"
)

func main() {
	// initialize db
	dal, _ := newDal("db.db", os.Getpagesize())

	// create new page
	page := dal.allocateEmptyPage()
	page.num = dal.freelist.getNextPage()
	copy(page.data[:], "data")

	// commit it
	_ = dal.writePage(page)
}
