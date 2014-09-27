ttunnel
=======

TLS tunneling software written in Golang.

## Overview

The ttunnel software provides functionality that is similar to
`stunnel`. The important differences are

* Clients are authenticated using tokens that have an expiration and
  can be easily revoked. 
* It uses Go's TLS library, providing forward secrecy by default. 

The ttunnel bin directory contains three applications:

* ttunnel-client.go
* ttunnel-config.go
* ttunnel-server.go

These should be fairly self-explanitory. The `ttunnel-config` program is
a small text-base interface for adminstering the server and
generating client tokens.

## Configuration

### Server

#### Create the initial configuration

In order to create the initial configuration file for the server, run
`ttunnel-config `.

Follow the prompts to create a new configuration file. Exit the
program when done.

The new configuration file will be in `~/.ttunnel/server.json`.

#### Install certificates

You'll need a key and signed certificate for the server. Instructions
for generating certificates can be found on the web. The key and
certificate file should be placed in the following locations:

```
~/.ttunnel/server.key
~/.ttunnel/server.crt
```

##### Custom certificate authority

If you're using a custom CA, you'll need to place the root's
certificate in `~/.ttunnel/rootCA.crt` on the server and on all
clients.

#### Run the server

The server can be run by calling `ttunnel-server`. 

### Client

#### Generating the client configuration

The client's configuration will be generated on the *server*, and will
contain an token that allows the client to connect to a specified
address and port.

On the server run `ttunnel-config` and select `Add client`. Follow the
prompts to create a configuration file for the client. 

When you are done, a configuration file will be created in
`~/.ttunnel/tunnels/<name>.json`. Copy this file to the same location
on the client. 

Remember that you'll need the appropriate `~/.ttunnel/rootCA.crt` file
if you're using a custom CA.

#### Start the client

You can start the tunnel by calling `ttunnel-client <name>`, where
`<name>` is the name of the configuration file without the `.json`
extension.

## Notes

This is currently a work in progress. Comments and questions are
welcomed.
