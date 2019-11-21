package display

import (
	"context"
	"fmt"
	"time"

	"github.com/marcusolsson/tui-go"
	"github.com/sjeandeaux/access-log-monitor/pkg/alerting"
	"github.com/sjeandeaux/access-log-monitor/pkg/statistic"
)

const topSize = 5

// TerminalUI the UI in the terminal
type TerminalUI struct {
	ui           tui.UI
	title        *tui.Box
	summary      *tui.Label
	sections     *tui.List
	users        *tui.List
	methods      *tui.List
	hosts        *tui.List
	alerts       *tui.Label
	parsingError *tui.Label
	quit         *tui.Label
}

func createTopBox(title string, widget tui.Widget) *tui.Box {
	box := tui.NewVBox(widget)
	box.SetBorder(true)
	box.SetTitle(title)
	return box
}

// NewTerminalUI create a display
func NewTerminalUI() (*TerminalUI, error) {
	const (
		titleInformation    = "Information"
		titleAlerts         = "Alerts"
		titleFailureParsing = "Failure parsing"

		textToFill    = "..."
		textEscToQuit = "Esc to quit"
	)
	t := tui.NewTheme()
	t.SetStyle("normal", tui.Style{Bg: tui.ColorWhite, Fg: tui.ColorBlack})
	t.SetStyle("label.fatal", tui.Style{Bg: tui.ColorDefault, Fg: tui.ColorRed})

	//title and summary
	summary := tui.NewLabel(textToFill)
	summary.SetSizePolicy(tui.Expanding, tui.Expanding)
	title := tui.NewVBox(summary)
	title.SetBorder(true)
	title.SetTitle(titleInformation)

	//the tops
	sections := tui.NewList()
	sections.AddItems(make([]string, topSize)...)
	users := tui.NewList()
	users.AddItems(make([]string, topSize)...)
	methods := tui.NewList()
	methods.AddItems(make([]string, topSize)...)
	hosts := tui.NewList()
	hosts.AddItems(make([]string, topSize)...)
	tops := tui.NewHBox(
		createTopBox(fmt.Sprintf("Top %d sections", topSize), sections),
		createTopBox(fmt.Sprintf("Top %d users", topSize), users),
		createTopBox(fmt.Sprintf("Top %d methods", topSize), methods),
		createTopBox(fmt.Sprintf("Top %d hosts", topSize), hosts))

	//alerting
	alerts := tui.NewLabel(textToFill)
	alerts.SetSizePolicy(tui.Expanding, tui.Expanding)
	alertTitle := tui.NewVBox(alerts)
	alertTitle.SetBorder(true)
	alertTitle.SetTitle(titleAlerts)

	//parsing error
	parsingError := tui.NewLabel(textToFill)
	parsingError.SetSizePolicy(tui.Expanding, tui.Expanding)
	parsingErrorTitle := tui.NewVBox(parsingError)
	parsingErrorTitle.SetBorder(true)
	parsingErrorTitle.SetTitle(titleFailureParsing)
	parsingError.SetStyleName("fatal")

	quit := tui.NewLabel(textEscToQuit)
	quit.SetSizePolicy(tui.Expanding, tui.Expanding)

	root := tui.NewVBox(title, tops, alertTitle, parsingErrorTitle, quit)

	ui, err := tui.New(root)
	if err != nil {
		return nil, err
	}
	ui.SetTheme(t)
	return &TerminalUI{
		ui:           ui,
		title:        title,
		parsingError: parsingError,
		alerts:       alerts,
		summary:      summary,
		quit:         quit,
		sections:     sections,
		users:        users,
		methods:      methods,
		hosts:        hosts,
	}, nil

}

// Run the traffic, it blocks
func (t *TerminalUI) Run(ctx context.Context, stats <-chan statistic.Traffic, alerts <-chan alerting.Alert, errChan <-chan error) error {

	ctx, cancel := context.WithCancel(ctx)
	t.ui.SetKeybinding("Esc", func() {
		cancel()
		t.ui.Quit()
	})

	//event stat, error and alert
	go func() {

		for {
			select {
			case <-ctx.Done():
				return
			case alert := <-alerts:
				text := ""
				styleName := ""
				if alert.Status == alerting.HighTraffic {
					styleName = "fatal"
					text = fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", alert.Hits, alert.Date.Format(time.Stamp))
				} else {
					styleName = "normal"
					text = fmt.Sprintf("The traffic recovered - hits = %d, triggered at %s", alert.Hits, alert.Date.Format(time.Stamp))
				}
				t.alerts.SetStyleName(styleName)
				t.alerts.SetText(text)
				t.ui.Repaint()
			case stat := <-stats:
				t.title.SetTitle(fmt.Sprintf("Between %s and %s", stat.Start.Format(time.Stamp), stat.End.Format(time.Stamp)))
				t.summary.SetText(fmt.Sprintf("Nb: %d Nb2XX: %d Nb3XX: %d Nb4XX: %d Nb5XX: %d Size: %d", stat.Nb, stat.Nb2XX, stat.Nb3XX, stat.Nb4XX, stat.Nb5XX, stat.Size))

				t.sections.RemoveItems()
				t.sections.AddItems(topInRow(stat.HitBySections)...)

				t.users.RemoveItems()
				t.users.AddItems(topInRow(stat.HitByUsers)...)

				t.methods.RemoveItems()
				t.methods.AddItems(topInRow(stat.HitByMethods)...)

				t.hosts.RemoveItems()
				t.hosts.AddItems(topInRow(stat.HitByHosts)...)
				t.ui.Repaint()
			case err := <-errChan:
				if err != nil {
					t.parsingError.SetText(err.Error())
					t.ui.Repaint()
				}
			}
		}
	}()

	return t.ui.Run()

}

// topInRow returns the top X to hits map in printable arrays [key:value,...,key:value]
func topInRow(hits map[string]int64) []string {
	summary := make([]string, topSize)
	for i, s := range statistic.Order(hits) {
		if i == topSize {
			break
		}
		summary[i] = fmt.Sprintf("%s:%d", s.Name, s.Hits)
	}
	return summary
}
