package main

import (
	"encoding/binary"
)

// freelist manages the free and used pages
// new page ids are first given from the releasedPageIDs to avoid growing the file
// if it's empty, then maxPage is incremented and a new page is created thus increasing the file size
type freelist struct {
	// maxPage holds the latest page num allocated
	maxPage pgnum
	// releasedPages holds all the ids that were released during delete
	releasedPages []pgnum
}

func newFreelist() *freelist {
	return &freelist{
		maxPage:       metaPageNum,
		releasedPages: []pgnum{},
	}
}

func (freelist *freelist) getNextPage() pgnum {
	// if possible, fetch pages first from the released pages
	if len(freelist.releasedPages) > 0 {
		pageId := freelist.releasedPages[len(freelist.releasedPages)-1]
		freelist.releasedPages = freelist.releasedPages[:len(freelist.releasedPages)-1]
		return pageId
	}
	freelist.maxPage += 1
	return freelist.maxPage
}

func (freelist *freelist) releasePage(pageId pgnum) {
	freelist.releasedPages = append(freelist.releasedPages, pageId)
}

func (freelist *freelist) serialize(buffer []byte) []byte {
	pos := 0

	// serialize maxPage
	binary.LittleEndian.PutUint16(buffer, uint16(freelist.maxPage))
	pos += 2

	// serialize releasedPages
	binary.LittleEndian.PutUint16(buffer[pos:], uint16(len(freelist.releasedPages)))
	pos += 2

	for _, page := range freelist.releasedPages {
		binary.LittleEndian.PutUint64(buffer[pos:], uint64(page))
		pos += pageNumSize
	}

	return buffer
}

func (freelist *freelist) deserialize(buffer []byte) {
	pos := 0

	// deserialize maxPage
	freelist.maxPage = pgnum(binary.LittleEndian.Uint16(buffer[pos:]))
	pos += 2

	// deserialiwe releasedPages
	releasedPagesCount := int(binary.LittleEndian.Uint16(buffer[pos:]))
	pos += 2

	for i := 0; i < releasedPagesCount; i++ {
		freelist.releasedPages = append(freelist.releasedPages, pgnum(binary.LittleEndian.Uint64(buffer[pos:])))
		pos += pageNumSize
	}
}
