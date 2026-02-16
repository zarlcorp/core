package zstore

// options holds configuration for the store.
type options struct{}

// Option configures a Store.
type Option func(*options)

func applyOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
