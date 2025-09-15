package main

// testMail sends a test email to verify SMTP settings are working.
// It reuses sendMail() without modifying it.
func testMail(smtpHost string, smtpPort int, username, password, from, to string) error {
	subject := "[TEST] Maildir Watchdog SMTP Test"
	body := "This is a test email from the Maildir Watchdog program.\nIf you received this, SMTP settings are working correctly."

	return sendMail(smtpHost, smtpPort, username, password, from, to, subject, body)
}
