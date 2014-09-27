ttunnel
=======

TLS tunneling software written in Golang.

## Overview

The ttunnel bin directory contains three applications: 

* ttunnel-client.go
* ttunnel-config.go
* ttunnel-server.go

These should be fairly self-explanitory. The ttunnel-config program is
a small text-base interface for adminstering the server, and
generating client tokens.

## Configuration

### Server

#### Create the initial configuration

In order to create the initial configuration file for the server, run

```
ttunnel-config
``` 

Follow the prompts to create a new configuration file. Exit the
program when done.

The new configuration file will be in `~/.ttunnel/server.json`.

#### Install certificates

You'll need a key and signed certificate for the server. Instructions
for generating certificates can be easily found on the web. The key
and certificate file should be placed in the following files:

```
~/.ttunnel/server.key
~/.ttunnel/server.crt
```

### Client




