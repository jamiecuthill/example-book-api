package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnsBookNameForISBN(t *testing.T) {
	library := &bookrepo{
		server: &url.URL{},
	}

	name, err := library.GetBookName("9780134190440")
	assert.NoError(t, err)
	assert.Equal(t, "The Go programming language", name)
}

var fixtureISBNResponse = `{
	"ISBN:9780134190440" : {
	   "key" : "/books/OL27191446M",
	   "notes" : "\"First printing, October 2015\"--Title page verso.\n\nIncludes index.",
	   "pagination" : "xvii, 380 pages",
	   "number_of_pages" : 380,
	   "title" : "The Go programming language",
	   "publish_date" : "2015",
	   "url" : "https://openlibrary.org/books/OL27191446M/The_Go_programming_language",
	   "authors" : [
		  {
			 "url" : "https://openlibrary.org/authors/OL7603397A/Alan_A._A._Donovan",
			 "name" : "Alan A. A. Donovan"
		  }
	   ],
	   "identifiers" : {
		  "isbn_10" : [
			 "0134190440"
		  ],
		  "oclc" : [
			 "951142414",
			 "903635603"
		  ],
		  "isbn_13" : [
			 "9780134190440"
		  ],
		  "lccn" : [
			 "2015950709"
		  ],
		  "openlibrary" : [
			 "OL27191446M"
		  ]
	   },
	   "subjects" : [
		  {
			 "url" : "https://openlibrary.org/subjects/go_(computer_program_language)",
			 "name" : "Go (Computer program language)"
		  },
		  {
			 "url" : "https://openlibrary.org/subjects/programming",
			 "name" : "Programming"
		  },
		  {
			 "url" : "https://openlibrary.org/subjects/open_source_software",
			 "name" : "Open source software"
		  }
	   ],
	   "by_statement" : "Alan A. A. Donovan, Brian W. Kernighan",
	   "classifications" : {
		  "dewey_decimal_class" : [
			 "005.13/3"
		  ]
	   }
	}
 }`
