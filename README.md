# CNAME Flattener for Mail-In-A-Box

This application flattens CNAME records to A and AAAA records for Mail-In-A-Box by using the API and checking for TXT records named `_cname_flatten.yourdomain.com`. It then looks up the value of the TXT record on a DNS server and sets the A and AAAA records of the parent domain (e.g., `yourdomain.com`) to those of the target in the TXT record.

## Prerequisites

- Go 1.22 or later

## Installation

You can install this application using the Go Install command:

```bash
go install github.com/CADawg/CNameFlatten@latest
```

## Usage

You can run the application in two ways:

### One-time Execution

To run the application once, use the following command:

```bash
CNameFlatten
```

### Scheduled Execution

To run the application repeatedly on a schedule, provide a cron expression as an argument. For example, to run the application every minute, use the following command:

```bash
CNameFlatten "*/1 * * * *"
```

Or, to run the application every hour, use:

```bash
CNameFlatten "0 * * * *"
```

The cron expression follows the standard format:

```
*    *    *    *    *
┬    ┬    ┬    ┬    ┬
│    │    │    │    │
│    │    │    │    └ day of week (0 - 7) (0 or 7 is Sun)
│    │    │    └───── month (1 - 12)
│    │    └────────── day of month (1 - 31)
│    └─────────────── hour (0 - 23)
└──────────────────── minute (0 - 59)
```

## Example Setup

![Demo Image](./Screenshot_20240307_203950.png)

Setup a TXT record as shown above, and the A and AAAA records will be created from the domain specified in the text record, in this case `domain-proxy.bearblog.dev`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the [GPL v3 Licence](LICENSE).
