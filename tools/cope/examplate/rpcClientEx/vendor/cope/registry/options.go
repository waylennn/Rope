package registry

import "time"

// Options for service registration
type Options struct {
	Addrs     []string
	HeartBeat int64
	TimeOut   time.Duration
	// example:  /xxx_company/app/kuaishou/service_A/10.192.1.1:8801
	// example:  /xxx_company/app/kuaishou/service_A/10.192.1.2:8801
	RegistryPath string
}

// Option ...
type Option func(opts *Options)

// WithTimeout add TimeOut
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.TimeOut = timeout
	}
}

// WithAddrs add addrs
func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

// WithHeartBeat add HeartBeat
func WithHeartBeat(heartHeat int64) Option {
	return func(opts *Options) {
		opts.HeartBeat = heartHeat
	}
}

// WithRegistryPath add RegistryPath
func WithRegistryPath(path string) Option {
	return func(opts *Options) {
		opts.RegistryPath = path
	}
}
