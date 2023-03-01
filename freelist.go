package main

// freelist manages the free and used pages
// new page ids are first given from the releasedPageIDs to avoid growing the file
// if it's empty, then maxPage is incremented and a new page is created thus increasing the file size
type freelist struct {
	// maxPage holds the latest page num allocated
	maxPage pgnum
	// releasedPages holds all the ids that were released during delete
	releasedPages []pgnum
}

// metaPage is the maximum pgnum that is used by the db for its own purposes
// page 0 is used as the header page
const metaPage = 0

func newFreelist() *freelist {
	return &freelist{
		maxPage:       metaPage,
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

func (freelist *freelist) releasePagepage(pageId pgnum) {
	freelist.releasedPages = append(freelist.releasedPages, pageId)
}
