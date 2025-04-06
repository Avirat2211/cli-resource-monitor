package main

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	cpuLine := widgets.NewPlot()
	cpuLine.PlotType = widgets.LineChart
	cpuLine.Title = "CPU Usage (%)"
	cpuLine.SetRect(0, 0, 80, 15)
	cpuLine.AxesColor = ui.ColorWhite
	cpuLine.LineColors[0] = ui.ColorCyan
	cpuLine.Data = [][]float64{{0, 0}}
	// cpuLine.Marker = widgets.MarkerDot

	memGauge := widgets.NewGauge()
	memGauge.Title = "Memory Usage"
	memGauge.SetRect(0, 15, 80, 20)
	memGauge.BarColor = ui.ColorGreen
	memGauge.BorderStyle.Fg = ui.ColorWhite
	memGauge.TitleStyle.Fg = ui.ColorCyan

	ui.Render(cpuLine, memGauge)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	cpuData := make([]float64, 0, 50)
	cpuData = append(cpuData, 0)
	cpuData = append(cpuData, 0)
	e := ui.PollEvents()
	for {
		select {
		case <-ticker.C:
			percent, _ := cpu.Percent(0, false)
			if len(cpuData) >= 50 {
				cpuData = cpuData[1:]
			}
			cpuData = append(cpuData, percent[0])
			cpuLine.Data[0] = cpuData
			cpuLine.Title = "CPU Usage (%) - " + fmt.Sprintf("%.2f%%", percent[0])
			vmStat, _ := mem.VirtualMemory()
			memGauge.Percent = int(vmStat.UsedPercent)
			memGauge.Label = fmt.Sprintf("%.2f%% (%.2f GB / %.2f GB)",
				vmStat.UsedPercent,
				float64(vmStat.Used)/1e9,
				float64(vmStat.Total)/1e9,
			)
			ui.Render(cpuLine, memGauge)
		case ev := <-e:
			if ev.Type == ui.KeyboardEvent {
				if ev.ID == "q" {
					return
				}
			}
		}
	}
}
