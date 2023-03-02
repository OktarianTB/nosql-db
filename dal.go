package main

import (
	"errors"
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

	*meta
	*freelist
}

func newDal(path string) (*dal, error) {
	dal := &dal{
		meta:     newEmptyMeta(),
		pageSize: os.Getpagesize(),
	}

	if _, err := os.Stat(path); err == nil {
		// file already exists
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = dal.close()
			return nil, err
		}

		meta, err := dal.readMeta()
		if err != nil {
			return nil, err
		}
		dal.meta = meta

		freelist, err := dal.readFreelist()
		if err != nil {
			return nil, err
		}
		dal.freelist = freelist
	} else if errors.Is(err, os.ErrNotExist) {
		// file does not exist
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = dal.close()
			return nil, err
		}

		dal.freelist = newFreelist()
		dal.freelistPage = dal.getNextPage()
		_, err := dal.writeFreelist()
		if err != nil {
			return nil, err
		}

		_, err = dal.writeMeta(dal.meta)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
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
		data: make([]byte, dal.pageSize, dal.pageSize),
	}
}

func (dal *dal) writeNode(node *Node) (*Node, error) {
	page := dal.allocateEmptyPage()
	if node.pageNum == 0 {
		page.num = dal.getNextPage()
		node.pageNum = page.num
	} else {
		page.num = node.pageNum
	}

	page.data = node.serialize(page.data)

	err := dal.writePage(page)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (dal *dal) getNode(pageNum pgnum) (*Node, error) {
	page, err := dal.readPage(pageNum)
	if err != nil {
		return nil, err
	}
	node := newEmptyNode()
	node.deserialize(page.data)
	node.pageNum = pageNum
	
	return node, nil
}

func (dal *dal) deleteNode(pageNum pgnum) {
	dal.releasePage(pageNum)
}

func (dal *dal) readPage(pageNum pgnum) (*page, error) {
	page := dal.allocateEmptyPage()

	// calculate offset
	offset := int64(pageNum) * int64(dal.pageSize)

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

func (dal *dal) writeFreelist() (*page, error) {
	page := dal.allocateEmptyPage()
	page.num = dal.freelistPage
	dal.freelist.serialize(page.data)

	err := dal.writePage(page)
	if err != nil {
		return nil, err
	}

	dal.freelistPage = page.num
	return page, nil
}

func (dal *dal) readFreelist() (*freelist, error) {
	page, err := dal.readPage(dal.freelistPage)
	if err != nil {
		return nil, err
	}

	freelist := newFreelist()
	freelist.deserialize(page.data)
	return freelist, nil
}

func (dal *dal) writeMeta(meta *meta) (*page, error) {
	page := dal.allocateEmptyPage()
	page.num = metaPageNum
	meta.serialize(page.data)

	err := dal.writePage(page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (dal *dal) readMeta() (*meta, error) {
	page, err := dal.readPage(metaPageNum)
	if err != nil {
		return nil, err
	}

	meta := newEmptyMeta()
	meta.deserialize(page.data)
	return meta, nil
}
