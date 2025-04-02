package modbus

import (
	"io"
	"time"
)

// Client 是 Modbus RTU 客户端
type Client struct {
	transport       io.ReadWriter // 通讯接口
	slaveID         byte          // 从站 ID
	timeout         time.Duration // 超时时间
	interFrameDelay time.Duration // 帧间延时
}

// NewClient 创建一个新的 Modbus RTU 客户端
func NewClient(transport io.ReadWriter, slaveID byte) *Client {
	return &Client{
		transport:       transport,
		slaveID:         slaveID,
		timeout:         1 * time.Second,
		interFrameDelay: 100 * time.Millisecond,
	}
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	return c
}

// SetInterFrameDelay 设置帧间延时
func (c *Client) SetInterFrameDelay(delay time.Duration) *Client {
	c.interFrameDelay = delay
	return c
}

// SetSlaveID 设置从站 ID
func (c *Client) SetSlaveID(slaveID byte) *Client {
	c.slaveID = slaveID
	return c
}

// 发送请求并读取响应
func (c *Client) sendAndReceive(request []byte, expectedFunctionCode byte) ([]byte, error) {
	// 清空接收缓冲区
	// 注意：这个步骤依赖于具体实现，可能需要根据实际情况进行调整或移除
	/*
		if flusher, ok := c.transport.(interface{ Flush() error }); ok {
			if err := flusher.Flush(); err != nil {
				return nil, err
			}
		}
	*/

	// 发送请求
	if _, err := c.transport.Write(request); err != nil {
		return nil, err
	}

	// 等待帧间延时
	time.Sleep(c.interFrameDelay)

	// 读取响应
	// 注意：实际应用中，你可能需要处理更复杂的读取逻辑，例如处理超时
	buffer := make([]byte, 256) // 足够大的缓冲区
	n, err := c.transport.Read(buffer)
	if err != nil {
		return nil, err
	}

	response := buffer[:n]

	// 验证响应
	if err := ValidateResponse(response, c.slaveID, expectedFunctionCode); err != nil {
		return nil, err
	}

	return response, nil
}

// ================= 位操作功能 =================

// ReadCoils 读取线圈状态
func (c *Client) ReadCoils(startAddress uint16, quantity uint16) ([]bool, error) {
	request := NewReadCoilsRequest(c.slaveID, startAddress, quantity)
	response, err := c.sendAndReceive(request, FuncReadCoils)
	if err != nil {
		return nil, err
	}

	bits, err := ParseReadBitsResponse(response, c.slaveID, FuncReadCoils)
	if err != nil {
		return nil, err
	}

	// 只返回请求的位数量
	if int(quantity) < len(bits) {
		bits = bits[:quantity]
	}

	return bits, nil
}

// ReadDiscreteInputs 读取离散输入状态
func (c *Client) ReadDiscreteInputs(startAddress uint16, quantity uint16) ([]bool, error) {
	request := NewReadDiscreteInputsRequest(c.slaveID, startAddress, quantity)
	response, err := c.sendAndReceive(request, FuncReadDiscreteInputs)
	if err != nil {
		return nil, err
	}

	bits, err := ParseReadBitsResponse(response, c.slaveID, FuncReadDiscreteInputs)
	if err != nil {
		return nil, err
	}

	// 只返回请求的位数量
	if int(quantity) < len(bits) {
		bits = bits[:quantity]
	}

	return bits, nil
}

// WriteSingleCoil 写单个线圈
func (c *Client) WriteSingleCoil(address uint16, value bool) error {
	request := NewWriteSingleCoilRequest(c.slaveID, address, value)
	response, err := c.sendAndReceive(request, FuncWriteSingleCoil)
	if err != nil {
		return err
	}

	return ParseWriteSingleCoilResponse(response, c.slaveID, address, value)
}

// WriteMultipleCoils 写多个线圈
func (c *Client) WriteMultipleCoils(startAddress uint16, values []bool) error {
	request := NewWriteMultipleCoilsRequest(c.slaveID, startAddress, values)
	if request == nil {
		return ErrInvalidLength
	}

	response, err := c.sendAndReceive(request, FuncWriteMultipleCoils)
	if err != nil {
		return err
	}

	return ParseWriteMultipleCoilsResponse(response, c.slaveID, startAddress, uint16(len(values)))
}

// ================= 字操作功能 =================

// ReadHoldingRegisters 读取保持寄存器
func (c *Client) ReadHoldingRegisters(startAddress uint16, quantity uint16) ([]uint16, error) {
	request := NewReadHoldingRegistersRequest(c.slaveID, startAddress, quantity)
	response, err := c.sendAndReceive(request, FuncReadHoldingRegisters)
	if err != nil {
		return nil, err
	}

	return ParseReadRegistersResponse(response, c.slaveID, FuncReadHoldingRegisters)
}

// ReadInputRegisters 读取输入寄存器
func (c *Client) ReadInputRegisters(startAddress uint16, quantity uint16) ([]uint16, error) {
	request := NewReadInputRegistersRequest(c.slaveID, startAddress, quantity)
	response, err := c.sendAndReceive(request, FuncReadInputRegisters)
	if err != nil {
		return nil, err
	}

	return ParseReadRegistersResponse(response, c.slaveID, FuncReadInputRegisters)
}

// WriteSingleRegister 写单个寄存器
func (c *Client) WriteSingleRegister(address uint16, value uint16) error {
	request := NewWriteSingleRegisterRequest(c.slaveID, address, value)
	response, err := c.sendAndReceive(request, FuncWriteSingleRegister)
	if err != nil {
		return err
	}

	return ParseWriteSingleRegisterResponse(response, c.slaveID, address, value)
}

// WriteMultipleRegisters 写多个寄存器
func (c *Client) WriteMultipleRegisters(startAddress uint16, values []uint16) error {
	request := NewWriteMultipleRegistersRequest(c.slaveID, startAddress, values)
	if request == nil {
		return ErrInvalidLength
	}

	response, err := c.sendAndReceive(request, FuncWriteMultipleRegisters)
	if err != nil {
		return err
	}

	return ParseWriteMultipleRegistersResponse(response, c.slaveID, startAddress, uint16(len(values)))
}
