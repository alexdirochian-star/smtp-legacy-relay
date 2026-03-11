# Chaos Test Plan

## Purpose

This document defines reliability validation tests for the relay.

The relay is designed as a store-and-forward message relay using a filesystem queue as the source of truth. Chaos tests validate that the system behaves correctly under common infrastructure failures.

Reliability invariant:

No accepted message may be lost under crashes, network outages, or provider errors.

A message is considered accepted only after durable queue persistence using the atomic pattern:

write -> fsync -> rename


# Failure Scenarios

## Worker Crash

Scenario:
The relay process is terminated while messages are queued or currently being processed.

Test Procedure:
1. Send several messages to the relay.
2. Verify messages appear in the queue directory.
3. Forcefully terminate the relay process.
4. Restart the relay.

Expected Result:
Messages remain in the queue and are delivered after restart.


## Network Failure

Scenario:
Connectivity to Microsoft Graph API is temporarily unavailable.

Test Procedure:
1. Send several messages to the relay.
2. Disable outbound connectivity or block the Graph endpoint.
3. Observe worker behavior.

Expected Result:
Messages remain queued and retry later when connectivity is restored.


## API Temporary Error

Scenario:
The provider returns temporary failures (5xx or retryable responses).

Test Procedure:
1. Simulate temporary provider errors.
2. Observe retry behavior.

Expected Result:
Worker retries delivery without message loss.


## Disk Full

Scenario:
Queue disk becomes full.

Test Procedure:
1. Fill the disk where the queue directory resides.
2. Attempt to submit new messages.

Expected Result:
SMTP ingest refuses new messages while preserving the existing queue.


## Queue Recovery

Scenario:
Relay restarts after crash while queue contains messages.

Test Procedure:
1. Send messages.
2. Crash the relay.
3. Restart the relay.

Expected Result:
Queued messages resume delivery automatically.


## Load Test

Scenario:
High message throughput.

Test Procedure:
Send 1000 messages through the relay.

Expected Result:
All messages are queued and processed successfully.


# Success Criteria

Lost messages = 0.


# Failure Mode Matrix

| Failure Mode | Trigger | Expected Behavior |
|--------------|--------|------------------|
| Worker Crash | Relay process killed | Messages remain in queue and deliver after restart |
| Network Failure | Graph API unreachable | Messages remain queued and retry later |
| API Temporary Error | Provider returns temporary failure | Worker retries delivery without message loss |
| Disk Full | Queue disk exhausted | SMTP ingest rejects new messages while preserving queue |
| Queue Recovery | Relay restart | Existing queue resumes processing |
| Load Test | 1000 messages sent | All messages processed with zero loss |
