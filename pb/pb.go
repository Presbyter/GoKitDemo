//go:generate protoc --go_out=plugins=grpc:. usersvc.proto
//go:generate protoc --go_out=plugins=grpc:. authsvc.proto
package pb
