package concurrency

import (
	"context"
	"sync"
)

func ProccessConcurrently[I any, O any](c context.Context, input <-chan I,
	process func(context.Context, I) (O, error), workers int) <-chan O {
	output := make(chan O)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for data := range input {
				res, err := process(c, data)
				if err != nil {
					select {
					case <-c.Done():
						return
					default:
						continue
					}
				}
				select {
				case output <- res:
				case <-c.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
