package tap

import (
	"log"

	"github.com/songgao/water"
)

var ifce *water.Interface
var ifceWait = make(chan struct{})

//Blocks (always returns a non nil error)
func StartInterfaceRead(capture func([]byte)) (err error) {

	config := water.Config{
		DeviceType: water.TAP,
	}
	configureForOS(&config)
	ifce, err = water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	close(ifceWait)

	for {
		var frame = make([]byte, 1500)
		r, err := ifce.Read([]byte(frame))
		if err != nil {
			return err
		}
		frame = frame[:r]
		capture(frame)
	}
}

func WaitForConnection() {
	<-ifceWait
}
