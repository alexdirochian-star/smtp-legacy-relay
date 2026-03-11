# Incident Postmortem

Template for analyzing relay incidents.

## Incident Summary

Short description of the incident.

## Timeline

Chronological sequence of events during the incident.

## Root Cause

Technical explanation of the failure.

## Impact

Description of affected systems or message delivery.

## Corrective Actions

Steps required to prevent recurrence.

## Chaos Drill: Relay Crash Test

Test scenario:
Start relay ? send 1000 messages ? kill process ? restart relay.

Result:
Relay restarted successfully.
Queue persisted across crash.
Worker resumed delivery.
Lost messages = 0.

