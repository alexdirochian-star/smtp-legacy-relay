# Architecture

## System Flow

Device
?
SMTP Ingest
?
Filesystem Queue
?
Worker
?
Provider Adapter
?
Microsoft Graph API

## Architecture Decisions

ADR-0001  
Maildir queue pattern  
tmp ? fsync ? rename ? new  

The filesystem queue uses the Maildir write pattern to guarantee atomic message persistence.

ADR-0002  
Delivery model = at-least-once  

The system guarantees that messages are delivered at least once, ensuring reliability even if retries are required.

ADR-0003  
SQLite used only as message index  

SQLite stores message metadata only, while the filesystem queue remains the source of truth for message bodies.

## Component Responsibilities

SMTP Ingest
Handles SMTP protocol and writes incoming messages to the queue.

Filesystem Queue
Persistent storage for all messages.

Delivery Worker
Reads queued messages and performs delivery attempts.

Provider Adapter
Converts internal message format into Microsoft Graph API calls.
