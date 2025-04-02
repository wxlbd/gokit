# Modbus RTU 协议实现

这个包提供了 Modbus RTU 协议的实现，包括请求帧生成和响应帧解析功能。该实现支持常见的 Modbus 功能码，包括位操作（线圈操作）和字操作（寄存器操作）。

## 功能特性

- 支持常见的 Modbus 功能码
- 提供请求帧生成和响应帧解析
- 内置 CRC-16 校验计算
- 异常处理和错误码解析
- 支持自定义传输接口
- 类型安全的 API

## 支持的功能码

### 位操作功能码
- `0x01`: 读取线圈状态 (Read Coils)
- `0x02`: 读取离散输入状态 (Read Discrete Inputs)
- `0x05`: 写单个线圈 (Write Single Coil)
- `0x0F`: 写多个线圈 (Write Multiple Coils)

### 字操作功能码
- `0x03`: 读取保持寄存器 (Read Holding Registers)
- `0x04`: 读取输入寄存器 (Read Input Registers)
- `0x06`: 写单个寄存器 (Write Single Register)
- `0x10`: 写多个寄存器 (Write Multiple Registers)

## 使用示例

### 初始化客户端

```go
// 客户端需要一个实现了 io.ReadWriter 接口的通信介质
client := modbus.NewClient(transport, 1) // 从站 ID = 1
client.SetTimeout(time.Second * 2)       // 设置超时时间
client.SetInterFrameDelay(time.Millisecond * 100) // 设置帧间延时
```

### 读取线圈状态

```go
// 从地址 0 开始读取 16 个线圈
coils, err := client.ReadCoils(0, 16)
if err != nil {
    log.Printf("读取线圈状态失败: %v", err)
    return
}

for i, status := range coils {
    fmt.Printf("线圈 %d: %v\n", i, status)
}
```

### 读取保持寄存器

```go
// 从地址 0 开始读取 10 个寄存器
registers, err := client.ReadHoldingRegisters(0, 10)
if err != nil {
    log.Printf("读取保持寄存器失败: %v", err)
    return
}

for i, reg := range registers {
    fmt.Printf("寄存器 %d: 0x%04X (%d)\n", i, reg, reg)
}
```

### 写单个线圈

```go
// 向地址 10 写入值 ON
err = client.WriteSingleCoil(10, true)
if err != nil {
    log.Printf("写线圈失败: %v", err)
    return
}
```

### 写单个寄存器

```go
// 向地址 1 写入值 123
err = client.WriteSingleRegister(1, 123)
if err != nil {
    log.Printf("写寄存器失败: %v", err)
    return
}
```

### 写多个线圈

```go
// 从地址 20 开始写入 5 个线圈状态
err = client.WriteMultipleCoils(20, []bool{true, false, true, false, true})
if err != nil {
    log.Printf("写多个线圈失败: %v", err)
    return
}
```

### 写多个寄存器

```go
// 从地址 10 开始写入 3 个寄存器值
err = client.WriteMultipleRegisters(10, []uint16{111, 222, 333})
if err != nil {
    log.Printf("写多个寄存器失败: %v", err)
    return
}
```

## 异常处理

当 Modbus 设备返回异常响应时，客户端会返回 `ModbusError` 类型的错误，其中包含功能码和异常码：

```go
// 尝试读取超出范围的寄存器
_, err := client.ReadHoldingRegisters(65535, 10)
if err != nil {
    if modbusErr, ok := err.(*modbus.ModbusError); ok {
        fmt.Printf("Modbus 错误 - 功能码: 0x%02X, 异常码: 0x%02X\n", 
            modbusErr.FunctionCode, modbusErr.ExceptionCode)
    }
}
```

## 自定义传输接口

你可以使用任何实现了 `io.ReadWriter` 接口的类型作为传输介质，例如串口、TCP 连接或自定义实现。

```go
// 串口示例
config := &serial.Config{
    Name: "/dev/ttyUSB0",
    Baud: 9600,
    // 其他配置...
}
port, err := serial.OpenPort(config)
if err != nil {
    log.Fatal(err)
}
client := modbus.NewClient(port, 1)

// TCP 示例
conn, err := net.Dial("tcp", "192.168.1.100:502")
if err != nil {
    log.Fatal(err)
}
client := modbus.NewClient(conn, 1)
```

## 低级 API

该包也提供了低级 API，允许直接生成请求帧和解析响应帧：

```go
// 生成读取保持寄存器的请求帧
request := modbus.NewReadHoldingRegistersRequest(1, 0, 10)

// 解析读取保持寄存器的响应
registers, err := modbus.ParseReadRegistersResponse(response, 1, modbus.FuncReadHoldingRegisters)
``` 