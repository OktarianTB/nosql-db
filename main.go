package main

func main() {
	// initialize db
	dal, _ := newDal("db.db")

	// create new page
	page := dal.allocateEmptyPage()
	page.num = dal.getNextPage()
	copy(page.data[:], "data")

	// commit it
	_ = dal.writePage(page)
	_, _ = dal.writeFreelist()

	// close the db
	_ = dal.close()

	// open db again
	// we expect the freelist state to have been saved
	dal, _ = newDal("db.db")
	page = dal.allocateEmptyPage()
	page.num = dal.getNextPage()
	copy(page.data[:], "data2")
	_ = dal.writePage(page)

	// Create a page and free it so the released pages will be updated
	pageNum := dal.getNextPage()
	dal.releasePage(pageNum)

	// commit it
	_, _ = dal.writeFreelist()
}
