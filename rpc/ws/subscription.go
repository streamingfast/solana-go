package ws

import (
	"context"
	"reflect"
)

type Subscription struct {
	req               *request
	subID             uint64
	stream            chan result
	err               chan error
	reflectType       reflect.Type
	closeFunc         func(err error)
	unsubscribeMethod string
}

func newSubscription(req *request, reflectType reflect.Type, closeFunc func(err error), unsubscribeMethod string) *Subscription {
	return &Subscription{
		req:               req,
		reflectType:       reflectType,
		stream:            make(chan result, 200_000),
		err:               make(chan error, 100_000),
		closeFunc:         closeFunc,
		unsubscribeMethod: unsubscribeMethod,
	}
}

// Recv waits for a message to arrive through the WebSocket connection for this exact subscription.
//
// It will either returns:
//  - `<value>, nil` If the stream has a message available
//  - `nil, err` If the subscription encounteted an error
//  - `nil, context.Canceled` If the context received was canceled
//  - `nil, context.DeadlineExceed` If the context timeout was reached
//
// Upon receiving a `context.Canceled` or `context.DeadlineExceed`, the subscription is
// automatically unsubscribed.
//
// *Future* It's not clear if `ctx context.Context` is appropriate here. Indeed, the gRPC
// way of doing things it to accept a `ctx context.Context` for the whole subscription
// lifecycle. The subscription is then tied to the context and closes automatically. The
// `Recv` method in this case does not accept a `ctx` object and uses the subscription
// one instead.
func (s *Subscription) Recv(ctx context.Context) (interface{}, error) {
	select {
	case d := <-s.stream:
		return d, nil
	case err := <-s.err:
		return nil, err
	case <-ctx.Done():
		s.unsubscribe(ctx.Err())
		return nil, ctx.Err()
	}
}

func (s *Subscription) Unsubscribe() {
	s.unsubscribe(nil)
}

func (s *Subscription) unsubscribe(err error) {
	s.closeFunc(err)

}
