package display

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/image/font/basicfont"
)

// InteractiveDisplay provides menu-driven display functionality
type InteractiveDisplay struct {
	display       DisplayInterface
	menuStack     []Menu
	currentMenu   *Menu
	selectedIndex int
	lastUpdate    time.Time
	animations    map[string]*Animation
}

// Menu represents a display menu
type Menu struct {
	Title     string
	Items     []MenuItem
	Parent    *Menu
	OnSelect  func(item MenuItem) error
}

// MenuItem represents a menu item
type MenuItem struct {
	Text     string
	Value    interface{}
	Type     MenuItemType
	SubMenu  *Menu
	Action   func() error
}

// MenuItemType defines different menu item types
type MenuItemType int

const (
	MenuItemAction MenuItemType = iota
	MenuItemSubmenu
	MenuItemValue
	MenuItemToggle
	MenuItemSlider
)

// Animation represents a display animation
type Animation struct {
	StartTime time.Time
	Duration  time.Duration
	Type      AnimationType
	Progress  float64
	Data      interface{}
}

// AnimationType defines animation types
type AnimationType int

const (
	AnimationFadeIn AnimationType = iota
	AnimationFadeOut
	AnimationSlideLeft
	AnimationSlideRight
	AnimationPulse
	AnimationRotate
)

// NewInteractiveDisplay creates a new interactive display
func NewInteractiveDisplay(display DisplayInterface) *InteractiveDisplay {
	id := &InteractiveDisplay{
		display:     display,
		menuStack:   make([]Menu, 0),
		animations:  make(map[string]*Animation),
		lastUpdate:  time.Now(),
	}
	
	// Create main menu
	mainMenu := id.createMainMenu()
	id.setMenu(mainMenu)
	
	return id
}

// createMainMenu creates the main navigation menu
func (id *InteractiveDisplay) createMainMenu() Menu {
	return Menu{
		Title: "Robot Control",
		Items: []MenuItem{
			{Text: "Robot Status", Type: MenuItemSubmenu, SubMenu: id.createStatusMenu()},
			{Text: "Joint Control", Type: MenuItemSubmenu, SubMenu: id.createJointMenu()},
			{Text: "Workspace Map", Type: MenuItemAction, Action: id.showWorkspaceMap},
			{Text: "Performance", Type: MenuItemSubmenu, SubMenu: id.createPerformanceMenu()},
			{Text: "System Health", Type: MenuItemAction, Action: id.showSystemHealth},
			{Text: "Settings", Type: MenuItemSubmenu, SubMenu: id.createSettingsMenu()},
			{Text: "Diagnostics", Type: MenuItemAction, Action: id.showDiagnostics},
		},
	}
}

// createStatusMenu creates the robot status submenu
func (id *InteractiveDisplay) createStatusMenu() *Menu {
	return &Menu{
		Title: "Robot Status",
		Items: []MenuItem{
			{Text: "Joint Angles", Type: MenuItemAction, Action: id.showJointAngles},
			{Text: "Motor Status", Type: MenuItemAction, Action: id.showMotorStatus},
			{Text: "Position Info", Type: MenuItemAction, Action: id.showPositionInfo},
			{Text: "Error Log", Type: MenuItemAction, Action: id.showErrorLog},
			{Text: "Back", Type: MenuItemAction, Action: id.goBack},
		},
	}
}

// createJointMenu creates the joint control submenu
func (id *InteractiveDisplay) createJointMenu() *Menu {
	return &Menu{
		Title: "Joint Control",
		Items: []MenuItem{
			{Text: "Home All", Type: MenuItemAction, Action: id.homeAllJoints},
			{Text: "Enable All", Type: MenuItemAction, Action: id.enableAllJoints},
			{Text: "Disable All", Type: MenuItemAction, Action: id.disableAllJoints},
			{Text: "Emergency Stop", Type: MenuItemAction, Action: id.emergencyStop},
			{Text: "Back", Type: MenuItemAction, Action: id.goBack},
		},
	}
}

// createPerformanceMenu creates the performance monitoring submenu
func (id *InteractiveDisplay) createPerformanceMenu() *Menu {
	return &Menu{
		Title: "Performance",
		Items: []MenuItem{
			{Text: "CPU Usage", Type: MenuItemAction, Action: id.showCPUChart},
			{Text: "Memory Usage", Type: MenuItemAction, Action: id.showMemoryChart},
			{Text: "Update Rate", Type: MenuItemAction, Action: id.showUpdateRate},
			{Text: "Motor Load", Type: MenuItemAction, Action: id.showMotorLoad},
			{Text: "Back", Type: MenuItemAction, Action: id.goBack},
		},
	}
}

// createSettingsMenu creates the settings submenu
func (id *InteractiveDisplay) createSettingsMenu() *Menu {
	return &Menu{
		Title: "Settings",
		Items: []MenuItem{
			{Text: "Auto Screen", Type: MenuItemToggle, Value: true},
			{Text: "Update Rate", Type: MenuItemSlider, Value: 200},
			{Text: "Brightness", Type: MenuItemSlider, Value: 100},
			{Text: "Animations", Type: MenuItemToggle, Value: true},
			{Text: "Back", Type: MenuItemAction, Action: id.goBack},
		},
	}
}

// setMenu sets the current menu
func (id *InteractiveDisplay) setMenu(menu Menu) {
	id.currentMenu = &menu
	id.selectedIndex = 0
	id.render()
}

// navigateUp moves selection up in current menu
func (id *InteractiveDisplay) NavigateUp() {
	if id.currentMenu != nil && len(id.currentMenu.Items) > 0 {
		id.selectedIndex = (id.selectedIndex - 1 + len(id.currentMenu.Items)) % len(id.currentMenu.Items)
		id.render()
	}
}

// navigateDown moves selection down in current menu
func (id *InteractiveDisplay) NavigateDown() {
	if id.currentMenu != nil && len(id.currentMenu.Items) > 0 {
		id.selectedIndex = (id.selectedIndex + 1) % len(id.currentMenu.Items)
		id.render()
	}
}

// select activates the currently selected menu item
func (id *InteractiveDisplay) Select() error {
	if id.currentMenu == nil || id.selectedIndex >= len(id.currentMenu.Items) {
		return fmt.Errorf("no valid selection")
	}
	
	item := id.currentMenu.Items[id.selectedIndex]
	
	switch item.Type {
	case MenuItemAction:
		if item.Action != nil {
			return item.Action()
		}
	case MenuItemSubmenu:
		if item.SubMenu != nil {
			id.menuStack = append(id.menuStack, *id.currentMenu)
			item.SubMenu.Parent = id.currentMenu
			id.setMenu(*item.SubMenu)
		}
	case MenuItemToggle:
		// Toggle boolean value
		if val, ok := item.Value.(bool); ok {
			item.Value = !val
			id.render()
		}
	case MenuItemSlider:
		// For now, just show current value
		id.showValue(item.Text, fmt.Sprintf("%v", item.Value))
	}
	
	return nil
}

// goBack returns to the parent menu
func (id *InteractiveDisplay) goBack() error {
	if len(id.menuStack) > 0 {
		// Pop from stack
		parentMenu := id.menuStack[len(id.menuStack)-1]
		id.menuStack = id.menuStack[:len(id.menuStack)-1]
		id.setMenu(parentMenu)
	}
	return nil
}

// render draws the current menu
func (id *InteractiveDisplay) render() {
	if id.currentMenu == nil {
		return
	}
	
	id.display.Clear()
	
	// Draw title with underline
	titleWidth := len(id.currentMenu.Title) * 6
	titleX := (128 - titleWidth) / 2
	id.display.DrawText(id.currentMenu.Title, titleX, 0, basicfont.Face7x13)
	id.display.DrawLine(titleX, 10, titleX+titleWidth, 10)
	
	// Draw menu items
	startY := 16
	visibleItems := 4 // Maximum items visible at once
	startIndex := 0
	
	// Adjust start index if selected item is off-screen
	if id.selectedIndex >= visibleItems {
		startIndex = id.selectedIndex - visibleItems + 1
	}
	
	for i := 0; i < visibleItems && (startIndex+i) < len(id.currentMenu.Items); i++ {
		itemIndex := startIndex + i
		item := id.currentMenu.Items[itemIndex]
		y := startY + i*12
		
		// Highlight selected item
		if itemIndex == id.selectedIndex {
			// Create highlight with simple indicators since we don't have FillRect
			id.display.DrawText("> ", 0, y, basicfont.Face7x13)
		}
		
		// Draw item text
		text := item.Text
		if item.Type == MenuItemToggle {
			if val, ok := item.Value.(bool); ok {
				if val {
					text += " [ON]"
				} else {
					text += " [OFF]"
				}
			}
		} else if item.Type == MenuItemSlider {
			text += fmt.Sprintf(" [%v]", item.Value)
		} else if item.Type == MenuItemSubmenu {
			text += " >"
		}
		
		id.display.DrawText(text, 12, y, basicfont.Face7x13)
	}
	
	// Draw scroll indicators if needed
	if startIndex > 0 {
		id.display.DrawText("^", 120, 16, basicfont.Face7x13)
	}
	if startIndex+visibleItems < len(id.currentMenu.Items) {
		id.display.DrawText("v", 120, 52, basicfont.Face7x13)
	}
	
	id.display.Update()
	id.lastUpdate = time.Now()
}

// Action implementations

func (id *InteractiveDisplay) showWorkspaceMap() error {
	id.display.Clear()
	id.display.DrawText("Workspace Map", 0, 0, basicfont.Face7x13)
	
	// Draw circular workspace boundary
	centerX, centerY := 64, 32
	radius := 25
	
	// Draw boundary circle
	for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
		x := centerX + int(float64(radius)*math.Cos(angle))
		y := centerY + int(float64(radius)*math.Sin(angle))
		if x >= 0 && x < 128 && y >= 0 && y < 64 {
			id.display.SetPixel(x, y, true)
		}
	}
	
	// Draw current position (center for demo)
	id.display.DrawRectangle(centerX-2, centerY-2, 4, 4, true)
	
	// Draw safety zones
	safeRadius := 15
	for angle := 0.0; angle < 2*math.Pi; angle += 0.2 {
		x := centerX + int(float64(safeRadius)*math.Cos(angle))
		y := centerY + int(float64(safeRadius)*math.Sin(angle))
		if x >= 0 && x < 128 && y >= 0 && y < 64 {
			id.display.SetPixel(x, y, true)
		}
	}
	
	id.display.DrawText("Press any key to return", 0, 56, basicfont.Face7x13)
	id.display.Update()
	
	time.Sleep(5 * time.Second)
	id.render()
	return nil
}

func (id *InteractiveDisplay) showSystemHealth() error {
	id.display.Clear()
	id.display.DrawText("System Health", 0, 0, basicfont.Face7x13)
	
	// Health indicators
	healthItems := []struct {
		name   string
		status string
		ok     bool
	}{
		{"CPU Temp", "45°C", true},
		{"Motor 1", "OK", true},
		{"Motor 2", "OK", true},
		{"Motor 3", "WARN", false},
		{"I2C Bus", "OK", true},
		{"Power", "12.1V", true},
	}
	
	for i, item := range healthItems {
		y := 12 + i*8
		if y > 56 {
			break
		}
		
		// Status indicator
		if item.ok {
			id.display.DrawRectangle(0, y, 4, 4, true)
		} else {
			id.display.DrawRectangle(0, y, 4, 4, false)
		}
		
		// Item name and status
		id.display.DrawText(fmt.Sprintf("%s: %s", item.name, item.status), 8, y, basicfont.Face7x13)
	}
	
	time.Sleep(5 * time.Second)
	id.render()
	return nil
}

func (id *InteractiveDisplay) showCPUChart() error {
	id.display.Clear()
	id.display.DrawText("CPU Usage", 0, 0, basicfont.Face7x13)
	
	// Simulate CPU usage data
	data := []int{30, 45, 38, 52, 41, 35, 48, 44, 39, 42, 36, 50}
	maxVal := 100
	chartHeight := 40
	chartY := 20
	
	// Draw chart
	for i, val := range data {
		if i >= 120 { // Limit to display width
			break
		}
		
		barHeight := (val * chartHeight) / maxVal
		x := i * 10
		y := chartY + chartHeight - barHeight
		
		id.display.DrawLine(x, y, x, chartY+chartHeight)
	}
	
	// Draw scale
	id.display.DrawText("0%", 0, 56, basicfont.Face7x13)
	id.display.DrawText("100%", 110, 56, basicfont.Face7x13)
	
	time.Sleep(5 * time.Second)
	id.render()
	return nil
}

func (id *InteractiveDisplay) showMemoryChart() error {
	id.display.Clear()
	id.display.DrawText("Memory Usage", 0, 0, basicfont.Face7x13)
	
	// Memory usage pie chart
	centerX, centerY := 64, 35
	radius := 20
	
	// Used memory (60%)
	usedAngle := 2 * math.Pi * 0.6
	for angle := 0.0; angle < usedAngle; angle += 0.1 {
		for r := 0; r < radius; r++ {
			x := centerX + int(float64(r)*math.Cos(angle-math.Pi/2))
			y := centerY + int(float64(r)*math.Sin(angle-math.Pi/2))
			if x >= 0 && x < 128 && y >= 0 && y < 64 {
				id.display.SetPixel(x, y, true)
			}
		}
	}
	
	// Labels
	id.display.DrawText("Used: 60%", 0, 56, basicfont.Face7x13)
	id.display.DrawText("Free: 40%", 70, 56, basicfont.Face7x13)
	
	time.Sleep(5 * time.Second)
	id.render()
	return nil
}

func (id *InteractiveDisplay) showUpdateRate() error {
	id.showValue("Update Rate", "200ms")
	return nil
}

func (id *InteractiveDisplay) showMotorLoad() error {
	id.display.Clear()
	id.display.DrawText("Motor Load", 0, 0, basicfont.Face7x13)
	
	// Motor load bars
	motors := []struct {
		name string
		load int
	}{
		{"Pan", 45},
		{"Tilt", 32},
		{"Elbow", 67},
		{"Wrist1", 23},
		{"Wrist2", 41},
		{"Gripper", 15},
	}
	
	for i, motor := range motors {
		y := 12 + i*8
		if y > 56 {
			break
		}
		
		// Motor name
		id.display.DrawText(motor.name, 0, y, basicfont.Face7x13)
		
		// Load bar
		barWidth := 60
		loadWidth := (motor.load * barWidth) / 100
		id.display.DrawRectangle(50, y, barWidth, 6, false)
		id.display.DrawRectangle(50, y, loadWidth, 6, true)
		
		// Percentage
		id.display.DrawText(fmt.Sprintf("%d%%", motor.load), 115, y, basicfont.Face7x13)
	}
	
	time.Sleep(5 * time.Second)
	id.render()
	return nil
}

func (id *InteractiveDisplay) showJointAngles() error {
	id.showValue("Joint Angles", "Feature Active")
	return nil
}

func (id *InteractiveDisplay) showMotorStatus() error {
	id.showValue("Motor Status", "All Enabled")
	return nil
}

func (id *InteractiveDisplay) showPositionInfo() error {
	id.showValue("Position", "X:0 Y:0 Z:200")
	return nil
}

func (id *InteractiveDisplay) showErrorLog() error {
	id.showValue("Error Log", "No Errors")
	return nil
}

func (id *InteractiveDisplay) showDiagnostics() error {
	id.showValue("Diagnostics", "All Systems OK")
	return nil
}

func (id *InteractiveDisplay) homeAllJoints() error {
	id.showValue("Home All", "Executing...")
	return nil
}

func (id *InteractiveDisplay) enableAllJoints() error {
	id.showValue("Enable All", "Motors Enabled")
	return nil
}

func (id *InteractiveDisplay) disableAllJoints() error {
	id.showValue("Disable All", "Motors Disabled")
	return nil
}

func (id *InteractiveDisplay) emergencyStop() error {
	id.display.Clear()
	id.display.DrawText("EMERGENCY", 10, 20, basicfont.Face7x13)
	id.display.DrawText("STOP", 30, 35, basicfont.Face7x13)
	id.display.Update()
	time.Sleep(3 * time.Second)
	id.render()
	return nil
}

func (id *InteractiveDisplay) showValue(title, value string) {
	id.display.Clear()
	id.display.DrawText(title, 0, 0, basicfont.Face7x13)
	id.display.DrawText(value, 0, 25, basicfont.Face7x13)
	id.display.DrawText("Press any key to return", 0, 56, basicfont.Face7x13)
	id.display.Update()
	time.Sleep(3 * time.Second)
	id.render()
}

// StartAnimation starts a new animation
func (id *InteractiveDisplay) StartAnimation(name string, animType AnimationType, duration time.Duration) {
	id.animations[name] = &Animation{
		StartTime: time.Now(),
		Duration:  duration,
		Type:      animType,
		Progress:  0.0,
	}
}

// UpdateAnimations updates all active animations
func (id *InteractiveDisplay) UpdateAnimations() {
	now := time.Now()
	for name, anim := range id.animations {
		elapsed := now.Sub(anim.StartTime)
		anim.Progress = math.Min(1.0, elapsed.Seconds()/anim.Duration.Seconds())
		
		if anim.Progress >= 1.0 {
			delete(id.animations, name)
		}
	}
}

// Close cleans up the interactive display
func (id *InteractiveDisplay) Close() {
	// Clean up any resources
}
