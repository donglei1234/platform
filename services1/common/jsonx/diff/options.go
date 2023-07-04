package diff

// Options contains all the various options that the provided WithXyz functions construct.
type Options struct {
	IncludePaths []string
}

func (o *Options) IsIncludePath(path string) bool {
	if len(o.IncludePaths) == 0 {
		return true
	}
	for _, p := range o.IncludePaths {
		if len(p) > len(path) {
			continue
		} else if p == path {
			return true
		} else if len(path) > len(p) &&
			p == path[0:len(p)] &&
			(path[len(p):len(p)+1] == "." || path[len(p):len(p)+1] == "[") {
			return true
		}
	}
	return false
}

// Option is a closure that updates Options.
type Option func(o *Options)

// NewOptions constructs an Options struct from the provided Option closures and returns it.
func NewOptions(opts ...Option) (options *Options) {
	options = new(Options)
	for _, opt := range opts {
		opt(options)
	}
	return
}

// WithIncludePaths provides a []string value option. This is optionally used when diff documents with included paths only.
// the correct version before updating.
func WithIncludePaths(paths []string) Option {
	return func(o *Options) {
		o.IncludePaths = append(o.IncludePaths, paths...)
		return
	}
}
