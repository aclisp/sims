PROTO=./proto

.PHONY: all
all:
	cd client && go build -o ../bin; cd ..
	cd server && go build -o ../bin; cd ..
	cd pub && go build -o ../bin; cd ..

.PHONY: proto
proto:
	protoc --proto_path=${PROTO} \
	--micro_out=${PROTO} --micro_opt=paths=source_relative \
	--go_out=${PROTO} --go_opt=paths=source_relative \
	${PROTO}/*.proto
