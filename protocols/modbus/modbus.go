package modbus

import (
	"errors"
	"fmt"
)

// 常见的 Modbus 功能码
const (
	// 位操作功能码
	FuncReadCoils          byte = 0x01 // 读取线圈状态
	FuncReadDiscreteInputs byte = 0x02 // 读取离散输入状态
	FuncWriteSingleCoil    byte = 0x05 // 写单个线圈
	FuncWriteMultipleCoils byte = 0x0F // 写多个线圈

	// 字操作功能码
	FuncReadHoldingRegisters   byte = 0x03 // 读取保持寄存器
	FuncReadInputRegisters     byte = 0x04 // 读取输入寄存器
	FuncWriteSingleRegister    byte = 0x06 // 写单个寄存器
	FuncWriteMultipleRegisters byte = 0x10 // 写多个寄存器

	// 其他功能码
	FuncReadExceptionStatus byte = 0x07 // 读取异常状态
	FuncDiagnostic          byte = 0x08 // 诊断
	FuncGetCommEventCounter byte = 0x0B // 获取通信事件计数
	FuncGetCommEventLog     byte = 0x0C // 获取通信事件日志
)

// 异常码
const (
	ExcIllegalFunction                    byte = 0x01 // 非法功能
	ExcIllegalDataAddress                 byte = 0x02 // 非法数据地址
	ExcIllegalDataValue                   byte = 0x03 // 非法数据值
	ExcServerDeviceFailure                byte = 0x04 // 服务器设备故障
	ExcAcknowledge                        byte = 0x05 // 确认
	ExcServerDeviceBusy                   byte = 0x06 // 服务器设备忙
	ExcMemoryParityError                  byte = 0x08 // 内存奇偶校验错误
	ExcGatewayPathUnavailable             byte = 0x0A // 网关路径不可用
	ExcGatewayTargetDeviceFailedToRespond byte = 0x0B // 网关目标设备无响应
)

// ModbusError 表示 Modbus 错误
type ModbusError struct {
	FunctionCode  byte
	ExceptionCode byte
}

// Error 实现 error 接口
func (e *ModbusError) Error() string {
	return fmt.Sprintf("modbus: function code %#.2x exception: %#.2x", e.FunctionCode, e.ExceptionCode)
}

// IsError 检查响应是否为错误
func IsError(functionCode byte) bool {
	return functionCode&0x80 != 0
}

// ParseError 从响应中解析错误
func ParseError(functionCode, exceptionCode byte) error {
	return &ModbusError{
		FunctionCode:  functionCode &^ 0x80,
		ExceptionCode: exceptionCode,
	}
}

// RTUFrame 表示 Modbus RTU 帧
type RTUFrame struct {
	SlaveID      byte
	FunctionCode byte
	Data         []byte
}

// ErrCRCMismatch 表示 CRC 校验不匹配
var ErrCRCMismatch = errors.New("modbus: CRC mismatch")

// ErrResponseTooShort 表示响应帧太短
var ErrResponseTooShort = errors.New("modbus: response too short")

// ErrInvalidSlaveID 表示从站 ID 不匹配
var ErrInvalidSlaveID = errors.New("modbus: invalid slave ID")

// ErrInvalidFunction 表示功能码不匹配
var ErrInvalidFunction = errors.New("modbus: invalid function")

// ErrInvalidLength 表示长度无效
var ErrInvalidLength = errors.New("modbus: invalid length")
