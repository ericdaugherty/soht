SOHT - Socket over HTTP Application in GoLang
----

Pseudo-port of original Java/.Net SOHT application: http://www.ericdaugherty.com/dev/soht

This is an application written in Go that proxies TCP Socket connections over an HTTP connection.  The goal of this tool is to provide Socket connectivity from a location that only provides outbound HTTP connections.  It uses a client that runs locally, and a remote server, to proxy the socket connection from the local computer over HTTP to the server, and then via a TCP Socket to the target computer.

This version is not compatible with the original implementation, but borrows heavily from its protocol specification and general approach.

