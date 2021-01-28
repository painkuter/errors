package errors

import (
	"strconv"
)

// Property ...
type Property struct {
	_ interface{} // fix issue with empty struct always has same addr
}

// PropertyValue ...
type PropertyValue struct {
	value interface{}
}

// PropertyMap ...
type PropertyMap map[*Property]*PropertyValue

// Predefined properties
var (
	ServiceName  = NewProperty()
	GRPCCode     = NewProperty()
	GRPCMessage  = NewProperty()
	HTTPCode     = NewProperty()
	HTTPBody     = NewProperty()
	ResponseCode = NewProperty()
)

// NewProperty ...
func NewProperty() *Property {
	return &Property{}
}

// Value ...
func (v *PropertyValue) Value() interface{} {
	if v == nil || v.value == nil {
		return nil
	}
	return v.value
}

// String ...
func (v *PropertyValue) String() string {
	if v == nil {
		return ""
	}
	return toString(v.value)
}

// Int ...
func (v *PropertyValue) Int() int {
	if v == nil || v.value == nil {
		return 0
	}
	switch typed := v.value.(type) {
	case stringer:
		if val, err := strconv.Atoi(typed.String()); err == nil {
			return val
		}
	case string:
		if val, err := strconv.Atoi(typed); err == nil {
			return val
		}
	case int:
		return typed
	case int8:
		return int(typed)
	case int16:
		return int(typed)
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case uint:
		return int(typed)
	case uint8:
		return int(typed)
	case uint16:
		return int(typed)
	case uint32:
		return int(typed)
	case uint64:
		return int(typed)
	}
	return 0
}
