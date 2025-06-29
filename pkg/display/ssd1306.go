package display

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// SSD1306 represents an SSD1306 OLED display
type SSD1306 struct {
	Width      int
	Height     int
	I2CAddress uint8
	I2CBus     int
	buffer     *image.RGBA
	mock       bool
	verbose    bool
	lastUpdate time.Time
}

// DisplayConfig holds the configuration for the SSD1306 display
type DisplayConfig struct {
	Width      int    // Display width in pixels (128, 64)
	Height     int    // Display height in pixels (64, 32)
	I2CAddress uint8  // I2C address (0x3C or 0x3D)
	I2CBus     int    // I2C bus number (usually 1)
	Mock       bool   // Use mock mode for testing
	Verbose    bool   // Enable verbose logging
}

// NewSSD1306 creates a new SSD1306 display instance
func NewSSD1306(config DisplayConfig) (*SSD1306, error) {
	display := &SSD1306{
		Width:      config.Width,
		Height:     config.Height,
		I2CAddress: config.I2CAddress,
		I2CBus:     config.I2CBus,
		mock:       config.Mock,
		verbose:    config.Verbose,
		buffer:     image.NewRGBA(image.Rect(0, 0, config.Width, config.Height)),
		lastUpdate: time.Now(),
	}

	if err := display.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize SSD1306: %w", err)
	}

	return display, nil
}

// Initialize sets up the SSD1306 display
func (d *SSD1306) Initialize() error {
	if d.mock {
		if d.verbose {
			log.Printf("[DISPLAY] Mock SSD1306 initialized: %dx%d @ I2C 0x%02X", 
				d.Width, d.Height, d.I2CAddress)
		}
		return nil
	}

	// In real implementation, this would initialize I2C communication
	// and send initialization commands to the SSD1306
	if d.verbose {
		log.Printf("[DISPLAY] SSD1306 initialized: %dx%d @ I2C bus %d, addr 0x%02X", 
			d.Width, d.Height, d.I2CBus, d.I2CAddress)
	}

	return nil
}

// Clear clears the display buffer
func (d *SSD1306) Clear() {
	// Fill buffer with black
	draw.Draw(d.buffer, d.buffer.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)
}

// SetPixel sets a pixel at the given coordinates
func (d *SSD1306) SetPixel(x, y int, on bool) {
	if x < 0 || x >= d.Width || y < 0 || y >= d.Height {
		return
	}

	var c color.RGBA
	if on {
		c = color.RGBA{255, 255, 255, 255} // White
	} else {
		c = color.RGBA{0, 0, 0, 255} // Black
	}

	d.buffer.Set(x, y, c)
}

// DrawText draws text at the specified position
func (d *SSD1306) DrawText(text string, x, y int, fontFace font.Face) {
	if fontFace == nil {
		fontFace = basicfont.Face7x13 // Default font
	}

	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	drawer := &font.Drawer{
		Dst:  d.buffer,
		Src:  image.NewUniform(color.RGBA{255, 255, 255, 255}),
		Face: fontFace,
		Dot:  point,
	}

	drawer.DrawString(text)
}

// DrawLine draws a line between two points
func (d *SSD1306) DrawLine(x1, y1, x2, y2 int) {
	// Simple line drawing using Bresenham's algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 >= x2 {
		sx = -1
	}
	sy := 1
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	x, y := x1, y1
	for {
		d.SetPixel(x, y, true)
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// DrawRectangle draws a rectangle
func (d *SSD1306) DrawRectangle(x, y, width, height int, filled bool) {
	if filled {
		for i := x; i < x+width; i++ {
			for j := y; j < y+height; j++ {
				d.SetPixel(i, j, true)
			}
		}
	} else {
		// Draw outline
		d.DrawLine(x, y, x+width-1, y)                 // Top
		d.DrawLine(x, y+height-1, x+width-1, y+height-1) // Bottom
		d.DrawLine(x, y, x, y+height-1)                 // Left
		d.DrawLine(x+width-1, y, x+width-1, y+height-1) // Right
	}
}

// DrawProgressBar draws a progress bar
func (d *SSD1306) DrawProgressBar(x, y, width, height int, progress float64) {
	// Draw border
	d.DrawRectangle(x, y, width, height, false)
	
	// Draw filled portion
	fillWidth := int(float64(width-2) * progress)
	if fillWidth > 0 {
		d.DrawRectangle(x+1, y+1, fillWidth, height-2, true)
	}
}

// Update sends the buffer to the display
func (d *SSD1306) Update() error {
	if d.mock {
		if d.verbose {
			elapsed := time.Since(d.lastUpdate)
			log.Printf("[DISPLAY] Mock update - %d ms since last update", elapsed.Milliseconds())
		}
		d.lastUpdate = time.Now()
		return nil
	}

	// In real implementation, this would send the buffer data via I2C
	if d.verbose {
		elapsed := time.Since(d.lastUpdate)
		log.Printf("[DISPLAY] SSD1306 updated via I2C - %d ms since last update", elapsed.Milliseconds())
	}
	d.lastUpdate = time.Now()
	
	return nil
}

// Close closes the display connection
func (d *SSD1306) Close() error {
	if d.mock {
		if d.verbose {
			log.Println("[DISPLAY] Mock SSD1306 closed")
		}
		return nil
	}

	// In real implementation, this would close I2C connection
	if d.verbose {
		log.Println("[DISPLAY] SSD1306 I2C connection closed")
	}
	
	return nil
}

// GetBuffer returns the current display buffer for advanced operations
func (d *SSD1306) GetBuffer() *image.RGBA {
	return d.buffer
}

// SetBrightness sets the display brightness (0-255)
func (d *SSD1306) SetBrightness(brightness uint8) error {
	if d.mock {
		if d.verbose {
			log.Printf("[DISPLAY] Mock brightness set to %d", brightness)
		}
		return nil
	}

	// In real implementation, this would send brightness command via I2C
	if d.verbose {
		log.Printf("[DISPLAY] SSD1306 brightness set to %d", brightness)
	}
	
	return nil
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DisplayInterface defines the interface for display modules
type DisplayInterface interface {
	Initialize() error
	Clear()
	SetPixel(x, y int, on bool)
	DrawText(text string, x, y int, fontFace font.Face)
	DrawLine(x1, y1, x2, y2 int)
	DrawRectangle(x, y, width, height int, filled bool)
	DrawProgressBar(x, y, width, height int, progress float64)
	Update() error
	Close() error
}
