## elconv
- interface{}转换成指定类型
- convert interface{} to other type
- interface{}を特定の型に変換する
```go
import "github.com/el-ideal-ideas/ellib/elconv"

elconv.AsBool(v interface{}) bool
elconv.AsFloat64(v interface{}) float64
elconv.AsFloat32(v interface{}) float32
elconv.AsInt(v interface{}) int
elconv.AsInt64(v interface{}) int64
elconv.AsStr(v interface{}) string
elconv.AsUint(v interface{}) uint
elconv.AsUint64(v interface{}) uint64
```