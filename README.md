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

To run the application you will first need to set the following environment variables:

- `MAILINABOX_USER` - A user with access to the Mail-In-A-Box API (admin panel)
- `MAILINABOX_PASSWORD` - The password of the user
- `MAILINABOX_HOSTNAME` - The hostname of the Mail-In-A-Box server (i.e. `box.example.com`)

You can provide these via a `.env` file in the same directory as the application, or by setting them in the environment.

Then, simply run the application using the following command (it will run once and exit):

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


### Running with PM2

You can also run the application with PM2 by using the following command:

```bash
pm2 start --name "CNameFlatten" ./CNameFlatten -- "*/5 * * * *"
```

This will run the process every 5 minutes.

## Example Setup

![Demo Image](https://github.com/CADawg/CNameFlatten/blob/main/demo/Screenshot_20240307_203950.png?raw=true)

Set up a TXT record as shown above, and the A and AAAA records will be created from the domain specified in the text record, in this case `domain-proxy.bearblog.dev`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the [GPL v3 Licence](LICENSE).
