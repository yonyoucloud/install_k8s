package args

import (
	"net"
)

var Holder = &holder{}

type holder struct {
	port        int
	bindAddress net.IP
}

// GetPort 'port' argument of Dashboard binary.
func (self *holder) GetPort() int {
	return self.port
}

// GetBindAddress 'bind-address' argument of Dashboard binary.
func (self *holder) GetBindAddress() net.IP {
	return self.bindAddress
}
