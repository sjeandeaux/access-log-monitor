package file

import (
	"context"

	"github.com/hpcloud/tail"
	log "github.com/sirupsen/logrus"
)

// Tailer the watcher it sends message in channel
type Tailer interface {
	Tail(context.Context) (<-chan string, <-chan error)
}

// DefaulTailer the file
type defaulTailer struct {
	tailer *tail.Tail
}

// NewDefaultTailer creates a watcher on the file.
func NewDefaultTailer(file string) (Tailer, error) {
	t, err := tail.TailFile(file, tail.Config{ReOpen: true, Follow: true, Logger: tail.DiscardingLogger})
	if err != nil {
		return nil, err
	}
	return &defaulTailer{tailer: t}, err
}

// Tail it runs the process of tailing.
// there are two channels:
// - chan string it is for the lines
// - chan error it is for the errors when it tails the file
func (dw *defaulTailer) Tail(ctx context.Context) (<-chan string, <-chan error) {
	lines := make(chan string, 20) //TODO identify the right number for the buffering
	errs := make(chan error, 20)   //TODO identify the right number for the buffering

	//It listen to the context and the tailing
	go func() {
		//it has to close the channels when it exists for this function
		defer func() {
			close(lines)
			close(errs)
		}()
		for {
			select {
			case <-ctx.Done():
				log.Info("tail context is done")
				return
			case line, ok := <-dw.tailer.Lines:
				if !ok { //the channel is closed
					return
				}
				if line.Err != nil {
					errs <- line.Err
				} else {
					lines <- line.Text
				}
			}
		}
	}()
	return lines, errs
}
