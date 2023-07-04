package keys

var namespace = ""
var namespaceKeyPrefix = ""

// Namespace returns the current global application namespace.
func Namespace() string {
	return namespace
}

func NamespaceKeyPrefix() string {
	return namespaceKeyPrefix
}

// SetNamespace sets the global application namespace.
func SetNamespace(ns string) {
	namespace = ns
	namespaceKeyPrefix = ns + KeySeparator
}
