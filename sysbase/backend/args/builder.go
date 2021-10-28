package args

import "net"

var builder = &holderBuilder{holder: Holder}

// Used to build argument holder structure. It is private to make sure that only 1 instance can be created
// that modifies singleton instance of argument holder.
type holderBuilder struct {
	holder *holder
}

// SetPort 'port' argument of Dashboard binary.
func (self *holderBuilder) SetPort(port int) *holderBuilder {
	self.holder.port = port
	return self
}

// SetBindAddress 'bind-address' argument of Dashboard binary.
func (self *holderBuilder) SetBindAddress(ip net.IP) *holderBuilder {
	self.holder.bindAddress = ip
	return self
}

// GetHolderBuilder returns singleton instance of argument holder builder.
func GetHolderBuilder() *holderBuilder {
	return builder
}
