# Slice.StartsWith

```go
func (s Slice) StartsWith(prefix []T) bool
```

StartsWith returns true if slice starts with the given prefix slice. If prefix is empty, it returns true.

Generic types: T.

## Structs

| Struct | T type |
| ------ | ------ |
| SliceBool | bool |
| SliceByte | byte |
| SliceString | string |
| SliceFloat32 | float32 |
| SliceFloat64 | float64 |
| SliceInt | int |
| SliceInt8 | int8 |
| SliceInt16 | int16 |
| SliceInt32 | int32 |
| SliceInt64 | int64 |
| SliceUint | uint |
| SliceUint8 | uint8 |
| SliceUint16 | uint16 |
| SliceUint32 | uint32 |
| SliceUint64 | uint64 |
| SliceInterface | interface{} |

## Source

```go
// StartsWith returns true if slice starts with the given prefix slice.
// If prefix is empty, it returns true.
func (s Slice) StartsWith(prefix []T) bool {
	if len(prefix) > len(s.Data) {
		return false
	}
	for i, el := range prefix {
		if el != s.Data[i] {
			return false
		}
	}
	return true
}
```
