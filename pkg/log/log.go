// Package log contains the parsing and the normalization.
package log

import "time"

// Entry contains the line information
type Entry struct {
	RemoteHost    string
	RemoteLogName string
	AuthUser      string
	Datetime      time.Time

	Method  string
	Request string
	Status  int16
	Size    int64
}

// Section returns the section of the request
// /pages/subpage1/..../subpageN => /pages
// /pages/ => /pages
// /pages => /pages
// / => /
//  => /
func (e *Entry) Section() string {
	const slash = "/"
	size := len(e.Request)
	//
	if size == 0 {
		return slash
	}

	indexSlash := 0
	for i := 0; i < size; i++ {
		if e.Request[i:i+1] == slash {
			indexSlash++
			//it meets the second slash, it can return the section
			if indexSlash == 2 {
				return e.Request[:i]
			}
		}
	}
	return e.Request
}
