package modbus

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// 测试 CRC 计算功能
func TestCRC16(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected uint16
	}{
		{
			name:     "空数据",
			data:     []byte{},
			expected: 0xFFFF,
		},
		{
			name:     "读单个寄存器",
			data:     []byte{0x01, 0x03, 0x00, 0x6B, 0x00, 0x01},
			expected: 0xD6F5, // Modbus CRC-16
		},
		{
			name:     "写单个寄存器",
			data:     []byte{0x01, 0x06, 0x00, 0x01, 0x00, 0x03},
			expected: 0x0B98, // Modbus CRC-16
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CRC16(tt.data)
			if result != tt.expected {
				t.Errorf("CRC16() = %04X, want %04X", result, tt.expected)
			}
		})
	}
}

// 测试 AppendCRC16 和 CheckCRC16 功能
func TestAppendAndCheckCRC16(t *testing.T) {
	data := []byte{0x01, 0x03, 0x00, 0x6B, 0x00, 0x01}

	// 计算并附加 CRC
	withCRC := AppendCRC16(data)

	// 验证长度
	if len(withCRC) != len(data)+2 {
		t.Errorf("AppendCRC16() returned data with incorrect length: got %d, want %d", len(withCRC), len(data)+2)
	}

	// 验证 CRC 值
	crc := uint16(withCRC[len(data)]) | uint16(withCRC[len(data)+1])<<8
	expected := CRC16(data)
	if crc != expected {
		t.Errorf("AppendCRC16() appended incorrect CRC: got %04X, want %04X", crc, expected)
	}

	// 验证 CheckCRC16 功能
	if !CheckCRC16(withCRC) {
		t.Error("CheckCRC16() returned false for valid CRC")
	}

	// 修改数据，验证 CheckCRC16 能够检测到
	withCRC[0] = 0x02
	if CheckCRC16(withCRC) {
		t.Error("CheckCRC16() returned true for invalid CRC")
	}
}

// 模拟传输接口
type mockTransport struct {
	t          *testing.T
	expectedTx []byte
	mockRx     []byte
}

func (m *mockTransport) Write(p []byte) (n int, err error) {
	// 验证发送的数据是否符合预期
	if len(m.expectedTx) > 0 && !bytes.Equal(p, m.expectedTx) {
		m.t.Errorf("Expected to send:\n%s\nActually sent:\n%s",
			hex.Dump(m.expectedTx), hex.Dump(p))
	}
	return len(p), nil
}

func (m *mockTransport) Read(p []byte) (n int, err error) {
	// 模拟接收数据
	if len(p) < len(m.mockRx) {
		return 0, bytes.ErrTooLarge
	}
	n = copy(p, m.mockRx)
	return n, nil
}

// 测试读取保持寄存器请求和响应
func TestReadHoldingRegisters(t *testing.T) {
	// 模拟请求和响应
	// 注意：我们使用我们自己的CRC计算函数来生成正确的CRC值
	requestData := []byte{0x01, 0x03, 0x00, 0x6B, 0x00, 0x03}
	requestCRC := CRC16(requestData)
	expectedRequest := append(requestData, byte(requestCRC), byte(requestCRC>>8))

	// 响应数据也使用我们的CRC函数计算
	responseData := []byte{0x01, 0x03, 0x06, 0x02, 0x2B, 0x00, 0x00, 0x00, 0x64}
	responseCRC := CRC16(responseData)
	mockResponse := append(responseData, byte(responseCRC), byte(responseCRC>>8))

	transport := &mockTransport{
		t:          t,
		expectedTx: expectedRequest,
		mockRx:     mockResponse,
	}

	client := NewClient(transport, 0x01)

	// 读取寄存器
	registers, err := client.ReadHoldingRegisters(0x6B, 0x03)
	if err != nil {
		t.Fatalf("ReadHoldingRegisters() error = %v", err)
	}

	// 验证结果
	expectedRegisters := []uint16{0x022B, 0x0000, 0x0064}
	if len(registers) != len(expectedRegisters) {
		t.Errorf("ReadHoldingRegisters() returned %d registers, want %d", len(registers), len(expectedRegisters))
	}

	for i, reg := range registers {
		if reg != expectedRegisters[i] {
			t.Errorf("ReadHoldingRegisters()[%d] = 0x%04X, want 0x%04X", i, reg, expectedRegisters[i])
		}
	}
}

// 测试写单个线圈请求和响应
func TestWriteSingleCoil(t *testing.T) {
	// 模拟请求和响应 (写 ON)
	requestData := []byte{0x01, 0x05, 0x00, 0xAC, 0xFF, 0x00}
	requestCRC := CRC16(requestData)
	expectedRequest := append(requestData, byte(requestCRC), byte(requestCRC>>8))

	// 响应与请求相同
	mockResponse := expectedRequest

	transport := &mockTransport{
		t:          t,
		expectedTx: expectedRequest,
		mockRx:     mockResponse,
	}

	client := NewClient(transport, 0x01)

	// 写线圈
	err := client.WriteSingleCoil(0xAC, true)
	if err != nil {
		t.Fatalf("WriteSingleCoil() error = %v", err)
	}

	// 模拟请求和响应 (写 OFF)
	requestData = []byte{0x01, 0x05, 0x00, 0xAC, 0x00, 0x00}
	requestCRC = CRC16(requestData)
	expectedRequest = append(requestData, byte(requestCRC), byte(requestCRC>>8))

	// 响应与请求相同
	mockResponse = expectedRequest

	transport.expectedTx = expectedRequest
	transport.mockRx = mockResponse

	// 写线圈
	err = client.WriteSingleCoil(0xAC, false)
	if err != nil {
		t.Fatalf("WriteSingleCoil() error = %v", err)
	}
}

// 测试写多个寄存器请求和响应
func TestWriteMultipleRegisters(t *testing.T) {
	// 模拟请求和响应
	requestData := []byte{0x01, 0x10, 0x00, 0x01, 0x00, 0x02, 0x04, 0x00, 0x0A, 0x01, 0x02}
	requestCRC := CRC16(requestData)
	expectedRequest := append(requestData, byte(requestCRC), byte(requestCRC>>8))

	responseData := []byte{0x01, 0x10, 0x00, 0x01, 0x00, 0x02}
	responseCRC := CRC16(responseData)
	mockResponse := append(responseData, byte(responseCRC), byte(responseCRC>>8))

	transport := &mockTransport{
		t:          t,
		expectedTx: expectedRequest,
		mockRx:     mockResponse,
	}

	client := NewClient(transport, 0x01)

	// 写寄存器
	err := client.WriteMultipleRegisters(0x01, []uint16{0x000A, 0x0102})
	if err != nil {
		t.Fatalf("WriteMultipleRegisters() error = %v", err)
	}
}

// 测试异常响应处理
func TestExceptionResponse(t *testing.T) {
	// 模拟请求和异常响应 (0x83 = 0x03 | 0x80, 错误码 0x02 = 非法数据地址)
	requestData := []byte{0x01, 0x03, 0x00, 0x6B, 0x00, 0x03}
	requestCRC := CRC16(requestData)
	expectedRequest := append(requestData, byte(requestCRC), byte(requestCRC>>8))

	// 异常响应
	responseData := []byte{0x01, 0x83, 0x02}
	responseCRC := CRC16(responseData)
	mockResponse := append(responseData, byte(responseCRC), byte(responseCRC>>8))

	transport := &mockTransport{
		t:          t,
		expectedTx: expectedRequest,
		mockRx:     mockResponse,
	}

	client := NewClient(transport, 0x01)

	// 读取寄存器，预期会返回错误
	_, err := client.ReadHoldingRegisters(0x6B, 0x03)
	if err == nil {
		t.Fatal("ReadHoldingRegisters() error = nil, expected an error")
	}

	// 验证错误类型
	modbusError, ok := err.(*ModbusError)
	if !ok {
		t.Fatalf("err is not a ModbusError, got %T", err)
	}

	if modbusError.FunctionCode != FuncReadHoldingRegisters {
		t.Errorf("ModbusError.FunctionCode = 0x%02X, want 0x%02X", modbusError.FunctionCode, FuncReadHoldingRegisters)
	}

	if modbusError.ExceptionCode != ExcIllegalDataAddress {
		t.Errorf("ModbusError.ExceptionCode = 0x%02X, want 0x%02X", modbusError.ExceptionCode, ExcIllegalDataAddress)
	}
}
