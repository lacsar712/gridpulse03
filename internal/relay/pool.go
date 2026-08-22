package relay

import "sync"

type TaskPool struct{}

func (TaskPool) Run(jobs []func()) {
	var wg sync.WaitGroup
	for _, job := range jobs {
		if job == nil {
			continue
		}
		wg.Add(1)
		go func(j func()) {
			defer wg.Done()
			j()
		}(job)
	}
	wg.Wait()
}
