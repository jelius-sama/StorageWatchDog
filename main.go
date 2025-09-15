package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"
	"time"

	gomail "gopkg.in/gomail.v2"
)

type Threshold struct {
	Percent    int
	NotifiedAt time.Time
}

func getDiskUsage(path string) (usedPercent float64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}
	total := float64(stat.Blocks) * float64(stat.Bsize)
	free := float64(stat.Bavail) * float64(stat.Bsize)
	used := total - free
	if total == 0 {
		return 0, fmt.Errorf("total blocks is zero")
	}
	return (used / total) * 100.0, nil
}

func sendMail(smtpHost string, smtpPort int, username, password, from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, username, password)
	return d.DialAndSend(m)
}

func main() {
	// Flags
	maildir := flag.String("path", "/var/vmail", "Path to monitor")
	interval := flag.Duration("interval", 10*time.Minute, "Check interval")

	warn80 := flag.Int("warn80", 80, "Warn threshold (percent)")
	warn85 := flag.Int("warn85", 85, "Warn threshold (percent)")
	warn90 := flag.Int("warn90", 90, "Warn threshold (percent)")
	warn95 := flag.Int("warn95", 95, "Warn threshold (percent)")
	warn98 := flag.Int("warn98", 98, "Critical warn threshold (percent)")

	repeatCritical := flag.Duration("criticalRepeat", 30*time.Minute, "Repeat interval for >= critical usage")

	growthThreshold := flag.Float64("growth", 5, "Growth percentage to trigger critical alert")
	growthInterval := flag.Duration("growthInterval", 5*time.Hour, "Interval to check growth")

	// SMTP config
	smtpHost := flag.String("smtpHost", "smtp.example.com", "SMTP host")
	smtpPort := flag.Int("smtpPort", 587, "SMTP port")
	smtpUser := flag.String("smtpUser", "user@example.com", "SMTP username")
	smtpPass := flag.String("smtpPass", "password", "SMTP password")
	mailFrom := flag.String("from", "alert@example.com", "From email")
	mailTo := flag.String("to", "admin@example.com", "To email")

	flag.Parse()

	// Setup thresholds
	thresholds := []Threshold{
		{Percent: *warn80},
		{Percent: *warn85},
		{Percent: *warn90},
		{Percent: *warn95},
		{Percent: *warn98},
	}

	lastUsage := 0.0
	lastUsageTime := time.Now()

	for {
		usedPercent, err := getDiskUsage(*maildir)
		if err != nil {
			log.Printf("Error checking disk usage: %v", err)
			time.Sleep(*interval)
			continue
		}

		now := time.Now()
		for i := range thresholds {
			t := &thresholds[i]
			if int(usedPercent) >= t.Percent {
				// Handle critical repeat rule
				if t.Percent >= *warn98 {
					if now.Sub(t.NotifiedAt) >= *repeatCritical {
						subject := fmt.Sprintf("[CRITICAL] Maildir usage %0.2f%% >= %d%%", usedPercent, t.Percent)
						body := fmt.Sprintf("Maildir at %s is at %0.2f%% usage.", *maildir, usedPercent)
						err := sendMail(*smtpHost, *smtpPort, *smtpUser, *smtpPass, *mailFrom, *mailTo, subject, body)
						if err != nil {
							log.Printf("Failed to send mail: %v", err)
						} else {
							t.NotifiedAt = now
							log.Printf("Sent critical mail: %s", subject)
						}
					}
				} else {
					// Normal thresholds
					if t.NotifiedAt.IsZero() || now.Sub(t.NotifiedAt) >= *interval {
						subject := fmt.Sprintf("[WARNING] Maildir usage %0.2f%% >= %d%%", usedPercent, t.Percent)
						body := fmt.Sprintf("Maildir at %s is at %0.2f%% usage.", *maildir, usedPercent)
						err := sendMail(*smtpHost, *smtpPort, *smtpUser, *smtpPass, *mailFrom, *mailTo, subject, body)
						if err != nil {
							log.Printf("Failed to send mail: %v", err)
						} else {
							t.NotifiedAt = now
							log.Printf("Sent mail: %s", subject)
						}
					}
				}
			}
		}

		// Check growth in usage
		if now.Sub(lastUsageTime) >= *growthInterval {
			delta := usedPercent - lastUsage
			if delta >= *growthThreshold {
				subject := fmt.Sprintf("[CRITICAL] Maildir grew %0.2f%% in %v", delta, *growthInterval)
				body := fmt.Sprintf("Maildir at %s grew by %0.2f%% in %v. Current usage: %0.2f%%", *maildir, delta, *growthInterval, usedPercent)
				err := sendMail(*smtpHost, *smtpPort, *smtpUser, *smtpPass, *mailFrom, *mailTo, subject, body)
				if err != nil {
					log.Printf("Failed to send growth alert: %v", err)
				} else {
					log.Printf("Sent growth alert: %s", subject)
				}
			}
			lastUsage = usedPercent
			lastUsageTime = now
		}

		time.Sleep(*interval)
	}
}
