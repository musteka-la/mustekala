## TODO

### server.go

* `getPrivateKey()`: The private key should not be hardcoded.
* `getClientName()`: Generate the client name using the git revision and version.

### handle-blockheader-msg.go

* `handleBlockHeaderMsg()`
  * Validate the sealing (pow) of each received header.
  * Ship received headers to our DB.