package main

import (
	"fmt"
	"os"
)

type pgnum uint64

type page struct {
	num  pgnum
	data []byte
}

type dal struct {
	file     *os.File
	pageSize int

	freelist *freelist
}

func newDal(path string, pageSize int) (*dal, error) {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	dal := &dal{
		file,
		pageSize,
		newFreelist(),
	}

	return dal, nil
}

func (dal *dal) close() error {
	if dal.file != nil {
		err := dal.file.Close()
		if err != nil {
			return fmt.Errorf("could not close file: %s", err)
		}
		dal.file = nil
	}

	return nil
}

func (dal *dal) allocateEmptyPage() *page {
	return &page{
		data: make([]byte, dal.pageSize),
	}
}

func (dal *dal) readPage(pageNum pgnum) (*page, error) {
	page := dal.allocateEmptyPage()

	// calculate offset
	offset := int64(page.num) * int64(dal.pageSize)

	// read data at offset
	_, err := dal.file.ReadAt(page.data, offset)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (dal *dal) writePage(page *page) error {
	// calculate offset
	offset := int64(page.num) * int64(dal.pageSize)

	// write at offset
	_, err := dal.file.WriteAt(page.data, offset)
	return err
}
