package modbus

// crc.go 实现了 Modbus RTU 协议中使用的 CRC-16 校验算法

// CRC16 计算给定数据的 Modbus CRC-16 校验值
// Modbus RTU 使用的是 CRC-16-ANSI 算法，多项式为 x^16 + x^15 + x^2 + 1 (0xA001)
func CRC16(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc = crc >> 1
			}
		}
	}
	return crc
}

// AppendCRC16 计算给定数据的 CRC-16 校验值并以小端序附加到数据末尾
func AppendCRC16(data []byte) []byte {
	crc := CRC16(data)
	// Modbus RTU 中 CRC 是以小端序方式添加的 (低字节在前，高字节在后)
	result := make([]byte, len(data)+2)
	copy(result, data)
	result[len(data)] = byte(crc)        // 低字节
	result[len(data)+1] = byte(crc >> 8) // 高字节
	return result
}

// CheckCRC16 验证数据的 CRC-16 校验值是否正确
// 输入数据必须包含 CRC 校验值（数据的最后两个字节）
func CheckCRC16(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	// 分离数据和 CRC
	actualData := data[:len(data)-2]
	receivedCRC := uint16(data[len(data)-2]) | uint16(data[len(data)-1])<<8

	// 计算 CRC
	calculatedCRC := CRC16(actualData)

	return calculatedCRC == receivedCRC
}

// ExtractWithoutCRC16 从包含 CRC 的数据中提取原始数据（不包含 CRC）
func ExtractWithoutCRC16(data []byte) []byte {
	if len(data) < 2 {
		return data
	}
	return data[:len(data)-2]
}
