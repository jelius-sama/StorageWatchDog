package main

import (
	"fmt"
	"syscall"
)

// diskSizeStats returns total GiB and used GiB (without recalculating percent).
func diskSizeStats(path string) (totalGiB float64, usedGiB float64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}
	total := float64(stat.Blocks) * float64(stat.Bsize)
	free := float64(stat.Bavail) * float64(stat.Bsize)
	used := total - free

	if total == 0 {
		return 0, 0, fmt.Errorf("total blocks is zero")
	}

	totalGiB = total / (1024 * 1024 * 1024)
	usedGiB = used / (1024 * 1024 * 1024)

	return totalGiB, usedGiB, nil
}

// testMail sends a test email to verify SMTP settings and report disk usage.
func testMail(smtpHost string, smtpPort int, username, password, from, to string) error {
	// reuse getDiskUsage for percentage
	usedPercent, err := getDiskUsage("/var/vmail")
	if err != nil {
		return err
	}

	// get GiB stats separately
	totalGiB, usedGiB, err := diskSizeStats("/var/vmail")
	if err != nil {
		return err
	}

	// sanity check: recompute percent from GiB to see if it matches getDiskUsage
	computedPercent := (usedGiB / totalGiB) * 100

	subject := "[TEST] Maildir Watchdog SMTP Test"
	body := fmt.Sprintf(
		"This is a test email from the Maildir Watchdog program.\n"+
			"If you received this, SMTP settings are working correctly.\n\n"+
			"Disk Usage Report:\n"+
			" - Used: %.2f GiB\n"+
			" - Total: %.2f GiB\n"+
			" - Usage (getDiskUsage): %.2f%%\n"+
			" - Usage (GiB recomputed): %.2f%%",
		usedGiB, totalGiB, usedPercent, computedPercent,
	)

	return sendMail(smtpHost, smtpPort, username, password, from, to, subject, body)
}
