# MySQL DATETIME 时区问题及解决方案

## 1. 问题背景

在数据库设计中，处理日期和时间是一个常见的需求，但时区问题常常会导致数据不一致和各种难以排查的 bug。

MySQL 的 `DATETIME` 类型存储 `'YYYY-MM-DD HH:MM:SS'` 格式的时间，但它**不包含任何时区信息**。这意味着，当一个值（如 `2023-10-27 10:00:00`）存入数据库时，我们无法仅从值本身判断它代表的是哪个时区的时间。当系统涉及多个时区时，这很容易引发混乱。

## 2. 推荐解决方案：存储 UTC，本地转换

为了彻底解决时区问题，业界公认的最佳实践是：

**后端统一使用 UTC 时间，将时间存入数据库，然后在展示层（前端或客户端）根据用户所在时区进行本地化转换。**

### 核心原则
1.  **数据库存储 UTC 时间**: 数据库中的 `DATETIME` 字段统一存储协调世界时 (UTC)，保证时间的绝对性和唯一性。
2.  **应用程序负责转换**:
    *   **写入时**: 应用程序将所有接收到的时间转换为 UTC 时间再存入数据库。
    *   **读取时**: 应用程序从数据库读取 UTC 时间后，根据业务需要转换为目标时区的时间再返回。
3.  **服务器时区统一**: 建议将数据库和应用服务器的系统时区都设置为 UTC，以避免环境差异。

## 3. Go (GORM) 实践示例

以下是在 Go 和 GORM 中实践此原则的两种方案。
在连接时， 需要将loc=UTC设置到连接dsn中，以确保gorm在执行查询等时， 自动将sql变量转为正确的值。

### 方案一：使用 `time.Time` 并配置 GORM (基础)

这是最简单的方法，适用于 `CreatedAt` 和 `UpdatedAt` 等由 GORM 自动管理的字段。

**1. GORM 模型定义:**
```go
import "time"

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string
    CreatedAt time.Time // GORM会自动处理
    UpdatedAt time.Time // GORM会自动处理
}
```

**2. GORM 初始化配置:**
通过覆盖 GORM 的 `NowFunc`，可以确保自动填充的时间戳是 UTC。

```go
import (
    "gorm.io/gorm"
    "time"
)

// 在初始化 GORM DB 实例时进行配置
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    NowFunc: func() time.Time {
        // 全局使用 UTC 时间
        return time.Now().UTC()
    },
})
```
**优点**: 简单快捷。
**缺点**: 对于自定义时间字段，需要开发者手动确保存入的是 UTC 时间 (`time.Now().UTC()`)，容易遗漏。

### 方案二：自定义 `UTCTime` 类型 (高级)

为了获得更强的类型安全和代码清晰度，我们可以创建一个自定义的 `UTCTime` 类型，它能自动处理与数据库之间的 UTC 时间转换。

**1. 定义 `UTCTime` 类型**

这个类型嵌入了 `time.Time` 并实现了 `sql.Scanner` 和 `driver.Valuer` 接口，从而接管 GORM 的读写过程。

```go
import (
    "database/sql/driver"
    "fmt"
    "time"
)

// UTCTime 是一个自定义类型，确保与数据库交互时使用 UTC 时间
type UTCTime struct {
	time.Time
}

// MarshalJSON 实现 json.Marshaler 接口，用于JSON序列化
func (t UTCTime) MarshalJSON() ([]byte, error) {
    if t.IsZero() {
        return []byte("null"), nil
    }
	// 格式化为标准的 RFC3339 格式，例如 "2006-01-02T15:04:05Z"
	formatted := fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))
	return []byte(formatted), nil
}

// Value 实现 driver.Valuer 接口，在写入数据库时被调用
func (t UTCTime) Value() (driver.Value, error) {
	// 如果时间是零值，则存入 NULL
	if t.IsZero() {
		return nil, nil
	}
	// 确保存入的是UTC时间
	return t.Time.UTC(), nil
}

// Scan 实现 sql.Scanner 接口，在从数据库读取时被调用
func (t *UTCTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{} // 存入零值
		return nil
	}

	var err error
	switch v := value.(type) {
	case time.Time:
		// 数据库驱动已将其解析为 time.Time (推荐在DSN中设置 parseTime=true&loc=UTC)
		t.Time = v.UTC()
	case []byte:
		// 可能是 "YYYY-MM-DD HH:MM:SS" 格式的字节数组
		t.Time, err = time.Parse("2006-01-02 15:04:05", string(v))
        if err == nil {
            t.Time = t.Time.UTC()
        }
	case string:
		// 可能是 "YYYY-MM-DD HH:MM:SS" 格式的字符串
		t.Time, err = time.Parse("2006-01-02 15:04:05", v)
        if err == nil {
            t.Time = t.Time.UTC()
        }
	default:
		return fmt.Errorf("cannot scan type %T into UTCTime: %v", value, value)
	}

	return err
}
```

**2. 在 GORM 模型中使用 `UTCTime`**

现在，可以在模型中直接使用 `UTCTime`，它会自动处理 UTC 转换。

```go
type MiniGameBetPreferenceTag struct {
	ID        uint64         `gorm:"primaryKey"`
	// ... 其他字段
	CreatedAt *UTCTime       `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt *UTCTime       `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**优点**:
- **类型安全**: 强制所有时间都以 UTC 方式处理，避免人为错误。
- **代码清晰**: 业务代码无需再关心时区转换，模型定义即文档。
- **自动转换**: 无论是读是写，都能确保时区正确。

## 总结

对于新项目或需要重构的项目，强烈推荐使用方案二（自定义 `UTCTime` 类型），因为它从根本上解决了问题，并提高了代码的健壮性和可维护性。
