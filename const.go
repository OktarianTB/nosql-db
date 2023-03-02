package main

const (
	// metaPageNum is the maximum pgnum that is used by the db for its own purposes
	// page 0 is used as the header page
	metaPageNum = 0

	// size of a page number (in bytes)
	pageNumSize = 8

	// size of the node header
	nodeHeaderSize = 3
)
