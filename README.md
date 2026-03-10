# Fix for 550 5.7.30 SMTP Auth Shutdown

Tiny SMTP store-and-forward relay for legacy devices and applications that can no longer send email through modern providers.

Many printers, scanners, ERP systems and old internal apps still support only basic SMTP authentication.

Modern providers reject them with errors like:

```
550 5.7.30 Basic authentication is not supported for Client Submission.
```

This relay accepts legacy SMTP locally, stores messages on disk, and forwards them to an upstream SMTP provider.

## Architecture

legacy device -> local relay -> disk queue -> worker -> upstream SMTP

## Use cases

This relay is useful when legacy systems can still send SMTP mail but cannot support modern authentication.

Typical examples:

- Office printers (scan-to-email)
- Network scanners
- Legacy ERP systems
- Monitoring systems
- Old internal scripts
