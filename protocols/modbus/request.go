package modbus

import (
	"encoding/binary"
)

// 请求帧生成器 - 位操作相关功能

// NewReadCoilsRequest 创建读取线圈状态请求
func NewReadCoilsRequest(slaveID byte, startAddress uint16, quantity uint16) []byte {
	if quantity < 1 || quantity > 2000 {
		quantity = 1 // 默认读取一个线圈
	}

	data := make([]byte, 5)
	data[0] = slaveID
	data[1] = FuncReadCoils
	binary.BigEndian.PutUint16(data[2:4], startAddress)
	data[4] = byte(quantity)
	if quantity > 255 {
		data = append(data[:4], byte(quantity>>8), byte(quantity))
	}

	return AppendCRC16(data)
}

// NewReadDiscreteInputsRequest 创建读取离散输入状态请求
func NewReadDiscreteInputsRequest(slaveID byte, startAddress uint16, quantity uint16) []byte {
	if quantity < 1 || quantity > 2000 {
		quantity = 1 // 默认读取一个输入
	}

	data := make([]byte, 5)
	data[0] = slaveID
	data[1] = FuncReadDiscreteInputs
	binary.BigEndian.PutUint16(data[2:4], startAddress)
	data[4] = byte(quantity)
	if quantity > 255 {
		data = append(data[:4], byte(quantity>>8), byte(quantity))
	}

	return AppendCRC16(data)
}

// NewWriteSingleCoilRequest 创建写单个线圈请求
// 线圈状态: true = ON (0xFF00), false = OFF (0x0000)
func NewWriteSingleCoilRequest(slaveID byte, address uint16, value bool) []byte {
	data := make([]byte, 6)
	data[0] = slaveID
	data[1] = FuncWriteSingleCoil
	binary.BigEndian.PutUint16(data[2:4], address)

	if value {
		data[4] = 0xFF
		data[5] = 0x00
	} else {
		data[4] = 0x00
		data[5] = 0x00
	}

	return AppendCRC16(data)
}

// NewWriteMultipleCoilsRequest 创建写多个线圈请求
func NewWriteMultipleCoilsRequest(slaveID byte, startAddress uint16, values []bool) []byte {
	if len(values) < 1 || len(values) > 1968 {
		return nil // 无效的线圈数量
	}

	// 计算字节数
	byteCount := (len(values) + 7) / 8

	// 请求头: slaveID + 功能码 + 起始地址(2字节) + 线圈数量(2字节) + 字节数
	data := make([]byte, 6+1+byteCount)
	data[0] = slaveID
	data[1] = FuncWriteMultipleCoils
	binary.BigEndian.PutUint16(data[2:4], startAddress)
	binary.BigEndian.PutUint16(data[4:6], uint16(len(values)))
	data[6] = byte(byteCount)

	// 填充线圈状态
	for i, value := range values {
		byteIdx := 7 + i/8
		bitIdx := uint(i % 8)

		if value {
			data[byteIdx] |= 1 << bitIdx
		}
	}

	return AppendCRC16(data)
}

// 请求帧生成器 - 字操作相关功能

// NewReadHoldingRegistersRequest 创建读取保持寄存器请求
func NewReadHoldingRegistersRequest(slaveID byte, startAddress uint16, quantity uint16) []byte {
	if quantity < 1 || quantity > 125 {
		quantity = 1 // 默认读取一个寄存器
	}

	data := make([]byte, 6)
	data[0] = slaveID
	data[1] = FuncReadHoldingRegisters
	binary.BigEndian.PutUint16(data[2:4], startAddress)
	binary.BigEndian.PutUint16(data[4:6], quantity)

	return AppendCRC16(data)
}

// NewReadInputRegistersRequest 创建读取输入寄存器请求
func NewReadInputRegistersRequest(slaveID byte, startAddress uint16, quantity uint16) []byte {
	if quantity < 1 || quantity > 125 {
		quantity = 1 // 默认读取一个寄存器
	}

	data := make([]byte, 6)
	data[0] = slaveID
	data[1] = FuncReadInputRegisters
	binary.BigEndian.PutUint16(data[2:4], startAddress)
	binary.BigEndian.PutUint16(data[4:6], quantity)

	return AppendCRC16(data)
}

// NewWriteSingleRegisterRequest 创建写单个寄存器请求
func NewWriteSingleRegisterRequest(slaveID byte, address uint16, value uint16) []byte {
	data := make([]byte, 6)
	data[0] = slaveID
	data[1] = FuncWriteSingleRegister
	binary.BigEndian.PutUint16(data[2:4], address)
	binary.BigEndian.PutUint16(data[4:6], value)

	return AppendCRC16(data)
}

// NewWriteMultipleRegistersRequest 创建写多个寄存器请求
func NewWriteMultipleRegistersRequest(slaveID byte, startAddress uint16, values []uint16) []byte {
	if len(values) < 1 || len(values) > 123 {
		return nil // 无效的寄存器数量
	}

	// 字节数 = 寄存器数量 * 2
	byteCount := len(values) * 2

	// 请求头: slaveID + 功能码 + 起始地址(2字节) + 寄存器数量(2字节) + 字节数
	data := make([]byte, 6+1+byteCount)
	data[0] = slaveID
	data[1] = FuncWriteMultipleRegisters
	binary.BigEndian.PutUint16(data[2:4], startAddress)
	binary.BigEndian.PutUint16(data[4:6], uint16(len(values)))
	data[6] = byte(byteCount)

	// 填充寄存器值
	for i, value := range values {
		offset := 7 + i*2
		binary.BigEndian.PutUint16(data[offset:offset+2], value)
	}

	return AppendCRC16(data)
}
