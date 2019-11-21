package statistic

import (
	"container/list"
	"context"
	"time"

	"github.com/sjeandeaux/access-log-monitor/pkg/log"
)

// Aggregator aggregates log to generat a traffic stat
type Aggregator interface {
	Aggregate(context.Context, <-chan log.Entry) <-chan Traffic
}

type defaultAggregator struct {
	start    time.Time
	duration time.Duration
}

// NewDefaultAggregator default aggregator with a start time and duration
func NewDefaultAggregator(start time.Time, duration time.Duration) Aggregator {
	return &defaultAggregator{
		duration: duration,
		start:    start,
	}
}

//aggregatorStage the current stage for the aggregation
type aggregatorStage struct {
	//The current interval
	start        time.Time
	duration     time.Duration
	queueEntries *list.List
}

//newAggregatorStage initialize a stage
func newAggregatorStage(start time.Time, duration time.Duration) aggregatorStage {
	return aggregatorStage{
		start:        start,
		duration:     duration,
		queueEntries: list.New(),
	}
}

//add an entry to compute
func (as *aggregatorStage) add(e log.Entry) {
	as.queueEntries.PushBack(e)
}

//compute the current traffic.
// it creates statistics for the traffic between the start time and start + duration
// it moves the start time to start + duration
func (as *aggregatorStage) computeTraffic() Traffic {
	// calculate the statistics
	end := as.start.Add(as.duration)
	statToSend := NewTraffic(as.start, end)

	//process the entries
	for as.queueEntries.Len() > 0 {
		e := as.queueEntries.Front()
		v := e.Value.(log.Entry)

		//the event isn't in the slot
		if v.Datetime.After(statToSend.End) {
			break
		}
		statToSend.Add(&v)
		as.queueEntries.Remove(e) //it can remove it
	}
	//update the next slot
	as.start = end
	return *statToSend
}

// Aggregate it aggregates entries by slot time and sends the result in traffic channel
func (a *defaultAggregator) Aggregate(ctx context.Context, entries <-chan log.Entry) <-chan Traffic {
	stats := make(chan Traffic, 20)
	go func() {
		//The ticker with slot duration
		ticker := time.NewTicker(a.duration)
		stage := newAggregatorStage(a.start, a.duration)
		for {
			select {
			case entry := <-entries: //a entry log, it adds it in the queue
				stage.add(entry)
			case <-ticker.C:
				//it sends the stat
				stats <- stage.computeTraffic()
			case <-ctx.Done():
				return
			}
		}
	}()
	return stats
}
