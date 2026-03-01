package display

import (
	"golang.org/x/image/font"
	"testing"
	"time"
)

// Mock display for testing
type MockDisplay struct {
	width  int
	height int
	pixels [][]bool
	clear  bool
}

func NewMockDisplay() *MockDisplay {
	md := &MockDisplay{
		width:  128,
		height: 64,
		pixels: make([][]bool, 64),
	}
	md.init()
	return md
}

func (md *MockDisplay) init() {
	for i := range md.pixels {
		md.pixels[i] = make([]bool, 128)
	}
}

func (md *MockDisplay) GetWidth() int  { return md.width }
func (md *MockDisplay) GetHeight() int { return md.height }
func (md *MockDisplay) Clear() error {
	md.init()
	md.clear = true
	return nil
}
func (md *MockDisplay) Display() error {
	return nil
}
func (md *MockDisplay) SetPixel(x, y int, on bool) {
	if x >= 0 && x < md.width && y >= 0 && y < md.height {
		md.pixels[y][x] = on
	}
}
func (md *MockDisplay) DrawText(text string, x, y int, fontFace font.Face) {}
func (md *MockDisplay) SetTextColor(white bool)                            {}
func (md *MockDisplay) SetTextSize(size int)                               {}
func (md *MockDisplay) DrawLine(x0, y0, x1, y1 int)                         {}
func (md *MockDisplay) DrawRect(x, y, width, height int)                     {}
func (md *MockDisplay) FillRect(x, y, width, height int)                     {}
func (md *MockDisplay) DrawCircle(x, y, radius int)                          {}
func (md *MockDisplay) Close() error                                         { return nil }
func (md *MockDisplay) Initialize() error                                    { return nil }
func (md *MockDisplay) DrawProgressBar(x, y, width, height int, progress float64) {}
func (md *MockDisplay) Update() error                                        { return nil }
func (md *MockDisplay) DrawRectangle(x, y, width, height int, filled bool)   {}

func TestInteractiveDisplay(t *testing.T) {
	mockDisplay := NewMockDisplay()
	interactive := NewInteractiveDisplay(mockDisplay)
	defer interactive.Close()

	// Test initial state
	if interactive.currentMenu == nil {
		t.Error("Expected current menu to be set")
	}

	if interactive.selectedIndex != 0 {
		t.Error("Expected initial selected index to be 0")
	}

	// Test navigation
	interactive.NavigateDown()
	if interactive.selectedIndex != 1 {
		t.Error("Expected selected index to be 1 after navigating down")
	}

	interactive.NavigateUp()
	if interactive.selectedIndex != 0 {
		t.Error("Expected selected index to be 0 after navigating up")
	}

	// Test menu selection
	err := interactive.Select()
	if err != nil {
		t.Errorf("Expected no error on selection, got: %v", err)
	}
}

func TestNotificationManager(t *testing.T) {
	mockDisplay := NewMockDisplay()
	nm := NewNotificationManager(mockDisplay)

	// Test adding notifications
	nm.ShowInfo("Test info message")
	if nm.GetNotificationCount() != 1 {
		t.Error("Expected 1 notification after adding info")
	}

	nm.ShowWarning("Test warning message")
	if nm.GetNotificationCount() != 2 {
		t.Error("Expected 2 notifications after adding warning")
	}

	// Test priority ordering
	nm.ShowError("Test error message")
	if nm.GetNotificationCount() != 3 {
		t.Error("Expected 3 notifications after adding error")
	}

	// Test critical notification
	nm.ShowCriticalError("Critical error")
	if !nm.HasCriticalNotifications() {
		t.Error("Expected to have critical notifications")
	}

	// Test progress notification
	nm.ShowProgress("Test progress", 0.5)
	count := nm.GetNotificationCount()
	if count < 1 {
		t.Error("Expected at least 1 notification after adding progress")
	}

	// Test dismiss
	nm.DismissCurrentNotification()
	newCount := nm.GetNotificationCount()
	if newCount >= count {
		t.Error("Expected notification count to decrease after dismiss")
	}

	// Test dismiss all
	nm.DismissAll()
	if nm.GetNotificationCount() != 0 {
		t.Error("Expected 0 notifications after dismiss all")
	}
}

func TestPerformanceMonitor(t *testing.T) {
	mockDisplay := NewMockDisplay()
	pm := NewPerformanceMonitor(mockDisplay)

	// Test metric simulation
	metrics := pm.SimulateMetrics()
	if metrics.CPUUsage < 0 || metrics.CPUUsage > 100 {
		t.Error("Expected CPU usage to be between 0 and 100")
	}

	if metrics.MemoryUsage < 0 || metrics.MemoryUsage > 100 {
		t.Error("Expected memory usage to be between 0 and 100")
	}

	if len(metrics.MotorLoad) == 0 {
		t.Error("Expected motor load data to be populated")
	}

	if len(metrics.Temperature) == 0 {
		t.Error("Expected temperature data to be populated")
	}

	// Test metric updates
	pm.UpdateMetrics(metrics)
	currentMetrics := pm.GetCurrentMetrics()
	if currentMetrics.CPUUsage != metrics.CPUUsage {
		t.Error("Expected CPU usage to match updated metrics")
	}

	// Test display methods (should not panic)
	pm.DisplayCPUChart()
	pm.DisplayMemoryChart()
	pm.DisplayUpdateRateChart()
	pm.DisplayMotorLoadChart()
	pm.DisplaySystemOverview()
}

func TestThemeManager(t *testing.T) {
	mockDisplay := NewMockDisplay()
	tm := NewThemeManager(mockDisplay)

	// Test default theme
	currentTheme := tm.GetCurrentTheme()
	if currentTheme.Name != "Default" {
		t.Error("Expected default theme to be 'Default'")
	}

	// Test theme switching
	err := tm.SetTheme("High Contrast")
	if err != nil {
		t.Errorf("Expected no error switching theme, got: %v", err)
	}

	newTheme := tm.GetCurrentTheme()
	if newTheme.Name != "High Contrast" {
		t.Error("Expected current theme to be 'High Contrast'")
	}

	// Test invalid theme
	err = tm.SetTheme("NonExistent")
	if err == nil {
		t.Error("Expected error when setting non-existent theme")
	}

	// Test available themes
	themes := tm.GetAvailableThemes()
	if len(themes) < 3 {
		t.Error("Expected at least 3 default themes")
	}

	// Test custom theme registration
	customTheme := ThemeConfig{
		Name:            "Custom",
		BackgroundColor: true,
		TextColor:       false,
		BorderStyle:     BorderRounded,
	}
	tm.RegisterTheme(customTheme)

	err = tm.SetTheme("Custom")
	if err != nil {
		t.Errorf("Expected no error setting custom theme, got: %v", err)
	}

	// Test themed drawing (should not panic)
	tm.DrawThemedText(0, 0, "Test text")
	tm.DrawThemedBorder(0, 0, 50, 30)
}

func TestCustomLayoutManager(t *testing.T) {
	mockDisplay := NewMockDisplay()
	tm := NewThemeManager(mockDisplay)
	clm := NewCustomLayoutManager(mockDisplay, tm)

	// Test region definition
	region := LayoutRegion{
		X: 0, Y: 0, Width: 50, Height: 20,
		Padding: 2, BorderStyle: BorderSingle,
		Content: "Test content", Visible: true, ZOrder: 1,
	}
	clm.DefineRegion("test", region)

	// Test region update
	clm.UpdateRegion("test", "Updated content")

	// Test visibility
	clm.ShowRegion("test")
	clm.HideRegion("test")

	// Test dashboard creation
	clm.CreateDashboard()

	// Test rendering (should not panic)
	clm.Render()
}

func TestAnimations(t *testing.T) {
	mockDisplay := NewMockDisplay()
	tm := NewThemeManager(mockDisplay)

	// Test animation start
	tm.StartAnimation("test", AnimationFadeIn, 1*time.Second)
	progress := tm.GetAnimationProgress("test")
	if progress != 0.0 {
		t.Error("Expected initial animation progress to be 0")
	}

	// Test animation update
	time.Sleep(100 * time.Millisecond)
	tm.UpdateAnimations()
	progress = tm.GetAnimationProgress("test")
	if progress <= 0.0 {
		t.Error("Expected animation progress to increase")
	}

	// Test easing functions
	easingFunctions := []EasingFunction{
		EaseInOutCubic,
		EaseInQuad,
		EaseOutQuad,
		EaseInOutSine,
		EaseBounce,
	}

	for _, easing := range easingFunctions {
		result := easing(0.5)
		if result < 0 || result > 1 {
			t.Errorf("Expected easing result to be between 0 and 1, got %f", result)
		}
	}
}

func TestNotificationExpiry(t *testing.T) {
	mockDisplay := NewMockDisplay()
	nm := NewNotificationManager(mockDisplay)

	// Add notification with short duration
	notification := Notification{
		Message:  "Short lived",
		Type:     NotificationInfo,
		Priority: PriorityLow,
		Duration: 10 * time.Millisecond,
	}
	nm.AddNotification(notification)

	if nm.GetNotificationCount() != 1 {
		t.Error("Expected 1 notification initially")
	}

	// Wait for expiry
	time.Sleep(20 * time.Millisecond)
	nm.Update()

	if nm.GetNotificationCount() != 0 {
		t.Error("Expected notification to expire")
	}
}

func TestPerformanceHistoryLimit(t *testing.T) {
	mockDisplay := NewMockDisplay()
	pm := NewPerformanceMonitor(mockDisplay)

	// Add more data points than the limit
	for i := 0; i < pm.maxHistory+10; i++ {
		metrics := SystemMetrics{
			CPUUsage:    float64(i),
			MemoryUsage: float64(i),
		}
		pm.UpdateMetrics(metrics)
	}

	// Check that history is limited
	if len(pm.cpuHistory) > pm.maxHistory {
		t.Errorf("Expected CPU history to be limited to %d, got %d", pm.maxHistory, len(pm.cpuHistory))
	}

	if len(pm.memoryHistory) > pm.maxHistory {
		t.Errorf("Expected memory history to be limited to %d, got %d", pm.maxHistory, len(pm.memoryHistory))
	}
}

func BenchmarkNotificationUpdate(b *testing.B) {
	mockDisplay := NewMockDisplay()
	nm := NewNotificationManager(mockDisplay)

	// Add some notifications
	for i := 0; i < 10; i++ {
		nm.ShowInfo("Benchmark notification")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nm.Update()
	}
}

func BenchmarkPerformanceMonitorUpdate(b *testing.B) {
	mockDisplay := NewMockDisplay()
	pm := NewPerformanceMonitor(mockDisplay)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics := pm.SimulateMetrics()
		pm.UpdateMetrics(metrics)
	}
}

func BenchmarkThemeManagerDrawing(b *testing.B) {
	mockDisplay := NewMockDisplay()
	tm := NewThemeManager(mockDisplay)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tm.DrawThemedText(0, 0, "Benchmark text")
		tm.DrawThemedBorder(0, 0, 50, 30)
	}
}
