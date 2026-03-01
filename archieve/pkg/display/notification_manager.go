package display

import (
	"fmt"
	"time"

	"golang.org/x/image/font/basicfont"
)

// NotificationPriority defines the priority level of notifications
type NotificationPriority int

const (
	PriorityLow NotificationPriority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// NotificationType defines different types of notifications
type NotificationType int

const (
	NotificationInfo NotificationType = iota
	NotificationWarning
	NotificationError
	NotificationSuccess
)

// Notification represents a display notification
type Notification struct {
	ID       string
	Message  string
	Type     NotificationType
	Priority NotificationPriority
	Time     time.Time
	Duration time.Duration
	Sticky   bool
	Progress float64 // For progress notifications (0.0 to 1.0)
}

// NotificationManager handles display notifications and alerts
type NotificationManager struct {
	display       DisplayInterface
	notifications []Notification
	currentIndex  int
	lastUpdate    time.Time
	isActive      bool
	blinkState    bool
	blinkTimer    time.Time
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(display DisplayInterface) *NotificationManager {
	return &NotificationManager{
		display:       display,
		notifications: make([]Notification, 0),
		lastUpdate:    time.Now(),
		blinkTimer:    time.Now(),
	}
}

// AddNotification adds a new notification
func (nm *NotificationManager) AddNotification(notification Notification) {
	// Set default values
	if notification.Time.IsZero() {
		notification.Time = time.Now()
	}
	if notification.Duration == 0 {
		notification.Duration = nm.getDefaultDuration(notification.Priority)
	}
	if notification.ID == "" {
		notification.ID = fmt.Sprintf("notif_%d", time.Now().UnixNano())
	}

	// Insert based on priority (higher priority first)
	inserted := false
	for i, existing := range nm.notifications {
		if notification.Priority > existing.Priority {
			// Insert at position i
			nm.notifications = append(nm.notifications[:i], append([]Notification{notification}, nm.notifications[i:]...)...)
			inserted = true
			break
		}
	}
	
	if !inserted {
		nm.notifications = append(nm.notifications, notification)
	}
	
	// Limit total notifications
	if len(nm.notifications) > 20 {
		nm.notifications = nm.notifications[:20]
	}
}

// getDefaultDuration returns default duration based on priority
func (nm *NotificationManager) getDefaultDuration(priority NotificationPriority) time.Duration {
	switch priority {
	case PriorityLow:
		return 3 * time.Second
	case PriorityMedium:
		return 5 * time.Second
	case PriorityHigh:
		return 8 * time.Second
	case PriorityCritical:
		return 0 // Sticky by default
	default:
		return 5 * time.Second
	}
}

// Update processes and displays notifications
func (nm *NotificationManager) Update() {
	now := time.Now()
	
	// Remove expired notifications
	nm.removeExpiredNotifications(now)
	
	// Update blink state for critical notifications
	if now.Sub(nm.blinkTimer) > 500*time.Millisecond {
		nm.blinkState = !nm.blinkState
		nm.blinkTimer = now
	}
	
	// Display current notification
	if len(nm.notifications) > 0 {
		nm.displayCurrentNotification()
		nm.isActive = true
	} else {
		nm.isActive = false
	}
	
	nm.lastUpdate = now
}

// removeExpiredNotifications removes notifications that have expired
func (nm *NotificationManager) removeExpiredNotifications(now time.Time) {
	filtered := make([]Notification, 0)
	
	for _, notif := range nm.notifications {
		if notif.Sticky || notif.Duration == 0 || now.Sub(notif.Time) < notif.Duration {
			filtered = append(filtered, notif)
		}
	}
	
	// Adjust current index if needed
	if nm.currentIndex >= len(filtered) {
		nm.currentIndex = 0
	}
	
	nm.notifications = filtered
}

// displayCurrentNotification displays the current notification (simplified version)
func (nm *NotificationManager) displayCurrentNotification() {
	if len(nm.notifications) == 0 {
		return
	}
	
	notif := nm.notifications[nm.currentIndex]
	
	// For critical notifications, blink the display
	if notif.Priority == PriorityCritical && nm.blinkState {
		nm.display.Clear()
		nm.display.Update()
		return
	}
	
	nm.display.Clear()
	
	// Simple notification display - just show the message
	nm.display.DrawText(notif.Message, 0, 20, basicfont.Face7x13)
	
	// Show type
	typeText := "INFO"
	switch notif.Type {
	case NotificationWarning:
		typeText = "WARN"
	case NotificationError:
		typeText = "ERR"
	case NotificationSuccess:
		typeText = "OK"
	}
	nm.display.DrawText(typeText, 0, 0, basicfont.Face7x13)
	
	// Show progress if applicable
	if notif.Progress > 0 {
		nm.display.DrawProgressBar(0, 50, 100, 8, notif.Progress)
	}
	
	nm.display.Update()
}

// Simple helper methods for text wrapping
func splitText(text string) []string {
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

// NextNotification moves to the next notification
func (nm *NotificationManager) NextNotification() {
	if len(nm.notifications) > 1 {
		nm.currentIndex = (nm.currentIndex + 1) % len(nm.notifications)
	}
}

// PreviousNotification moves to the previous notification
func (nm *NotificationManager) PreviousNotification() {
	if len(nm.notifications) > 1 {
		nm.currentIndex = (nm.currentIndex - 1 + len(nm.notifications)) % len(nm.notifications)
	}
}

// DismissCurrentNotification removes the current notification
func (nm *NotificationManager) DismissCurrentNotification() {
	if len(nm.notifications) > 0 {
		// Remove current notification
		nm.notifications = append(nm.notifications[:nm.currentIndex], nm.notifications[nm.currentIndex+1:]...)
		
		// Adjust index
		if nm.currentIndex >= len(nm.notifications) {
			nm.currentIndex = 0
		}
	}
}

// DismissAll removes all notifications
func (nm *NotificationManager) DismissAll() {
	nm.notifications = make([]Notification, 0)
	nm.currentIndex = 0
}

// GetNotificationCount returns the number of active notifications
func (nm *NotificationManager) GetNotificationCount() int {
	return len(nm.notifications)
}

// HasCriticalNotifications returns true if there are critical notifications
func (nm *NotificationManager) HasCriticalNotifications() bool {
	for _, notif := range nm.notifications {
		if notif.Priority == PriorityCritical {
			return true
		}
	}
	return false
}

// IsActive returns true if the notification manager is currently displaying notifications
func (nm *NotificationManager) IsActive() bool {
	return nm.isActive
}

// Predefined notification creators

// ShowInfo shows an info notification
func (nm *NotificationManager) ShowInfo(message string) {
	nm.AddNotification(Notification{
		Message:  message,
		Type:     NotificationInfo,
		Priority: PriorityLow,
	})
}

// ShowWarning shows a warning notification
func (nm *NotificationManager) ShowWarning(message string) {
	nm.AddNotification(Notification{
		Message:  message,
		Type:     NotificationWarning,
		Priority: PriorityMedium,
	})
}

// ShowError shows an error notification
func (nm *NotificationManager) ShowError(message string) {
	nm.AddNotification(Notification{
		Message:  message,
		Type:     NotificationError,
		Priority: PriorityHigh,
	})
}

// ShowCriticalError shows a critical error notification
func (nm *NotificationManager) ShowCriticalError(message string) {
	nm.AddNotification(Notification{
		Message:  message,
		Type:     NotificationError,
		Priority: PriorityCritical,
		Sticky:   true,
	})
}

// ShowSuccess shows a success notification
func (nm *NotificationManager) ShowSuccess(message string) {
	nm.AddNotification(Notification{
		Message:  message,
		Type:     NotificationSuccess,
		Priority: PriorityLow,
	})
}

// ShowProgress shows a progress notification
func (nm *NotificationManager) ShowProgress(message string, progress float64) {
	// Update existing progress notification or create new one
	id := "progress_main"
	
	// Remove existing progress notification
	for i, notif := range nm.notifications {
		if notif.ID == id {
			nm.notifications = append(nm.notifications[:i], nm.notifications[i+1:]...)
			break
		}
	}
	
	nm.AddNotification(Notification{
		ID:       id,
		Message:  message,
		Type:     NotificationInfo,
		Priority: PriorityMedium,
		Progress: progress,
		Duration: 10 * time.Second, // Longer duration for progress
	})
}
