# ID 生成器

提供分佈式ID生成方案，基於Twitter的Snowflake算法。

## 特性

- **全局唯一性**：分佈式環境下保證ID唯一
- **有序性**：生成的ID按時間遞增
- **高效率**：單機可生成超過400萬個ID/秒
- **線程安全**：支持併發訪問
- **易擴展**：支持多個數據中心和機器

## Snowflake算法結構

```
┌─────────────────────────────────────────────────────────────────┐
│0│         41 bits timestamp          │5 bits│5 bits│  12 bits   │
│ │     (milliseconds since epoch)      │  DC  │Machine│ sequence  │
└─────────────────────────────────────────────────────────────────┘
```

- **1 bit**：符號位，始終為0，確保生成的ID為正數
- **41 bits**：時間戳，精確到毫秒，可使用69年
- **5 bits**：數據中心ID (0-31)
- **5 bits**：機器ID (0-31)
- **12 bits**：序列號，同毫秒內的序列 (0-4095)

## 使用方法

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/FatSong207/common/id"
)

func main() {
    // 建立生成器 (數據中心1, 機器1)
    sf, err := id.NewSnowflake(1, 1)
    if err != nil {
        panic(err)
    }
    
    // 生成int64 ID
    id64 := sf.Generate()
    fmt.Println(id64)  // 輸出類似：6762906437922086913
    
    // 生成字符串ID
    idStr := sf.GenerateString()
    fmt.Println(idStr)  // 輸出類似：2Hn5Y8Kp3Q
}
```

### 實例

```go
// 單個實例
sf, _ := id.NewSnowflake(1, 1)

// 並發生成ID
for i := 0; i < 100; i++ {
    id := sf.Generate()
    // 使用ID
}
```

## API

### NewSnowflake(datacenter, machine int64) (*Snowflake, error)

建立一個新的Snowflake ID生成器。

**參數**：
- `datacenter`: 數據中心ID，範圍0-31
- `machine`: 機器ID，範圍0-31

**返回**：
- 生成器實例和可能的錯誤

### Generate() int64

生成一個int64類型的唯一ID。

### GenerateString() string

生成一個Base62編碼的字符串ID。

### Int64ToA(num int64) string

將int64轉換為Base62編碼的字符串。

### AToInt64(s string) int64

將Base62編碼的字符串轉換為int64。

## 性能

單機性能測試結果：

```
BenchmarkGenerate-8         50000000    23.2 ns/op
BenchmarkGenerateString-8   30000000    35.8 ns/op
```

## 注意事項

1. 同一datacenter/machine組合應該只在一個應用實例中使用
2. 時間戳以2022-01-01為基準，可使用69年
3. 時鐘回撥會導致ID重複，需要確保系統時間同步
