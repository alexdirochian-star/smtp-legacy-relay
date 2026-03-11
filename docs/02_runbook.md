# Runbook

Operational guide for running the relay.

This document is intended for system operators.

## 1 Starting the relay

Open a terminal and navigate to the relay directory.

Example:

cd C:\Users\diroc\Desktop\relay\cmd\relay

Start the relay:

relay.exe

If the relay starts successfully, the terminal will display:

SMTP relay listening on 127.0.0.1:2525

The relay is now ready to accept SMTP connections.

## 2 Queue directory location

All messages are stored in the queue directory.

Default location:

queue\

Example full path:

C:\Users\diroc\Desktop\relay\cmd\relay\queue\

Files in the queue represent messages waiting for delivery.

## 3 Log output

Logs are printed directly to the terminal.

Typical messages include:

SMTP relay listening
queued message
processing message
message delivered successfully
send error

Operators should monitor the terminal for repeated delivery errors.

## 4 Stopping the relay

To stop the relay press:

CTRL + C

The process will terminate immediately.

Queued messages remain on disk.

## 5 Restart procedure

If the relay must be restarted:

1 Navigate to the relay directory

cd C:\Users\diroc\Desktop\relay\cmd\relay

2 Start the relay again

relay.exe

The worker will automatically resume processing queued messages.

## 6 Checking message delivery

Delivery status can be verified in two ways.

Terminal logs

Successful delivery will appear as:

message delivered successfully

Failures will appear as:

send error

Queue directory

Check the queue directory.

Normal state:

queue gradually becomes empty.

Problem state:

many files remain in the queue for a long time.

This indicates an upstream SMTP or API delivery problem.
