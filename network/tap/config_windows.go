//go:build windows

package tap

import "github.com/songgao/water"

func configureForOS(config *water.Config) {
	config.ComponentID = "root\\tap0901"
}
