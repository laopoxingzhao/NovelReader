package test

import (
	"NovelReader/utils"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPublishAndSubscribe(t *testing.T) {

	var w sync.WaitGroup

	publisher := utils.GetPublisher()
	var a int32 = 1
	publisher.Subscribe("add", func(data any) {
		atomic.AddInt32(&a, 1)
	})
	publisher.Subscribe("reduce", func(data any) {
		atomic.AddInt32(&a, -1)
	})
	for i := 0; i < 1000; i++ {
		go publisher.Publish("add", func(data any) {
			w.Add(1)
			a++
			w.Done()

		})
		go publisher.Publish("reduce", func(data any) {
			w.Add(1)
			a--
			w.Done()
		})
	}

	w.Wait()
	fmt.Println(a)
	//canvas.NewImageFromURI()
}
