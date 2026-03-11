# Failure Modes

This document lists known failure scenarios and the expected behavior of the relay.

The goal is to ensure predictable system behavior during faults.

---

## 1 Network Failure

Scenario:

The relay loses network connectivity to the Microsoft Graph API.

Expected behavior:

• delivery attempts fail  
• error is logged  
• message remains in queue  
• worker retries later  

Messages must not be deleted.

---

## 2 Graph API Rejection

Scenario:

The Graph API returns authentication or API errors.

Expected behavior:

• delivery fails  
• error is logged  
• message remains in queue  
• retry occurs later  

Messages must not be silently discarded.

---

## 3 Relay Crash

Scenario:

The relay process terminates unexpectedly.

Expected behavior:

• queued messages remain on disk  
• no message loss  
• after restart the worker resumes delivery  

Filesystem queue guarantees persistence.

---

## 4 Disk Full

Scenario:

The filesystem cannot accept new writes.

Expected behavior:

• SMTP ingest cannot persist message  
• SMTP server must reject the message  
• operator must resolve disk issue  

Messages must never be partially written.

---

## 5 Worker Failure

Scenario:

The delivery worker stops processing.

Expected behavior:

• messages accumulate in the queue  
• system remains safe  
• operator restart restores processing  

Queue ensures message durability.

---

## 6 Provider Rate Limiting

Scenario:

Microsoft Graph API returns rate limit responses.

Expected behavior:

• delivery temporarily fails  
• worker retries later  
• messages remain in queue  

Delivery must resume automatically.

---

# Failure Handling Principle

The relay follows the rule:

Never lose a message.

Failure must result in retry, not deletion.
