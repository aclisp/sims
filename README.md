Simple IM Server
---

* [ ] Vendor go-micro and micro
* [ ] Pure golang grpc client, is it necessary?
  + No, you don't
  + Services are exposed via `api` using HTTP
  + Add another API gateway as grpc reverse proxy
* [ ] Simplified proto naming

Directory Layout
---

* `bin/` build output
* `client/`  client SDKs
  + `.` go-micro client for easy testing
  + `go` golang sdk
  + `js` typescript sdk
  + `java` java sdk
* `gate/` API gateway
* `pkg/` reusable lib
  + `codec` ???
  + `grpcproxy` grpc transparent reverse proxy
  + `go-micro` modified go-micro base on v2.9.1
* `proto/` protobuf definitions
* `pub/` event publisher
* `server/` the sims server
* `tools/` dependent tools
  + `bin` protoc-gen-go protoc-gen-micro micro
