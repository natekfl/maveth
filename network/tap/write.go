package tap

import (
	"errors"
)

func SendPacket(payload []byte) error {
	if ifce == nil {
		return errors.New("no interface")
	}
	_, err := ifce.Write(payload)
	return err
}
