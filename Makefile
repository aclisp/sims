PROTO=./proto

.PHONY: server
server:
	cd server && go build -o ../bin; cd ..

.PHONY: client
client:
	cd client && go build -o ../bin; cd ..

.PHONY: pub
pub:
	cd pub && go build -o ../bin; cd ..

.PHONY: all
all:
	cd client && go build -o ../bin; cd ..
	cd server && go build -o ../bin; cd ..
	cd pub && go build -o ../bin; cd ..

.PHONY: linux
linux:
	cd server && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/sims; \
	cd ..

.PHONY: proto
proto:
	protoc --proto_path=${PROTO} \
	--micro_out=${PROTO} --micro_opt=paths=source_relative \
	--go_out=${PROTO} --go_opt=paths=source_relative \
	${PROTO}/*.proto

.PHONY: protoc
protoc:
	protoc --proto_path=${PROTO} \
	--go_out=plugins=grpc:${PROTO}/go --go_opt=paths=source_relative \
	${PROTO}/*.proto

.PHONY: lint
lint:
	~/go/bin/golint server/... pub/... client/...
	gofmt -l -w -s server pub client
