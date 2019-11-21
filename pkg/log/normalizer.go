package log

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

// Normalizer it normalize the line in Entry
type Normalizer interface {
	Normalize(context.Context, <-chan string) (<-chan Entry, <-chan error)
}

//defaultNormalizer the default implementation
type defaultNormalizer struct {
	parser Parser
}

// NewDefaultNormalizer creates a normalizer on string
func NewDefaultNormalizer() Normalizer {
	return &defaultNormalizer{
		parser: NewDefaultParser(),
	}
}

// Normalize it normalizes the line in Entry in *Entry channel.
// If it fails it sends an error
func (dn *defaultNormalizer) Normalize(ctx context.Context, lines <-chan string) (<-chan Entry, <-chan error) {
	normalizedLines := make(chan Entry, 20) //TODO identify the right number for the buffering
	errs := make(chan error, 20)            //TODO identify the right number for the buffering

	//It listen to the context and the tailing
	go func() {
		//it has to close the channels when it exists for this function
		defer func() {
			close(normalizedLines)
			close(errs)
		}()
		for {
			select {
			case <-ctx.Done():
				log.Info("normalize context is done")
				return
			case line, ok := <-lines:
				if !ok { //the channel is closed
					return
				}
				if entry, err := dn.parser.Parse(line); err == nil && entry != nil { //no problem it sends the normalized line in the channel of entry
					normalizedLines <- *entry
				} else if err != nil { //Oupss an issue with the line
					errs <- err
				} else {
					// when the entry is
					errs <- errors.New("should never never happen")
				}

			}
		}
	}()
	return normalizedLines, errs
}
