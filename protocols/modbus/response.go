package modbus

import (
	"encoding/binary"
	"errors"
)

// 响应帧解析器 - 通用功能

// ValidateResponse 验证响应帧的基本有效性
// 检查 CRC、长度以及从站 ID 和功能码是否匹配
func ValidateResponse(response []byte, expectedSlaveID, expectedFunctionCode byte) error {
	// 检查长度是否足够
	if len(response) < 4 { // 至少需要从站 ID、功能码和 CRC (2字节)
		return ErrResponseTooShort
	}

	// 检查 CRC
	if !CheckCRC16(response) {
		return ErrCRCMismatch
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查从站 ID
	if frameData[0] != expectedSlaveID {
		return ErrInvalidSlaveID
	}

	// 检查功能码，如果高位为 1，则为异常响应
	if IsError(frameData[1]) {
		if len(frameData) < 3 {
			return ErrResponseTooShort
		}
		return ParseError(frameData[1], frameData[2])
	}

	// 检查功能码是否匹配
	if frameData[1] != expectedFunctionCode {
		return ErrInvalidFunction
	}

	return nil
}

// ParseReadBitsResponse 解析读取位状态（线圈或离散输入）的响应
// 适用于功能码 0x01 和 0x02
func ParseReadBitsResponse(response []byte, expectedSlaveID, expectedFunctionCode byte) ([]bool, error) {
	// 验证响应
	err := ValidateResponse(response, expectedSlaveID, expectedFunctionCode)
	if err != nil {
		return nil, err
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查字节计数
	if len(frameData) < 3 { // 从站 ID + 功能码 + 字节计数
		return nil, ErrResponseTooShort
	}

	byteCount := int(frameData[2])
	if len(frameData) < 3+byteCount {
		return nil, ErrResponseTooShort
	}

	// 将字节转换为位状态
	bitData := frameData[3 : 3+byteCount]

	// 由于 Modbus 协议中未指定总位数，所以我们返回所有位
	// 调用者需要根据请求的位数量来取用适当数量的位
	result := make([]bool, byteCount*8)

	for i := 0; i < byteCount; i++ {
		for j := 0; j < 8; j++ {
			idx := i*8 + j
			result[idx] = (bitData[i] & (1 << uint(j))) != 0
		}
	}

	return result, nil
}

// ParseReadRegistersResponse 解析读取寄存器（保持寄存器或输入寄存器）的响应
// 适用于功能码 0x03 和 0x04
func ParseReadRegistersResponse(response []byte, expectedSlaveID, expectedFunctionCode byte) ([]uint16, error) {
	// 验证响应
	err := ValidateResponse(response, expectedSlaveID, expectedFunctionCode)
	if err != nil {
		return nil, err
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查字节计数
	if len(frameData) < 3 { // 从站 ID + 功能码 + 字节计数
		return nil, ErrResponseTooShort
	}

	byteCount := int(frameData[2])
	if len(frameData) < 3+byteCount {
		return nil, ErrResponseTooShort
	}

	// 字节计数应该是偶数（每个寄存器 2 字节）
	if byteCount%2 != 0 {
		return nil, ErrInvalidLength
	}

	// 解析寄存器值
	registerCount := byteCount / 2
	registerData := frameData[3 : 3+byteCount]
	result := make([]uint16, registerCount)

	for i := 0; i < registerCount; i++ {
		offset := i * 2
		result[i] = binary.BigEndian.Uint16(registerData[offset : offset+2])
	}

	return result, nil
}

// ParseWriteSingleCoilResponse 解析写单个线圈的响应
func ParseWriteSingleCoilResponse(response []byte, expectedSlaveID byte, expectedAddress uint16, expectedValue bool) error {
	// 验证响应
	err := ValidateResponse(response, expectedSlaveID, FuncWriteSingleCoil)
	if err != nil {
		return err
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查响应长度
	if len(frameData) < 6 { // 从站 ID + 功能码 + 地址(2字节) + 值(2字节)
		return ErrResponseTooShort
	}

	// 检查地址是否匹配
	address := binary.BigEndian.Uint16(frameData[2:4])
	if address != expectedAddress {
		return errors.New("modbus: address mismatch in response")
	}

	// 检查值是否匹配
	value := frameData[4] == 0xFF && frameData[5] == 0x00
	expectedOff := frameData[4] == 0x00 && frameData[5] == 0x00

	if value != expectedValue && !(expectedValue == false && expectedOff) {
		return errors.New("modbus: value mismatch in response")
	}

	return nil
}

// ParseWriteSingleRegisterResponse 解析写单个寄存器的响应
func ParseWriteSingleRegisterResponse(response []byte, expectedSlaveID byte, expectedAddress, expectedValue uint16) error {
	// 验证响应
	err := ValidateResponse(response, expectedSlaveID, FuncWriteSingleRegister)
	if err != nil {
		return err
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查响应长度
	if len(frameData) < 6 { // 从站 ID + 功能码 + 地址(2字节) + 值(2字节)
		return ErrResponseTooShort
	}

	// 检查地址是否匹配
	address := binary.BigEndian.Uint16(frameData[2:4])
	if address != expectedAddress {
		return errors.New("modbus: address mismatch in response")
	}

	// 检查值是否匹配
	value := binary.BigEndian.Uint16(frameData[4:6])
	if value != expectedValue {
		return errors.New("modbus: value mismatch in response")
	}

	return nil
}

// ParseWriteMultipleCoilsResponse 解析写多个线圈的响应
func ParseWriteMultipleCoilsResponse(response []byte, expectedSlaveID byte, expectedAddress uint16, expectedQuantity uint16) error {
	// 验证响应
	err := ValidateResponse(response, expectedSlaveID, FuncWriteMultipleCoils)
	if err != nil {
		return err
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查响应长度
	if len(frameData) < 6 { // 从站 ID + 功能码 + 地址(2字节) + 数量(2字节)
		return ErrResponseTooShort
	}

	// 检查地址是否匹配
	address := binary.BigEndian.Uint16(frameData[2:4])
	if address != expectedAddress {
		return errors.New("modbus: address mismatch in response")
	}

	// 检查数量是否匹配
	quantity := binary.BigEndian.Uint16(frameData[4:6])
	if quantity != expectedQuantity {
		return errors.New("modbus: quantity mismatch in response")
	}

	return nil
}

// ParseWriteMultipleRegistersResponse 解析写多个寄存器的响应
func ParseWriteMultipleRegistersResponse(response []byte, expectedSlaveID byte, expectedAddress uint16, expectedQuantity uint16) error {
	// 验证响应
	err := ValidateResponse(response, expectedSlaveID, FuncWriteMultipleRegisters)
	if err != nil {
		return err
	}

	// 提取响应内容（不含 CRC）
	frameData := response[:len(response)-2]

	// 检查响应长度
	if len(frameData) < 6 { // 从站 ID + 功能码 + 地址(2字节) + 数量(2字节)
		return ErrResponseTooShort
	}

	// 检查地址是否匹配
	address := binary.BigEndian.Uint16(frameData[2:4])
	if address != expectedAddress {
		return errors.New("modbus: address mismatch in response")
	}

	// 检查数量是否匹配
	quantity := binary.BigEndian.Uint16(frameData[4:6])
	if quantity != expectedQuantity {
		return errors.New("modbus: quantity mismatch in response")
	}

	return nil
}
