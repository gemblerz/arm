package display

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/image/font/basicfont"
)

// PerformanceMonitor tracks and displays system performance metrics
type PerformanceMonitor struct {
	display        DisplayInterface
	cpuHistory     []float64
	memoryHistory  []float64
	updateHistory  []time.Duration
	lastUpdate     time.Time
	maxHistory     int
	currentMetrics SystemMetrics
}

// SystemMetrics represents current system performance metrics
type SystemMetrics struct {
	CPUUsage      float64
	MemoryUsage   float64
	MemoryTotal   int64
	MemoryUsed    int64
	UpdateRate    time.Duration
	MotorLoad     map[string]float64
	Temperature   map[string]float64
	Voltage       float64
	ErrorCount    int
	UptimeSeconds int64
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(display DisplayInterface) *PerformanceMonitor {
	return &PerformanceMonitor{
		display:        display,
		cpuHistory:     make([]float64, 0),
		memoryHistory:  make([]float64, 0),
		updateHistory:  make([]time.Duration, 0),
		maxHistory:     60, // Keep 60 data points
		currentMetrics: SystemMetrics{
			MotorLoad:   make(map[string]float64),
			Temperature: make(map[string]float64),
		},
	}
}

// UpdateMetrics updates the performance metrics
func (pm *PerformanceMonitor) UpdateMetrics(metrics SystemMetrics) {
	pm.currentMetrics = metrics
	
	// Add to history
	pm.addToHistory(&pm.cpuHistory, metrics.CPUUsage)
	pm.addToHistory(&pm.memoryHistory, metrics.MemoryUsage)
	
	// Simulate update rate calculation
	now := time.Now()
	if !pm.lastUpdate.IsZero() {
		updateDuration := now.Sub(pm.lastUpdate)
		pm.addToHistoryDuration(&pm.updateHistory, updateDuration)
	}
	pm.lastUpdate = now
}

// addToHistory adds a value to history and maintains max size
func (pm *PerformanceMonitor) addToHistory(history *[]float64, value float64) {
	*history = append(*history, value)
	if len(*history) > pm.maxHistory {
		*history = (*history)[1:]
	}
}

// addToHistoryDuration adds a duration to history
func (pm *PerformanceMonitor) addToHistoryDuration(history *[]time.Duration, value time.Duration) {
	*history = append(*history, value)
	if len(*history) > pm.maxHistory {
		*history = (*history)[1:]
	}
}

// DisplayCPUChart displays CPU usage chart
func (pm *PerformanceMonitor) DisplayCPUChart() {
	pm.display.Clear()
	pm.display.DrawText("CPU Usage", 0, 0, basicfont.Face7x13)
	
	// Draw current value
	currentText := fmt.Sprintf("%.1f%%", pm.currentMetrics.CPUUsage)
	pm.display.DrawText(currentText, 80, 0, basicfont.Face7x13)
	
	// Draw chart
	pm.drawLineChart(pm.cpuHistory, 12, 40, 100.0, "CPU %")
	
	pm.display.Update()
}

// DisplayMemoryChart displays memory usage chart
func (pm *PerformanceMonitor) DisplayMemoryChart() {
	pm.display.Clear()
	pm.display.DrawText("Memory Usage", 0, 0, basicfont.Face7x13)
	
	// Draw current value
	currentText := fmt.Sprintf("%.1f%%", pm.currentMetrics.MemoryUsage)
	pm.display.DrawText(currentText, 80, 0, basicfont.Face7x13)
	
	// Memory info
	if pm.currentMetrics.MemoryTotal > 0 {
		totalMB := pm.currentMetrics.MemoryTotal / (1024 * 1024)
		usedMB := pm.currentMetrics.MemoryUsed / (1024 * 1024)
		memText := fmt.Sprintf("%dMB/%dMB", usedMB, totalMB)
		pm.display.DrawText(memText, 0, 56, basicfont.Face7x13)
	}
	
	// Draw chart
	pm.drawLineChart(pm.memoryHistory, 12, 35, 100.0, "Mem %")
	
	pm.display.Update()
}

// DisplayUpdateRateChart displays update rate performance
func (pm *PerformanceMonitor) DisplayUpdateRateChart() {
	pm.display.Clear()
	pm.display.DrawText("Update Rate", 0, 0, basicfont.Face7x13)
	
	// Calculate average update rate
	if len(pm.updateHistory) > 0 {
		var total time.Duration
		for _, duration := range pm.updateHistory {
			total += duration
		}
		avg := total / time.Duration(len(pm.updateHistory))
		pm.display.DrawText(fmt.Sprintf("%.0fms", float64(avg.Nanoseconds())/1e6), 80, 0, basicfont.Face7x13)
	}
	
	// Convert duration history to float64 for chart
	rateHistory := make([]float64, len(pm.updateHistory))
	for i, duration := range pm.updateHistory {
		rateHistory[i] = float64(duration.Nanoseconds()) / 1e6 // Convert to milliseconds
	}
	
	// Draw chart (max 500ms scale)
	pm.drawLineChart(rateHistory, 12, 40, 500.0, "ms")
	
	pm.display.Update()
}

// DisplayMotorLoadChart displays motor load information
func (pm *PerformanceMonitor) DisplayMotorLoadChart() {
	pm.display.Clear()
	pm.display.DrawText("Motor Load", 0, 0, basicfont.Face7x13)
	
	// Display motor loads as bars
	motors := []string{"Pan", "Tilt", "Elbow", "Wrist1", "Wrist2", "Grip"}
	y := 12
	
	for i, motor := range motors {
		if i >= 5 { // Limit to fit on screen
			break
		}
		
		load := pm.currentMetrics.MotorLoad[motor]
		if load == 0 {
			load = float64(20 + i*10) // Demo values
		}
		
		// Motor name
		pm.display.DrawText(motor, 0, y, basicfont.Face7x13)
		
		// Load bar
		barWidth := 60
		loadWidth := int(load * float64(barWidth) / 100.0)
		pm.display.DrawRectangle(50, y, barWidth, 8, false)
		pm.display.DrawRectangle(50, y, loadWidth, 8, true)
		
		// Percentage
		pm.display.DrawText(fmt.Sprintf("%.0f%%", load), 115, y, basicfont.Face7x13)
		
		y += 10
	}
	
	pm.display.Update()
}

// DisplaySystemOverview displays a system overview
func (pm *PerformanceMonitor) DisplaySystemOverview() {
	pm.display.Clear()
	pm.display.DrawText("System Overview", 0, 0, basicfont.Face7x13)
	
	// Key metrics
	y := 12
	
	// CPU
	pm.display.DrawText(fmt.Sprintf("CPU: %.1f%%", pm.currentMetrics.CPUUsage), 0, y, basicfont.Face7x13)
	y += 10
	
	// Memory
	pm.display.DrawText(fmt.Sprintf("Mem: %.1f%%", pm.currentMetrics.MemoryUsage), 0, y, basicfont.Face7x13)
	y += 10
	
	// Voltage
	if pm.currentMetrics.Voltage > 0 {
		pm.display.DrawText(fmt.Sprintf("Volt: %.1fV", pm.currentMetrics.Voltage), 0, y, basicfont.Face7x13)
	} else {
		pm.display.DrawText("Volt: 12.1V", 0, y, basicfont.Face7x13) // Demo value
	}
	y += 10
	
	// Temperature
	temp := pm.currentMetrics.Temperature["CPU"]
	if temp == 0 {
		temp = 45.2 // Demo value
	}
	pm.display.DrawText(fmt.Sprintf("Temp: %.1f°C", temp), 0, y, basicfont.Face7x13)
	y += 10
	
	// Uptime
	uptime := pm.currentMetrics.UptimeSeconds
	if uptime == 0 {
		uptime = 3661 // Demo: 1 hour, 1 minute, 1 second
	}
	hours := uptime / 3600
	minutes := (uptime % 3600) / 60
	pm.display.DrawText(fmt.Sprintf("Up: %dh%dm", hours, minutes), 0, y, basicfont.Face7x13)
	
	pm.display.Update()
}

// drawLineChart draws a line chart for the given data
func (pm *PerformanceMonitor) drawLineChart(data []float64, startY, height int, maxValue float64, unit string) {
	if len(data) == 0 {
		return
	}
	
	chartWidth := 120
	chartHeight := height
	
	// Draw axes
	pm.display.DrawLine(8, startY+chartHeight, 8+chartWidth, startY+chartHeight) // X-axis
	pm.display.DrawLine(8, startY, 8, startY+chartHeight)                        // Y-axis
	
	// Draw scale labels
	pm.display.DrawText("0", 0, startY+chartHeight+5, basicfont.Face7x13)
	pm.display.DrawText(fmt.Sprintf("%.0f", maxValue), 0, startY-5, basicfont.Face7x13)
	
	// Calculate point spacing
	pointSpacing := float64(chartWidth) / float64(len(data)-1)
	if len(data) == 1 {
		pointSpacing = 0
	}
	
	// Draw data points and lines
	for i, value := range data {
		// Calculate position
		x := 8 + int(float64(i)*pointSpacing)
		normalizedValue := math.Min(value/maxValue, 1.0)
		y := startY + chartHeight - int(normalizedValue*float64(chartHeight))
		
		// Draw point
		pm.display.SetPixel(x, y, true)
		
		// Draw line to next point
		if i < len(data)-1 {
			nextValue := data[i+1]
			nextX := 8 + int(float64(i+1)*pointSpacing)
			nextNormalizedValue := math.Min(nextValue/maxValue, 1.0)
			nextY := startY + chartHeight - int(nextNormalizedValue*float64(chartHeight))
			
			pm.drawLine(x, y, nextX, nextY)
		}
	}
	
	// Draw current value
	if len(data) > 0 {
		currentValue := data[len(data)-1]
		valueText := fmt.Sprintf("%.1f%s", currentValue, unit)
		pm.display.DrawText(valueText, 80, startY+chartHeight+5, basicfont.Face7x13)
	}
}

// drawLine draws a line between two points using Bresenham's algorithm
func (pm *PerformanceMonitor) drawLine(x0, y0, x1, y1 int) {
	dx := absInt(x1 - x0)
	dy := absInt(y1 - y0)
	x, y := x0, y0
	n := 1 + dx + dy
	x_inc := 1
	if x1 < x0 {
		x_inc = -1
	}
	y_inc := 1
	if y1 < y0 {
		y_inc = -1
	}
	error := dx - dy
	
	dx *= 2
	dy *= 2
	
	for ; n > 0; n-- {
		pm.display.SetPixel(x, y, true)
		
		if error > 0 {
			x += x_inc
			error -= dy
		} else {
			y += y_inc
			error += dx
		}
	}
}

// absInt returns the absolute value of an integer
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetCurrentMetrics returns the current system metrics
func (pm *PerformanceMonitor) GetCurrentMetrics() SystemMetrics {
	return pm.currentMetrics
}

// SimulateMetrics generates simulated metrics for demo purposes
func (pm *PerformanceMonitor) SimulateMetrics() SystemMetrics {
	now := time.Now()
	
	// Simulate varying CPU usage
	cpuBase := 30.0
	cpuVariation := 20.0 * math.Sin(float64(now.Unix())/10.0)
	cpu := cpuBase + cpuVariation
	if cpu < 0 {
		cpu = 0
	}
	if cpu > 100 {
		cpu = 100
	}
	
	// Simulate memory usage
	memBase := 65.0
	memVariation := 10.0 * math.Sin(float64(now.Unix())/20.0)
	memory := memBase + memVariation
	
	// Simulate motor loads
	motorLoad := map[string]float64{
		"Pan":    30 + 20*math.Sin(float64(now.Unix())/5.0),
		"Tilt":   25 + 15*math.Cos(float64(now.Unix())/7.0),
		"Elbow":  40 + 25*math.Sin(float64(now.Unix())/6.0),
		"Wrist1": 20 + 10*math.Cos(float64(now.Unix())/8.0),
		"Wrist2": 35 + 20*math.Sin(float64(now.Unix())/9.0),
		"Grip":   15 + 10*math.Cos(float64(now.Unix())/4.0),
	}
	
	// Simulate temperatures
	temperature := map[string]float64{
		"CPU":    42.0 + 8.0*math.Sin(float64(now.Unix())/30.0),
		"Motor1": 38.0 + 5.0*math.Cos(float64(now.Unix())/25.0),
		"Motor2": 41.0 + 6.0*math.Sin(float64(now.Unix())/35.0),
	}
	
	return SystemMetrics{
		CPUUsage:      cpu,
		MemoryUsage:   memory,
		MemoryTotal:   2048 * 1024 * 1024, // 2GB
		MemoryUsed:    int64(memory * 2048 * 1024 * 1024 / 100),
		UpdateRate:    200 * time.Millisecond,
		MotorLoad:     motorLoad,
		Temperature:   temperature,
		Voltage:       12.1 + 0.3*math.Sin(float64(now.Unix())/40.0),
		ErrorCount:    0,
		UptimeSeconds: int64(now.Unix() % 86400), // Simulate daily uptime
	}
}
