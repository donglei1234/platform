package memory

import (
	"reflect"

	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"github.com/donglei1234/platform/services/common/nosql/memory/options"
)

//const MaxRetries = 10

// Memory provides a generic JSON object that's loaded from a DocumentStore.
type Memory struct {
	clear    func()
	dataType reflect.Type
	data     interface{}
	//version  Version

	MemoryStore MemoryStore
	Path        keys.Key
	Key         keys.Key
}

// Init performs an in-place initialization of a Memory.
func (m *Memory) Init(data interface{}, clear func(), store MemoryStore, path, key keys.Key) {
	m.clear = clear
	m.dataType = reflect.TypeOf(data)
	m.data = data
	m.MemoryStore = store
	m.Path = path
	m.Key = key
}

// LoadFromString loads this entity from the provided string.
func (m *Memory) LoadFromString(s string) (err error) {
	return jsonx.ParseString(s, m.data)
}

// Clear clears all data on this Memory.
func (m *Memory) Clear() {
	m.clear()
}

// Create creates this Memory if it doesn't already exist.
func (m *Memory) Create() (err error) {
	err = m.MemoryStore.HSet(
		m.Path,
		m.Key,
		options.WithSource(m.data),
	)
	return
}

//// CreateExpiry creates this Memory with an expiration value, if it doesn't already exist.
//func (m *Memory) CreateExpiry(expiry time.Duration) (err error) {
//	err = m.MemoryStore.Set(
//		m.Path,
//		m.Key,
//		options.WithSource(m.data),
//		options.WithTTL(expiry),
//	)
//	return
//}

// Load loads this Memory from its store if it exists.
func (m *Memory) Load() (err error) {
	m.clear()

	err = m.MemoryStore.HGet(
		m.Path,
		m.Key,
		options.WithDestination(m.data),
	)
	return
}

//// LoadAndTouch loads this Memory from its store if it exists while
//// simultaneously updating any Expiry setting.
//func (m *Memory) LoadAndTouch(expiry time.Duration) (err error) {
//	m.clear()
//
//	m.version, err = m.MemoryStore.Get(
//		m.Key,
//		WithDestination(m.data),
//		WithTTL(expiry),
//	)
//	return
//}

// Save saves this Memory to the database if it's based on the latest version that Redis knows about.
func (m *Memory) Save() (err error) {
	err = m.MemoryStore.HSet(
		m.Path,
		m.Key,
		options.WithSource(m.data),
	)
	return
}

//// SaveExpiry saves this Memory to the database, with a new expiration value,
//// if it's based on the latest version that Redis knows about.
//func (m *Memory) SaveExpiry(expiry time.Duration) (err error) {
//	m.version, err = m.MemoryStore.Set(
//		m.Key,
//		WithSource(m.data),
//		WithVersion(m.version),
//		WithTTL(expiry),
//	)
//	return
//}

// doUpdate executes the provided function until it returns true and
// the Memory can be stored using the provided update function in the database without issue.
// If too many attempts are made then the function fails with ErrTooManyRetries.
//func (m *Memory) doUpdate(f func() bool, u func() error) error {

//for r := 0; r < MaxRetries; r++ {
//	if f() {
//		if err := u(); err == nil {
//			return nil
//		} else {
//			// SNICHOLS: This random sleep is to help with updating the same Entity multiple times.  See
//			// https://github.com/89trillion/platform/services/issues/136 for more details.
//			time.Sleep(time.Millisecond * time.Duration(rand.Float32()*float32(r+1)*5))
//			if err := m.Load(); err != nil {
//				// a failed load is a real error
//				return err
//			}
//		}
//	} else {
//		return errors.ErrUpdateLogicFailed
//	}
//}
//
//	return errors.ErrTooManyRetries
//}

// Invokes doUpdate with Save() as the update function.
func (m *Memory) Update(f func() bool) error {
	//TODO add redis lock setnx
	if err := m.Load(); err != nil {
		// a failed load is a real error
		return err
	} else if f() {
		return m.Save()
	}
	return nil
}

//// Invokes doUpdate with SaveExpiry() as the update function.
//func (m *Memory) UpdateExpiry(f func() bool, expiry time.Duration) error {
//	return m.doUpdate(f, func() error {
//		return m.SaveExpiry(expiry)
//	})
//
//}

//// Invokes document sub-mutation to add entry to end of array.
//func (m *Memory) PushBack(key string, path string, item interface{}) error {
//	return m.MemoryStore.PushBack(Key{value: key}, path, item)
//}

// Remove removes the Memory from the database.
func (m *Memory) Remove() error {
	return m.MemoryStore.HDel(m.Path, m.Key)
}
