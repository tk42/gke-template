package threadsafe

import "context"

type RingBuffer[T any] struct {
	inputChannel  <-chan T
	outputChannel chan T
}

func NewRingBuffer[T any](inputChannel <-chan T, outputChannel chan T) *RingBuffer[T] {
	return &RingBuffer[T]{inputChannel, outputChannel}
}

func (r *RingBuffer[T]) Run(ctx context.Context) {
	defer close(r.outputChannel)
	for v := range r.inputChannel {
		select {
		case <-ctx.Done():
			return
		case r.outputChannel <- v:
		default:
			<-r.outputChannel
			r.outputChannel <- v
		}
	}
}
