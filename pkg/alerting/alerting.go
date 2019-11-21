package alerting

import (
	"container/list"
	"context"
	"time"

	"github.com/sjeandeaux/access-log-monitor/pkg/statistic"
)

// Status alert
type Status int

const (
	//Undefined nothing happens
	Undefined Status = iota
	//HighTraffic too much traffic
	HighTraffic
	// Recovered to the normal
	Recovered
)

// Alert sent by the listener
type Alert struct {
	//Status of alert
	Status Status
	//Hits number of connection
	Hits int64
	//Date when it triggered
	Date time.Time
}

//Alertor it watches the traffics to send the alert when it is necesary
type Alertor interface {
	Alert(context.Context, <-chan statistic.Traffic) <-chan Alert
}

//defaultAlertor default implementation
type defaultAlertor struct {
	nbTraffics int64
	threshold  float64
	duration   time.Duration
}

//NewDefaultAlertor creates an alertor
func NewDefaultAlertor(nbTraffics int64, threshold float64, duration time.Duration) Alertor {
	return &defaultAlertor{
		nbTraffics: nbTraffics,
		threshold:  threshold,
		duration:   duration,
	}
}

//alertorStage manages the stage alert to identify if it can trigger the alert or not
type alertorStage struct {
	nbTraffics    int64
	duration      time.Duration
	threshold     float64
	nbHits        int64
	hitsByTraffic *list.List
	currentStatus Status
}

//newAlertorStage creates a alertorStage
func newAlertorStage(nbTraffics int64, threshold float64, duration time.Duration) alertorStage {
	return alertorStage{
		nbTraffics:    nbTraffics,
		threshold:     threshold,
		duration:      duration,
		hitsByTraffic: list.New(),
		currentStatus: Undefined,
	}
}

//add a traffic stat
func (a *alertorStage) add(stat statistic.Traffic) {
	//It reaches the number of traffic stat
	if a.nbTraffics == int64(a.hitsByTraffic.Len()) {
		//it removes the older
		if older := a.hitsByTraffic.Front(); older != nil {
			a.hitsByTraffic.Remove(older)
			outOfDate := older.Value.(int64)
			//decrease its number of hits
			a.nbHits -= outOfDate
		}

	}
	//add the current traffic
	a.nbHits += stat.Nb
	a.hitsByTraffic.PushBack(stat.Nb)
}

//generateAlert it generates the alert if needed else returns nil
func (a *alertorStage) generateAlert() *Alert {
	hitsAverage := float64(a.nbHits) / a.duration.Seconds()

	switch a.currentStatus {
	case HighTraffic: //Check if it recovered
		if hitsAverage < a.threshold {
			a.currentStatus = Recovered
			return &Alert{Status: Recovered, Hits: int64(a.nbHits), Date: time.Now()}
		}
	default: //Check if it is in High Traffic
		if hitsAverage >= a.threshold {
			a.currentStatus = HighTraffic
			return &Alert{Status: HighTraffic, Hits: int64(a.nbHits), Date: time.Now()}
		}
	}

	return nil
}

//Alert it listens to the traffic and sends alert if there is hight traffic or recovery
func (d *defaultAlertor) Alert(ctx context.Context, stats <-chan statistic.Traffic) <-chan Alert {
	alerts := make(chan Alert, 20) //TODO identify the right number for the buffering

	go func() {
		defer close(alerts)
		stage := newAlertorStage(d.nbTraffics, d.threshold, d.duration)
		for {
			select {
			case <-ctx.Done():
				return
			case stat := <-stats:
				stage.add(stat)
				if alert := stage.generateAlert(); alert != nil {
					alerts <- *alert
				}
			}
		}
	}()
	return alerts
}
