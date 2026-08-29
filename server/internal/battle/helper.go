// Package battle 辅助：JSON 封装（便于替换/测试）与 gctx 快捷方式。
package battle

import (
	"encoding/json"
	"context"
)

func jsonMarshal(v interface{}) ([]byte, error)  { return json.Marshal(v) }
func jsonUnmarshal(b string, v interface{}) error { return json.Unmarshal([]byte(b), v) }
func gctxNew() context.Context                   { return context.Background() }
