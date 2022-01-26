//go:build linux

package tap

import "github.com/songgao/water"

func configureForOS(config *water.Config) {
	config.Name = "maveth"
}
