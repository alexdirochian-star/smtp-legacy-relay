\# IEG Project



\## Project Purpose



The Infrastructure Email Gateway (IEG) is a minimal compatibility layer designed to preserve email delivery from legacy systems that cannot support modern authentication methods required by cloud email providers.



IEG operates as a local store-and-forward SMTP relay. It accepts messages from legacy systems and forwards them to modern email providers using a compatible outbound connection.



The purpose of the project is to prevent operational email outages caused by mandatory authentication changes in cloud email platforms.



\## Infrastructure Failure



The infrastructure failure addressed by this project is the SMTP authentication rejection:



550 5.7.30 Authentication unsuccessful



This error occurs when a legacy device or application attempts to send email using basic SMTP authentication with a username and password, but the upstream cloud email provider no longer accepts this authentication method.



The rejection happens during the SMTP authentication phase, and the message is rejected before entering the mail delivery pipeline.



\## Affected Systems



This issue affects systems that can send SMTP mail but cannot be upgraded to support modern authentication mechanisms such as OAuth2.



Typical affected systems include:



\- printers and scanners  

\- industrial controllers  

\- monitoring and alerting systems  

\- building management systems  

\- medical and laboratory equipment  

\- backup and notification software  

\- legacy enterprise applications  



These systems typically support only SMTP AUTH using static credentials and cannot be modified or upgraded.



\## Compatibility Adapter Concept (IEG)



IEG acts as a compatibility adapter placed between legacy SMTP clients and modern cloud email infrastructure.



The gateway performs the following functions:



1\. Accept SMTP messages from legacy systems  

2\. Store messages in a durable local queue  

3\. Forward messages to the upstream email provider through a modern-compatible connection  

4\. Retry delivery if the upstream provider is temporarily unavailable  



The architecture follows a store-and-forward model similar to traditional mail transfer agents.



This approach ensures message persistence and delivery reliability even when upstream services experience temporary failures.



\## Why the Solution Exists



Many legacy systems remain operational but cannot be upgraded to meet modern cloud email authentication requirements.



In many environments:



\- firmware updates are unavailable  

\- vendor support has ended  

\- replacing the system is expensive or operationally risky  



When cloud providers disable legacy authentication methods, these systems immediately lose the ability to send alerts, reports, and automated notifications.



IEG restores this capability by introducing a small compatibility node that bridges the protocol gap between legacy SMTP clients and modern authenticated email services.



This allows organizations to maintain operational continuity without replacing existing infrastructure.

