// Package zoptions provides a generic option pattern for configuring structs.
//
// This package offers a standardized way to implement the functional options
// pattern across all projects, ensuring consistency and reusability.
//
// # Usage Example
//
//	type MyConfig struct {
//		Name    string
//		Timeout time.Duration
//	}
//
//	func NewMyConfig(opts ...zoptions.Option[MyConfig]) *MyConfig {
//		cfg := &MyConfig{
//			Name:    "default",
//			Timeout: 30 * time.Second,
//		}
//
//		for _, opt := range opts {
//			opt(cfg)
//		}
//
//		return cfg
//	}
//
//	func WithName(name string) zoptions.Option[MyConfig] {
//		return func(cfg *MyConfig) {
//			cfg.Name = name
//		}
//	}
//
//	func WithTimeout(timeout time.Duration) zoptions.Option[MyConfig] {
//		return func(cfg *MyConfig) {
//			cfg.Timeout = timeout
//		}
//	}
//
//	// Usage
//	cfg := NewMyConfig(
//		WithName("production"),
//		WithTimeout(60 * time.Second),
//	)
package zoptions

// Option is a generic functional option that modifies a struct of type T.
// It follows the standard Go functional options pattern where options
// are functions that modify the target struct in place.
type Option[T any] func(*T)
