# Chaos Testing

Failure simulation procedures for the relay system.

## Network outage

Simulate loss of network connectivity.

Expected behaviour:

Messages remain in the filesystem queue.
Delivery worker retries when connectivity returns.

## Graph API failure

Simulate Microsoft Graph API returning errors.

Expected behaviour:

Delivery attempts fail.
Messages remain in the queue.
Worker retries delivery later.

## Disk full

Simulate disk capacity exhaustion.

Expected behaviour:

SMTP ingest fails to write messages.
SMTP server must reject new messages.

## Worker crash

Simulate unexpected worker termination.

Expected behaviour:

Messages remain safely stored in the queue.
Worker resumes processing after restart.
