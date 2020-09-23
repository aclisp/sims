Simple IM Server
---

* [x] Vendor go-micro and micro
* [x] Pure golang grpc client, in addition to pure HTTP client
  + Services are exposed via `api` using HTTP
  + Could add another API gateway as [grpc reverse proxy](pkg/grpcproxy/README.md) (TODO)
  + SIMS can be started with bind address
* [x] Simplified proto naming
* [x] Reliable publishing: error correct and tracable
* [x] add go-micro & grpc debugging and tracing
* [ ] Real world deployment
* [ ] Performance testing
  + Preliminary conclusion: memory bound: 1G mem ~ 10k user
* [ ] Authentication at API gateway with wechat
* [ ] Authorization on event publishing

Directory Layout
---

* `bin/` build output
* `client/`  client SDKs
  + `.` go-micro client for easy testing
  + `go` golang sdk
  + `js` typescript sdk
  + `java` java sdk
* `pkg/` reusable lib
  + `codec` ???
  + `grpcproxy` grpc transparent reverse proxy
  + `go-micro` modified go-micro base on v2.9.1
* `proto/` protobuf definitions
* `pub/` event publisher
* `server/` the sims server
* `micro/` the micro cli (modified)
* `tools/` dependent tools
  + `bin` protoc-gen-go protoc-gen-micro micro

Debugging
---

1. Start server with `--pprof_address 127.0.0.1:6060`
2. look for
   + http://127.0.0.1:6060/debug/pprof/
   + http://127.0.0.1:6060/debug/requests
   + http://127.0.0.1:6060/debug/events
3. micro debug interface
   + micro call go.micro.srv.sims Debug.Health
   + micro call go.micro.srv.sims Debug.Stats
   + micro call go.micro.srv.sims Debug.Cache
4. start with debug level logging
   + MICRO_LOG_LEVEL=debug bin/server --server_address :18080 --pprof_address :6060
   + MICRO_LOG_LEVEL=debug ./micro api --type srv
