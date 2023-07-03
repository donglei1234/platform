package query

import (
	"reflect"

	"github.com/donglei1234/platform/services/common/jsonx"
)

type action func(obj interface{}) error

type options struct {
	input   []byte
	actions []action
}

type Option func(o *options)

func newOptions(opts ...Option) (o options) {
	for _, opt := range opts {
		opt(&o)
	}
	return
}

func WithInput(data interface{}) Option {
	return func(o *options) {
		switch dt := data.(type) {
		case string:
			o.input = []byte(dt)
		case []byte:
			o.input = dt
		case map[string]interface{}:
			o.input, _ = jsonx.Marshal(data)
		}
	}
}

func WithSetField(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					p[nodes[0].Key] = v
				case *interface{}:
					*p = v
				case []interface{}:
					if nodes[0].Index < len(p) {
						p[nodes[0].Index] = v
					} else {
						for loop := len(p); loop < nodes[0].Index; loop++ {
							p = append(p, nil)
						}
						p = append(p, v)
						n1 := nodes[1]
						if n.Index == -1 || n1.Key == "" {
							return ErrInvalidDestination
						} else {
							switch pp := n1.Parent.(type) {
							case map[string]interface{}:
								pp[n1.Key] = p
							default:
								return ErrInvalidDestination
							}
						}
					}
				default:
					return ErrInvalidDestination
				}
				return nil
			}
		})
	}
}

func WithIncrement(path []byte, increment int) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						value := int(p[n.Key].(uint64))
						p[n.Key] = value + increment
					}
				case []interface{}:
					if n.Index == -1 {
						return ErrInvalidDestination
					} else {
						value := int(p[n.Index].(uint64))
						p[n.Index] = value + increment
					}
				default:
					return ErrInvalidDestination
				}
				return nil
			}
		})
	}
}

func WithArrayInsert(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if len(nodes) < 2 {
				return ErrInvalidDestination
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				n1 := nodes[1]
				switch p := n1.Parent.(type) {
				case map[string]interface{}:
					if n1.Key == "" || n.Index == -1 {
						return ErrInvalidDestination
					} else {
						p[n1.Key] = append(p[n1.Key].([]interface{}), v)
						copy(p[n1.Key].([]interface{})[n.Index+1:], p[n1.Key].([]interface{})[n.Index:])
						p[n1.Key].([]interface{})[n.Index] = v
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func WithArrayPushFront(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						p[n.Key] = append([]interface{}{v}, p[n.Key].([]interface{})...)
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func WithArrayPushBack(path []byte, value []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						p[n.Key] = append(p[n.Key].([]interface{}), v)
					}
				case *interface{}:
					switch (*p).(type) {
					case []interface{}:
						*p = append((*p).([]interface{}), v)
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func WithArrayUnique(path []byte, value []byte) Option {

	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else if v, err := jsonx.ParseAndReturn(value); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						for _, a := range p[n.Key].([]interface{}) {
							if reflect.DeepEqual(a, v) {
								return nil
							}
						}
						p[n.Key] = append(p[n.Key].([]interface{}), v)
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func WithDelete(path []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodes, err := pathToProperty(root, path); err != nil {
				return err
			} else {
				n := nodes[0]
				switch p := n.Parent.(type) {
				case map[string]interface{}:
					if n.Key == "" {
						return ErrInvalidDestination
					} else {
						delete(p, n.Key)
					}
				case []interface{}:
					n1 := nodes[1]
					if n.Index == -1 || n1.Key == "" {
						return ErrInvalidDestination
					} else {
						switch p := n1.Parent.(type) {
						case map[string]interface{}:
							p[n1.Key] = append(p[n1.Key].([]interface{})[:n.Index], p[n1.Key].([]interface{})[n.Index+1:]...)
						case []interface{}:
							p[n1.Index] = append(p[n1.Index].([]interface{})[:n.Index], p[n1.Index].([]interface{})[n.Index+1:]...)
						default:
							return ErrInvalidDestination
						}
					}

				default:
					return ErrInvalidDestination
				}
			}

			return nil
		})
	}
}

func WithCopy(pathFrom []byte, pathTo []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodesFrom, err := pathToProperty(root, pathFrom); err != nil {
				return err
			} else if nodesTo, err := pathToProperty(root, pathTo); err != nil {
				return err
			} else {
				to := nodesTo[0]
				from := nodesFrom[0]
				switch p := to.Parent.(type) {
				case map[string]interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Key == "" || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Key]
						}
					case []interface{}:
						if to.Key == "" || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Index]
						}
					default:
						return ErrInvalidDestination
					}
				case []interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Index == -1 || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Key]
						}
					case []interface{}:
						if to.Index == -1 || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Index]
						}
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func WithMove(pathFrom []byte, pathTo []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodesFrom, err := pathToProperty(root, pathFrom); err != nil {
				return err
			} else if nodesTo, err := pathToProperty(root, pathTo); err != nil {
				return err
			} else {
				to := nodesTo[0]
				from := nodesFrom[0]
				switch p := to.Parent.(type) {
				case map[string]interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Key == "" || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Key]
							delete(f, from.Key)
						}
					case []interface{}:
						if to.Key == "" || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Key] = f[from.Index]
							f = append(f[:from.Index], f[from.Index+1:]...)
						}
					default:
						return ErrInvalidDestination
					}
				case []interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Index == -1 || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Key]
							delete(f, from.Key)
						}
					case []interface{}:
						if to.Index == -1 || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Index] = f[from.Index]
							f = append(f[:from.Index], f[from.Index+1:]...)
						}
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}

func WithSwap(pathFrom []byte, pathTo []byte) Option {
	return func(o *options) {
		o.actions = append(o.actions, func(root interface{}) error {
			if nodesFrom, err := pathToProperty(root, pathFrom); err != nil {
				return err
			} else if nodesTo, err := pathToProperty(root, pathTo); err != nil {
				return err
			} else {
				to := nodesTo[0]
				from := nodesFrom[0]
				switch p := to.Parent.(type) {
				case map[string]interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Key == "" || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Key], f[from.Key] = f[from.Key], p[to.Key]
						}
					case []interface{}:
						if to.Key == "" || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Key], f[from.Index] = f[from.Index], p[to.Key]
						}
					default:
						return ErrInvalidDestination
					}
				case []interface{}:
					switch f := from.Parent.(type) {
					case map[string]interface{}:
						if to.Index == -1 || from.Key == "" {
							return ErrInvalidDestination
						} else {
							p[to.Index], f[from.Key] = f[from.Key], p[to.Index]
						}
					case []interface{}:
						if to.Index == -1 || from.Index == -1 {
							return ErrInvalidDestination
						} else {
							p[to.Index], f[from.Index] = f[from.Index], p[to.Index]
						}
					default:
						return ErrInvalidDestination
					}
				default:
					return ErrInvalidDestination
				}
			}
			return nil
		})
	}
}
