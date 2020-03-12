package controller

import (
	"fmt"
	"runtime"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

/*
	Implement terminal-based controller
	xterm color reference https://jonasjacek.github.io/colors/
*/

const (
	// terminalWidth, width of terminal UI
	terminalWidth     = 120
	heapAllocBarCount = 6
)

type ControllerInterface interface {
	Render(*runtime.MemStats)
	Resize()
}

type Controller struct {
	Grid *ui.Grid

	HeapObjectsSparkline      *widgets.Sparkline
	HeapObjectsSparklineGroup *widgets.SparklineGroup
	HeapObjectsData           *statRing

	SysText *widgets.Paragraph // OK

	// GCCPUFraction *widgets.Gauge // TODO crashes the whole application!

	HeapAllocBarChart     *widgets.BarChart
	HeapAllocBarChartData *statRing

	HeapPie *widgets.PieChart // OK
}

func NewController() *Controller {
	ctl := &Controller{
		Grid: ui.NewGrid(),

		HeapObjectsSparkline: widgets.NewSparkline(),
		HeapObjectsData:      newChartRing(terminalWidth),

		SysText: widgets.NewParagraph(),

		// GCCPUFraction: widgets.NewGauge(),

		HeapAllocBarChart:     widgets.NewBarChart(),
		HeapAllocBarChartData: newChartRing(heapAllocBarCount),

		HeapPie: widgets.NewPieChart(),
	}
	ctl.initUI()
	return ctl
}

func (c *Controller) initUI() {
	c.resize()

	c.HeapObjectsSparkline.LineColor = ui.Color(89) // xterm color DeepPink4
	c.HeapObjectsSparklineGroup = widgets.NewSparklineGroup(c.HeapObjectsSparkline)

	c.SysText.Title = "Sys, the total bytes of memory obtained from the OS"
	c.SysText.PaddingLeft = 25
	c.SysText.PaddingTop = 3

	// c.GCCPUFraction.Title = "GCCPUFraction 0%~100%"
	// c.GCCPUFraction.BarColor = ui.Color(50) // xterm color Cyan2

	c.HeapAllocBarChart.BarGap = 2
	c.HeapAllocBarChart.BarWidth = 8
	c.HeapAllocBarChart.Title = "HeapAlloc, bytes of allocated heap objects"
	c.HeapAllocBarChart.NumFormatter = func(f float64) string { return "" }

	c.HeapPie.Title = "HeapInuse vs HeapIdle"
	c.HeapPie.LabelFormatter = func(idx int, _ float64) string { return []string{"Idle", "Inuse"}[idx] }

	// TODO use a flag/env-var to choose between minimal and extended configuration
	c.Grid.Set(
		// MINIMAL
		// ui.NewRow(.5,
		// 	ui.NewCol(.5, c.SysText),
		// 	ui.NewCol(.5, c.HeapPie),
		// ),
		// ui.NewRow(.5, c.HeapAllocBarChart),

		// EXTENDED
		ui.NewRow(.2, c.HeapObjectsSparklineGroup),
		ui.NewRow(.8,
			ui.NewCol(.5,
				ui.NewRow(.2, c.SysText),
				ui.NewRow(.8, c.HeapAllocBarChart),
			),
			ui.NewCol(.5, c.HeapPie),
		),
	)
}

func (c *Controller) Render(data *runtime.MemStats) {
	c.HeapObjectsSparklineGroup.Title = fmt.Sprintf("HeapObjects, live heap object count: %d", data.HeapObjects)
	c.HeapObjectsSparkline.Data = c.HeapObjectsData.normalizedData()
	c.HeapObjectsData.push(data.HeapObjects)

	c.SysText.Text = fmt.Sprint(byteCountBinary(data.Sys))

	// c.GCCPUFraction.Label = fmt.Sprintf("%.2f%%", data.GCCPUFraction*100)
	// c.GCCPUFraction.Percent = fNormalize(data)

	c.HeapAllocBarChartData.push(data.HeapAlloc)
	c.HeapAllocBarChart.Data = c.HeapAllocBarChartData.convertData()
	c.HeapAllocBarChart.Labels = nil
	for _, v := range c.HeapAllocBarChart.Data {
		c.HeapAllocBarChart.Labels = append(c.HeapAllocBarChart.Labels, byteCountBinary(uint64(v)))
	}

	c.HeapPie.Data = []float64{float64(data.HeapIdle), float64(data.HeapInuse)}

	ui.Render(c.Grid)
}

func (c *Controller) Resize() {
	c.resize()
	ui.Render(c.Grid)
}

func (c *Controller) resize() {
	_, h := ui.TerminalDimensions()
	c.Grid.SetRect(0, 0, terminalWidth, h)
}

func fNormalize(data *runtime.MemStats) int {
	f := data.GCCPUFraction
	if f < 0.01 {
		for f < 1 {
			f = f * 10.0
		}
	}
	return int(f)
}

// source: https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func byteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
