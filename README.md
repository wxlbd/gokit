# GoKit

GoKit 是一个功能丰富的 Go 语言工具库，提供了常用的工具函数和组件。

## 项目结构

```
gokit/
├── protocols/     # 协议实现
│   ├── modbus/   # Modbus 协议实现
│   ├── ymodem/   # YModem 协议实现
│   └── at/       # AT 指令实现
├── ds/           # 数据结构
│   └── bimap/    # 双向映射实现
├── utils/        # 通用工具
│   ├── bytex/    # 字节处理工具
│   ├── intx/     # 整数处理工具
│   └── genericx/ # 泛型工具
└── middleware/   # 中间件组件
    ├── jwt/      # JWT 认证
    ├── logger/   # 日志处理
    └── singleflight/ # 请求合并
```

## 使用说明

```go
import "github.com/wxlbd/gokit/utils/bytex"
import "github.com/wxlbd/gokit/protocols/modbus"
// ... 其他包的导入
```

## 特性

- 模块化设计，每个包都是独立的，可以单独使用
- 完整的测试覆盖
- 详细的文档和示例
- 持续维护和更新

## 许可证

MIT License
