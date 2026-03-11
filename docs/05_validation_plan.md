# Validation Plan

This document defines the engineering validation procedure for the relay.

The goal is to verify that the system behaves according to the documented architecture and reliability model.

All tests must be executed before production deployment.

---

# 1. Environment Setup

Requirements:

• Go installed
• relay binary compiled
• Microsoft 365 mailbox configured
• network access available

Build the relay:

go build .

Start the relay:

.\relay.exe

Expected output:

SMTP relay listening on 127.0.0.1:2525

---

# 2. SMTP Ingest Test

Verify that the relay accepts SMTP connections.

Test using telnet:

telnet localhost 2525

Execute SMTP session:

HELO localhost
MAIL FROM:<test@example.com>
RCPT TO:<recipient@example.com>
DATA
test message
.
QUIT

Expected behavior:

• message accepted
• message written to queue directory

Log output example:

queued message: queue\<timestamp>.eml

---

# 3. Queue Persistence Test

Verify that messages remain on disk until delivered.

Procedure:

1. stop the relay
2. send a message to the queue
3. confirm message file exists

Expected result:

queue directory contains:

*.eml

Restart relay.

Expected behavior:

worker resumes processing queued message.

---

# 4. Crash Recovery Test

Verify crash safety of the filesystem queue.

Procedure:

1. send test message
2. terminate relay during delivery
3. restart relay

Expected behavior:

• message is not lost
• message delivery resumes

---

# 5. Worker Retry Behavior

Simulate delivery failure.

Procedure:

1. disable network access
2. send test message
3. observe delivery failure

Expected log output:

send error: <provider response>

Expected behavior:

message remains in queue for retry.

---

# 6. Provider Adapter Test

Verify Graph API delivery.

Procedure:

1. configure Graph API credentials
2. send test message
3. confirm delivery to mailbox

Expected result:

recipient mailbox receives message.

---

# 7. Delivery Confirmation

Verify queue cleanup.

Expected behavior:

successful delivery removes message file.

Queue directory should not accumulate messages.

---

# 8. Validation Success Criteria

The relay passes validation when:

• SMTP ingestion works
• messages persist on disk
• crash recovery succeeds
• delivery retries function
• Graph API delivery succeeds
• queue is cleaned after delivery

---

# Validation Result

PARTIAL PASS
