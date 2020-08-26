module github.com/aclisp/sims

go 1.14

require (
	go.uber.org/atomic v1.5.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/grpc v1.26.0
	google.golang.org/protobuf v1.23.0
)

replace github.com/micro/go-micro/v2 => ./pkg/go-micro
