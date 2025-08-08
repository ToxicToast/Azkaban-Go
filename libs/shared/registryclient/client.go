package registryclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type RequestFactory func() proto.Message
type GrpcInvoker func(ctx context.Context, cc *grpc.ClientConn, req proto.Message) (proto.Message, error)

type Entry struct {
	NewReq RequestFactory
	Invoke GrpcInvoker
}

type Registry map[string]Entry

func NewRegistry() Registry {
	return make(Registry)
}

func (r Registry) Register(key string, e Entry) {
	r[key] = e
}

func (r Registry) Get(key string) (Entry, bool) {
	e, ok := r[key]
	return e, ok
}
