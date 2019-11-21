// Package statistic contains the aggregator to generate statistic on the traffic
package statistic

import (
	"net/http"
	"sort"
	"time"

	"github.com/sjeandeaux/access-log-monitor/pkg/log"
)

// Traffic contains statistics about traffic
type Traffic struct {
	//The start time of captation
	Start time.Time
	//The end time of captation
	End                    time.Time
	HitBySections          map[string]int64
	HitByUsers             map[string]int64
	HitUnauthorizedByUsers map[string]int64
	HitForbiddenByUsers    map[string]int64
	HitByMethods           map[string]int64
	HitByHosts             map[string]int64

	//Nb hits
	Nb    int64
	Nb2XX int64
	Nb3XX int64
	Nb4XX int64
	Nb5XX int64

	//The global size
	Size int64
}

// NewTraffic create a traffic static (init map)
func NewTraffic(start, end time.Time) *Traffic {
	return &Traffic{
		HitBySections:          make(map[string]int64),
		HitByUsers:             make(map[string]int64),
		HitUnauthorizedByUsers: make(map[string]int64),
		HitForbiddenByUsers:    make(map[string]int64),
		HitByMethods:           make(map[string]int64),
		HitByHosts:             make(map[string]int64),
		Start:                  start,
		End:                    end,
	}
}

// Add it adds a entry in statistic
// If you want to create a Traffic by your own, don't forget to init the map.
func (t *Traffic) Add(le *log.Entry) {
	if le == nil {
		return
	}
	//track only log between start and end time.
	if le.Datetime.Before(t.Start) || le.Datetime.After(t.End) {
		return
	}

	t.Nb++
	t.Size += le.Size
	t.HitBySections[le.Section()]++
	t.HitByUsers[le.AuthUser]++
	t.HitByMethods[le.Method]++
	t.HitByHosts[le.RemoteHost]++

	if le.Status == http.StatusForbidden {
		t.HitForbiddenByUsers[le.AuthUser]++
	}

	if le.Status == http.StatusUnauthorized {
		t.HitUnauthorizedByUsers[le.AuthUser]++
	}

	switch status := le.Status; {
	case status >= 200 && status <= 299:
		t.Nb2XX++
	case status >= 300 && status <= 399:
		t.Nb3XX++
	case status >= 400 && status <= 499:
		t.Nb4XX++
	case status >= 500 && status <= 599:
		t.Nb5XX++
	}
}

//Hit the name ex POST or a section name and the number of hits
type Hit struct {
	Name string
	Hits int64
}

// ByHits implements sort.Interface for []hit based on
// the hists field.
type ByHits []Hit

func (a ByHits) Len() int           { return len(a) }
func (a ByHits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHits) Less(i, j int) bool { return a[i].Hits > a[j].Hits }

//Order it orders the map
func Order(hits map[string]int64) []Hit {
	if hits == nil {
		return nil
	}
	sortedHits := make([]Hit, len(hits))
	i := 0
	for k, v := range hits {
		sortedHits[i] = Hit{Name: k, Hits: v}
		i++
	}
	sort.Sort(ByHits(sortedHits))
	return sortedHits
}
