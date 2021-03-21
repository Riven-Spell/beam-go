# This package is purely internal.

It specifies how we interact with the various Beam API executables (`wallet-api`, `explorer-node`, etc.) over HTTP and over TCP.

It enables the use of HTTP/S, and TCP over TLS. Both HTTPEndpoint and TCPEndpoint are valid to supply to the `wallet.Client` and `explorer.Client`.

