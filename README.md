# goqog
goqog (pronounced 'go cog') is a process launcher designed for embedded devices that need configuration. It provides an HTTP server for control and configuration, and a stdio interface for launched processes to indicate status, receive configuration, and register for behavior triggers.

## Feature Goals
* Launch applets/commands on boot (and re-launch on crash)
* Store configuration
* Permit editing configuration via HTTP
* Permit control and feedback of applets via signals exchanged on stdio

