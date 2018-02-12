# tmt-go-sdk
golang for 腾讯机器翻译sdk，https://cloud.tencent.com/document/product/551/7380

## Example
```go
package main

import (
	"github.com/dreamCodeMan/tmt-go-sdk"
	"fmt"
)

func main() {
	client := translate.New("AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA", "Gu5t9xGARNpq86cd98joQYCN3Cozk1qA", "gz")
	result, err := client.Do("你好")
	fmt.Println(result.TargetText, err)
}
```
