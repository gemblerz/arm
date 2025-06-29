package display

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/image/font/basicfont"
)

// ThemeConfig defines display theme configuration
type ThemeConfig struct {
	Name            string
	BackgroundColor bool
	TextColor       bool
	AccentColor     bool
	BorderStyle     BorderStyle
	AnimationSpeed  time.Duration
}

// BorderStyle defines different border styles
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderSingle
	BorderDouble
	BorderRounded
	BorderDashed
)

// ThemeManager manages display themes and visual customizations
type ThemeManager struct {
	display      DisplayInterface
	currentTheme ThemeConfig
	themes       map[string]ThemeConfig
	animations   map[string]*ThemeAnimation
}

// ThemeAnimation represents a visual animation
type ThemeAnimation struct {
	StartTime time.Time
	Duration  time.Duration
	Type      AnimationType
	Progress  float64
	Easing    EasingFunction
}

// EasingFunction defines animation easing
type EasingFunction func(float64) float64

// Predefined themes
var (
	DefaultTheme = ThemeConfig{
		Name:            "Default",
		BackgroundColor: false,
		TextColor:       true,
		AccentColor:     true,
		BorderStyle:     BorderSingle,
		AnimationSpeed:  300 * time.Millisecond,
	}
	
	HighContrastTheme = ThemeConfig{
		Name:            "High Contrast",
		BackgroundColor: false,
		TextColor:       true,
		AccentColor:     true,
		BorderStyle:     BorderDouble,
		AnimationSpeed:  200 * time.Millisecond,
	}
	
	MinimalTheme = ThemeConfig{
		Name:            "Minimal",
		BackgroundColor: false,
		TextColor:       true,
		AccentColor:     false,
		BorderStyle:     BorderNone,
		AnimationSpeed:  500 * time.Millisecond,
	}
)

// NewThemeManager creates a new theme manager
func NewThemeManager(display DisplayInterface) *ThemeManager {
	tm := &ThemeManager{
		display:    display,
		themes:     make(map[string]ThemeConfig),
		animations: make(map[string]*ThemeAnimation),
	}
	
	// Register default themes
	tm.RegisterTheme(DefaultTheme)
	tm.RegisterTheme(HighContrastTheme)
	tm.RegisterTheme(MinimalTheme)
	
	// Set default theme
	tm.SetTheme("Default")
	
	return tm
}

// RegisterTheme registers a new theme
func (tm *ThemeManager) RegisterTheme(theme ThemeConfig) {
	tm.themes[theme.Name] = theme
}

// SetTheme activates a theme
func (tm *ThemeManager) SetTheme(name string) error {
	if theme, exists := tm.themes[name]; exists {
		tm.currentTheme = theme
		return nil
	}
	return fmt.Errorf("theme '%s' not found", name)
}

// GetCurrentTheme returns the current theme
func (tm *ThemeManager) GetCurrentTheme() ThemeConfig {
	return tm.currentTheme
}

// GetAvailableThemes returns list of available themes
func (tm *ThemeManager) GetAvailableThemes() []string {
	themes := make([]string, 0, len(tm.themes))
	for name := range tm.themes {
		themes = append(themes, name)
	}
	return themes
}

// DrawThemedText draws text with current theme
func (tm *ThemeManager) DrawThemedText(x, y int, text string) {
	// Note: DisplayInterface doesn't support text colors, so we just draw text normally
	tm.display.DrawText(text, x, y, basicfont.Face7x13)
}

// DrawThemedBorder draws a border with current theme style
func (tm *ThemeManager) DrawThemedBorder(x, y, width, height int) {
	switch tm.currentTheme.BorderStyle {
	case BorderSingle:
		tm.display.DrawRectangle(x, y, width, height, false)
	case BorderDouble:
		tm.display.DrawRectangle(x, y, width, height, false)
		tm.display.DrawRectangle(x+1, y+1, width-2, height-2, false)
	case BorderRounded:
		tm.drawRoundedRect(x, y, width, height, 3)
	case BorderDashed:
		tm.drawDashedRect(x, y, width, height)
	case BorderNone:
		// No border
	}
}

// drawRoundedRect draws a rectangle with rounded corners
func (tm *ThemeManager) drawRoundedRect(x, y, width, height, radius int) {
	// Draw straight lines
	tm.display.DrawLine(x+radius, y, x+width-radius, y)                    // Top
	tm.display.DrawLine(x+radius, y+height, x+width-radius, y+height)      // Bottom
	tm.display.DrawLine(x, y+radius, x, y+height-radius)                   // Left
	tm.display.DrawLine(x+width, y+radius, x+width, y+height-radius)       // Right
	
	// Draw rounded corners (simplified)
	for i := 0; i < radius; i++ {
		for j := 0; j < radius; j++ {
			// Top-left corner
			if i*i+j*j >= (radius-1)*(radius-1) {
				tm.display.SetPixel(x+radius-i, y+radius-j, true)
			}
			// Top-right corner
			if i*i+j*j >= (radius-1)*(radius-1) {
				tm.display.SetPixel(x+width-radius+i, y+radius-j, true)
			}
			// Bottom-left corner
			if i*i+j*j >= (radius-1)*(radius-1) {
				tm.display.SetPixel(x+radius-i, y+height-radius+j, true)
			}
			// Bottom-right corner
			if i*i+j*j >= (radius-1)*(radius-1) {
				tm.display.SetPixel(x+width-radius+i, y+height-radius+j, true)
			}
		}
	}
}

// drawDashedRect draws a dashed rectangle
func (tm *ThemeManager) drawDashedRect(x, y, width, height int) {
	dashLength := 3
	spaceLength := 2
	
	// Top and bottom lines
	for i := 0; i < width; i += dashLength + spaceLength {
		endX := i + dashLength
		if endX > width {
			endX = width
		}
		tm.display.DrawLine(x+i, y, x+endX, y)           // Top
		tm.display.DrawLine(x+i, y+height, x+endX, y+height) // Bottom
	}
	
	// Left and right lines
	for i := 0; i < height; i += dashLength + spaceLength {
		endY := i + dashLength
		if endY > height {
			endY = height
		}
		tm.display.DrawLine(x, y+i, x, y+endY)           // Left
		tm.display.DrawLine(x+width, y+i, x+width, y+endY) // Right
	}
}

// StartAnimation starts a themed animation
func (tm *ThemeManager) StartAnimation(name string, animType AnimationType, duration time.Duration) {
	tm.animations[name] = &ThemeAnimation{
		StartTime: time.Now(),
		Duration:  duration,
		Type:      animType,
		Progress:  0.0,
		Easing:    EaseInOutCubic,
	}
}

// UpdateAnimations updates all active animations
func (tm *ThemeManager) UpdateAnimations() {
	now := time.Now()
	for name, anim := range tm.animations {
		elapsed := now.Sub(anim.StartTime)
		rawProgress := math.Min(1.0, elapsed.Seconds()/anim.Duration.Seconds())
		anim.Progress = anim.Easing(rawProgress)
		
		if rawProgress >= 1.0 {
			delete(tm.animations, name)
		}
	}
}

// GetAnimationProgress returns the progress of a named animation
func (tm *ThemeManager) GetAnimationProgress(name string) float64 {
	if anim, exists := tm.animations[name]; exists {
		return anim.Progress
	}
	return 0.0
}

// Easing functions

// EaseInOutCubic provides smooth acceleration and deceleration
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

// EaseInQuad provides gradual acceleration
func EaseInQuad(t float64) float64 {
	return t * t
}

// EaseOutQuad provides gradual deceleration
func EaseOutQuad(t float64) float64 {
	return 1 - (1-t)*(1-t)
}

// EaseInOutSine provides smooth sine-wave easing
func EaseInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1) / 2
}

// EaseBounce provides bouncing effect
func EaseBounce(t float64) float64 {
	if t < 1/2.75 {
		return 7.5625 * t * t
	} else if t < 2/2.75 {
		t -= 1.5 / 2.75
		return 7.5625*t*t + 0.75
	} else if t < 2.5/2.75 {
		t -= 2.25 / 2.75
		return 7.5625*t*t + 0.9375
	} else {
		t -= 2.625 / 2.75
		return 7.5625*t*t + 0.984375
	}
}

// CustomLayoutManager provides advanced layout capabilities
type CustomLayoutManager struct {
	display     DisplayInterface
	themeManager *ThemeManager
	regions     map[string]LayoutRegion
}

// LayoutRegion defines a screen region
type LayoutRegion struct {
	X, Y          int
	Width, Height int
	Padding       int
	BorderStyle   BorderStyle
	Content       string
	Visible       bool
	ZOrder        int
}

// NewCustomLayoutManager creates a new layout manager
func NewCustomLayoutManager(display DisplayInterface, themeManager *ThemeManager) *CustomLayoutManager {
	return &CustomLayoutManager{
		display:      display,
		themeManager: themeManager,
		regions:      make(map[string]LayoutRegion),
	}
}

// DefineRegion defines a new layout region
func (clm *CustomLayoutManager) DefineRegion(name string, region LayoutRegion) {
	clm.regions[name] = region
}

// UpdateRegion updates content in a region
func (clm *CustomLayoutManager) UpdateRegion(name, content string) {
	if region, exists := clm.regions[name]; exists {
		region.Content = content
		clm.regions[name] = region
	}
}

// ShowRegion makes a region visible
func (clm *CustomLayoutManager) ShowRegion(name string) {
	if region, exists := clm.regions[name]; exists {
		region.Visible = true
		clm.regions[name] = region
	}
}

// HideRegion makes a region invisible
func (clm *CustomLayoutManager) HideRegion(name string) {
	if region, exists := clm.regions[name]; exists {
		region.Visible = false
		clm.regions[name] = region
	}
}

// Render renders all visible regions
func (clm *CustomLayoutManager) Render() {
	// Sort regions by Z-order
	sortedRegions := make([]LayoutRegion, 0)
	for _, region := range clm.regions {
		if region.Visible {
			sortedRegions = append(sortedRegions, region)
		}
	}
	
	// Simple bubble sort by Z-order
	for i := 0; i < len(sortedRegions)-1; i++ {
		for j := 0; j < len(sortedRegions)-i-1; j++ {
			if sortedRegions[j].ZOrder > sortedRegions[j+1].ZOrder {
				sortedRegions[j], sortedRegions[j+1] = sortedRegions[j+1], sortedRegions[j]
			}
		}
	}
	
	// Render each region
	for _, region := range sortedRegions {
		clm.renderRegion(region)
	}
}

// renderRegion renders a single region
func (clm *CustomLayoutManager) renderRegion(region LayoutRegion) {
	// Draw border if specified
	if region.BorderStyle != BorderNone {
		oldBorderStyle := clm.themeManager.currentTheme.BorderStyle
		clm.themeManager.currentTheme.BorderStyle = region.BorderStyle
		clm.themeManager.DrawThemedBorder(region.X, region.Y, region.Width, region.Height)
		clm.themeManager.currentTheme.BorderStyle = oldBorderStyle
	}
	
	// Draw content with padding
	contentX := region.X + region.Padding
	contentY := region.Y + region.Padding
	maxWidth := region.Width - 2*region.Padding
	
	// Simple text wrapping
	clm.drawWrappedText(contentX, contentY, region.Content, maxWidth)
}

// drawWrappedText draws text with word wrapping within a region
func (clm *CustomLayoutManager) drawWrappedText(x, y int, text string, maxWidth int) {
	words := splitTextIntoWords(text)
	currentLine := ""
	lineY := y
	lineHeight := 12
	
	for _, word := range words {
		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word
		
		// Estimate width (6 pixels per character for basic font)
		if len(testLine)*6 <= maxWidth {
			currentLine = testLine
		} else {
			// Draw current line and start new line
			if currentLine != "" {
				clm.themeManager.DrawThemedText(x, lineY, currentLine)
				lineY += lineHeight
			}
			currentLine = word
		}
	}
	
	// Draw remaining text
	if currentLine != "" {
		clm.themeManager.DrawThemedText(x, lineY, currentLine)
	}
}

// splitTextIntoWords splits text into words for wrapping
func splitTextIntoWords(text string) []string {
	words := make([]string, 0)
	currentWord := ""
	
	for _, char := range text {
		if char == ' ' || char == '\n' || char == '\t' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	
	if currentWord != "" {
		words = append(words, currentWord)
	}
	
	return words
}

// CreateDashboard creates a predefined dashboard layout
func (clm *CustomLayoutManager) CreateDashboard() {
	// Header region
	clm.DefineRegion("header", LayoutRegion{
		X: 0, Y: 0, Width: 128, Height: 12,
		Padding: 1, BorderStyle: BorderSingle,
		Content: "Robot Status", Visible: true, ZOrder: 1,
	})
	
	// Main content region
	clm.DefineRegion("main", LayoutRegion{
		X: 0, Y: 14, Width: 96, Height: 36,
		Padding: 2, BorderStyle: BorderSingle,
		Content: "Main content area", Visible: true, ZOrder: 1,
	})
	
	// Sidebar region
	clm.DefineRegion("sidebar", LayoutRegion{
		X: 98, Y: 14, Width: 30, Height: 36,
		Padding: 1, BorderStyle: BorderSingle,
		Content: "Side info", Visible: true, ZOrder: 1,
	})
	
	// Footer region
	clm.DefineRegion("footer", LayoutRegion{
		X: 0, Y: 52, Width: 128, Height: 12,
		Padding: 1, BorderStyle: BorderNone,
		Content: "Ready", Visible: true, ZOrder: 1,
	})
}
