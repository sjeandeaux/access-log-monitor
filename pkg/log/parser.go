package log

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Parser parses of line
type Parser interface {
	Parse(string) (*Entry, error)
}

type regexParser struct {
	regex *regexp.Regexp
}

var _ Parser = &regexParser{}

//NewDefaultParser the default parser
func NewDefaultParser() Parser {
	re := regexp.MustCompile(`^(?P<RemoteHost>\S+) (?P<RemoteLogName>\S+) (?P<AuthUser>\S+) \[(?P<Datetime>[^\]]+)\] "(?P<Method>[A-Z]+) (?P<Request>[^ "]+)? HTTP\/[0-9.]+" (?P<Status>[0-9]{3}) (?P<Size>[0-9]+|-)`)
	return &regexParser{regex: re}
}

func (r *regexParser) Parse(line string) (*Entry, error) {
	const (
		groupRemoteHost    = 1
		groupRemoteLogName = 2
		groupAuthUser      = 3
		groupDatetime      = 4
		groupMethod        = 5
		groupRequest       = 6
		groupStatus        = 7
		groupSize          = 8
		nbGroups           = 9
		layoutTime         = "02/Jan/2006:15:04:05 -0700"
	)
	values := r.regex.FindStringSubmatch(line)
	if len(values) != nbGroups {
		return nil, fmt.Errorf("could not parse the line %q", line)
	}

	le := &Entry{}
	le.RemoteHost = values[groupRemoteHost]
	le.RemoteLogName = values[groupRemoteLogName]
	le.AuthUser = values[groupAuthUser]

	datetime, err := time.Parse(layoutTime, values[groupDatetime])
	if err != nil {
		return nil, fmt.Errorf("could not parse the line %q time: %v", line, err)
	}
	le.Datetime = datetime

	le.Method = values[groupMethod]
	le.Request = values[groupRequest]

	status, _ := strconv.ParseInt(values[groupStatus], 10, 16)
	le.Status = int16(status)

	if values[8] != "-" {
		size, _ := strconv.ParseInt(values[groupSize], 10, 64)
		le.Size = size
	}

	return le, nil
}
