# GORM Gen 代碼生成工具

## 概述

本包裝提供了簡化的 [GORM Gen](https://gorm.io/gen) 代碼生成接口。GORM Gen 是一個強大的代碼生成工具，能夠根據你的數據模型自動生成類型安全的查詢代碼，提供編譯時檢查和自動完成支持。

## 快速開始

### 1. 定義你的數據模型

```go
package model

type User struct {
	ID    int    `gorm:"primaryKey"`
	Name  string
	Email string
}

type Product struct {
	ID    int     `gorm:"primaryKey"`
	Name  string
	Price float64
}
```

### 2. 創建代碼生成腳本

創建一個 `main.go` 文件（例如在 `cmd/gen/main.go`）：

```go
package main

import (
	"your-module/db/gorm/gen"
	"your-module/model"
)

func main() {
	// 生成查詢代碼到 query 目錄
	gen.G("query", &model.User{}, &model.Product{})
}
```

### 3. 執行代碼生成

```bash
go run cmd/gen/main.go
```

這會在 `query` 目錄下生成查詢模型文件。

## 使用生成的查詢代碼

### 初始化

```go
import (
	"your-module/db/gorm"
	"your-module/query"
)

func init() {
	db := gorm.NewDB(dsn)
	query.SetDefault(db.Db)
    query.Use(db) // 也可以
}
```

### 查詢示例

```go
import "your-module/query"

// 簡單查詢
user, err := query.User.Where(query.User.Name.Eq("Alice")).First()

// 複雜查詢
users, err := query.User.
	Where(query.User.Email.Like("%@gmail.com")).
	Order(query.User.ID.Desc()).
	Limit(10).
	Find()

// 統計
count, err := query.User.Where(query.User.Age.Gt(18)).Count()
```

## 函數簽名

```go
func G(output string, models ...any)
```

### 參數

- **output** (string): 生成文件的輸出目錄路徑
- **models** (...any): 需要生成查詢代碼的數據模型，可傳入多個

### 生成配置

當前使用的配置：

```go
gen.Config{
	OutPath: output,
	Mode: gen.WithoutContext |         // 不生成 context 版本
	      gen.WithDefaultQuery |       // 生成默認查詢實例
	      gen.WithQueryInterface,      // 生成查詢接口
}
```

## 工作流程建議

### 開發階段

1. **修改數據模型** - 更新 struct 定義
2. **重新生成代碼** - 執行生成腳本
3. **更新查詢代碼** - 使用新生成的查詢方法

### 最佳實踐

- ✅ 將生成腳本放在 `cmd/gen/` 或 `scripts/` 目錄下
- ✅ 每次修改模型後重新生成
- ✅ 生成的代碼應該納入版本控制（commit 到 git）
- ✅ 為不同的服務分開執行生成腳本

### 示例項目結構

```
project/
├── cmd/
│   └── gen/
│       └── main.go           # 代碼生成腳本
├── model/
│   ├── user.go               # User 模型
│   └── product.go            # Product 模型
├── query/
│   ├── gen.go                # 生成的代碼
│   ├── user.gen.go           # User 查詢模型
│   └── product.gen.go        # Product 查詢模型
└── db/
    └── init.go               # 初始化查詢實例
```

## 常見問題

### Q: 生成的文件應該提交到 git 嗎？

**A:** 是的。生成的代碼應該提交到版本控制系統，這樣其他開發者不需要在本地重新生成。

### Q: 如何為不同的數據庫生成不同的代碼？

**A:** 在生成腳本中傳入 `output` 參數即可：

```go
gen.G("query/mysql", models...)    // MySQL
gen.G("query/postgres", models...) // PostgreSQL
```

### Q: 生成的代碼可以自定義嗎？

**A:** 基礎配置由本包裝提供，如需高級配置，可直接使用 `gorm.io/gen` 包。詳見 [GORM Gen 文檔](https://gorm.io/gen)。

## 下一步

- 查看 [GORM Gen 文檔](https://gorm.io/gen) 了解更多高級用法
- 探索生成代碼的查詢方法和功能
- 集成到 CI/CD 流程中自動化代碼生成