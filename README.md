# Storage Watchdog

A lightweight Go service that monitors disk usage of a mail storage directory (or any directory) and sends alert emails when usage thresholds are reached or when rapid growth is detected.  

## Features

- Monitors disk usage of a specified path (default: `/var/vmail`).
- Sends email alerts via SMTP when usage exceeds thresholds (80%, 85%, 90%, 95%, 98% by default).
- Critical alerts (≥ 98%) are repeated at a configurable interval.
- Detects abnormal disk growth (default: 5% increase over 5 hours).
- Sends a test email on startup to verify SMTP configuration.
- Configurable via command-line flags.

## Installation

Clone the repository and build the binary or download any of the prebuilt binaries in the releases (Make sure it matches your OS and CPU architecture):

**Note: Of the prebuilt binaries only Linux x86_64 is tested on actual hardware (AWS EC2 t2.micro instance)**

```bash
git clone https://github.com/jelius-sama/StorageWatchDog.git
cd StorageWatchDog
CGO_ENABLED=0 go build -trimpath -buildvcs=false -o ./WatchDog ./
```

## Usage

Run the watchdog with appropriate flags:

```bash
./WatchDog \
  -path /var/vmail \
  -interval 10m \
  -warn80 80 \
  -warn85 85 \
  -warn90 90 \
  -warn95 95 \
  -warn98 98 \
  -criticalRepeat 30m \
  -growth 5 \
  -growthInterval 5h \
  -smtpHost smtp.example.com \
  -smtpPort 587 \
  -smtpUser user@example.com \
  -smtpPass password \
  -from alert@example.com \
  -to admin@example.com
```

### Key Flags

| Flag                  | Default             | Description                                    |
| --------------------- | ------------------- | ---------------------------------------------- |
| `-path`               | `/var/vmail`        | Directory to monitor                           |
| `-interval`           | `10m`               | Disk usage check interval                      |
| `-warn80` … `-warn98` | `80 … 98`           | Usage thresholds for warnings/critical alerts  |
| `-criticalRepeat`     | `30m`               | Repeat interval for critical alerts            |
| `-growth`             | `5`                 | Growth percentage threshold for critical alert |
| `-growthInterval`     | `5h`                | Interval to check growth                       |
| `-smtpHost`           | `smtp.example.com`  | SMTP server host                               |
| `-smtpPort`           | `587`               | SMTP server port                               |
| `-smtpUser`           | `user@example.com`  | SMTP username                                  |
| `-smtpPass`           | `password`          | SMTP password                                  |
| `-from`               | `alert@example.com` | Sender email address                           |
| `-to`                 | `admin@example.com` | Recipient email address                        |

## Example

To monitor `/data/mail` with custom thresholds and Gmail SMTP:

```bash
./WatchDog \
  -path /data/mail \
  -warn90 90 \
  -warn95 95 \
  -warn98 98 \
  -smtpHost smtp.gmail.com \
  -smtpPort 587 \
  -smtpUser myuser@gmail.com \
  -smtpPass myapppassword \
  -from myuser@gmail.com \
  -to sysadmin@company.com
```

## How It Works

1. On startup, a test email is sent to verify SMTP configuration and report initial disk usage.
2. The service checks disk usage at the configured interval.
3. If usage exceeds a threshold, an email alert is sent.
4. Critical alerts (≥ 98%) are repeated at the `-criticalRepeat` interval until resolved.
5. Periodically checks for disk usage growth. If usage grows more than the configured percentage within the interval, a critical alert is sent.

## Requirements

* Go 1.18+
* SMTP credentials with permission to send email

## License

MIT License

