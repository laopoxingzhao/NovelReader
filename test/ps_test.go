package main

import (
	"awesomeProject/utils"
	"fmt"
	"testing"
)

func TestPs(t *testing.T) {
	publisher := utils.NewPublisher()
	publisher.Subscribe("test", func(data interface{}) {
		fmt.Println(data.(string))
	})
	publisher.Publish("test", "hello world")
	//canvas.NewImageFromURI()
}
