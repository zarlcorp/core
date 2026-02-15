package zapp

// WithName overrides the default application name.
func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}
