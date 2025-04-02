package main

// 运行此示例前，请安装以下依赖：
// go get github.com/tarm/serial

import (
	"fmt"
	"log"
	"time"

	// 用于串口通信
	"github.com/wxlbd/gokit/protocols/modbus"
)

// MockTransport 实现了一个模拟的传输接口，用于演示
type MockTransport struct {
	ExpectedRequest []byte
	MockResponse    []byte
}

func (m *MockTransport) Write(p []byte) (n int, err error) {
	// 打印发送的请求数据
	fmt.Printf("发送请求: %X\n", p)
	return len(p), nil
}

func (m *MockTransport) Read(p []byte) (n int, err error) {
	// 返回模拟的响应数据
	if len(p) < len(m.MockResponse) {
		return 0, fmt.Errorf("缓冲区太小")
	}
	n = copy(p, m.MockResponse)
	return n, nil
}

// 此示例展示了如何使用 Modbus RTU 客户端与设备通信
// 实际使用时需要根据设备的具体参数进行配置

func main() {
	fmt.Println("Modbus RTU 客户端示例")
	fmt.Println("=================================")

	// 创建模拟传输接口
	// 这里我们模拟读取保持寄存器的请求和响应
	mock := &MockTransport{
		// 模拟的响应: 从站ID=1, 功能码=3, 字节数=4, 2个寄存器值(0x1234, 0x5678), CRC
		MockResponse: []byte{0x01, 0x03, 0x04, 0x12, 0x34, 0x56, 0x78, 0x73, 0x8F},
	}

	// 创建 Modbus 客户端
	client := modbus.NewClient(mock, 1)               // 从站 ID = 1
	client.SetTimeout(time.Second * 2)                // 设置超时时间
	client.SetInterFrameDelay(time.Millisecond * 100) // 设置帧间延时

	fmt.Println("\n1. 读取保持寄存器示例")
	registers, err := client.ReadHoldingRegisters(0, 2) // 从地址 0 开始读取 2 个寄存器
	if err != nil {
		log.Printf("读取保持寄存器失败: %v", err)
	} else {
		fmt.Println("读取成功:")
		for i, reg := range registers {
			fmt.Printf("  寄存器 %d: 0x%04X (%d)\n", i, reg, reg)
		}
	}

	// 模拟读取线圈状态
	fmt.Println("\n2. 读取线圈状态示例")
	mock.MockResponse = []byte{0x01, 0x01, 0x01, 0x05, 0x8D, 0xAC} // 响应：5个线圈 (0b00000101)

	coils, err := client.ReadCoils(0, 8) // 从地址 0 开始读取 8 个线圈
	if err != nil {
		log.Printf("读取线圈状态失败: %v", err)
	} else {
		fmt.Println("读取成功:")
		for i, status := range coils {
			fmt.Printf("  线圈 %d: %v\n", i, status)
		}
	}

	// 模拟写单个寄存器
	fmt.Println("\n3. 写单个寄存器示例")
	mock.MockResponse = []byte{0x01, 0x06, 0x00, 0x01, 0x00, 0x03, 0x9A, 0x9B} // 响应：写入成功

	err = client.WriteSingleRegister(1, 3) // 向地址 1 写入值 3
	if err != nil {
		log.Printf("写单个寄存器失败: %v", err)
	} else {
		fmt.Println("写入成功")
	}

	// 模拟写单个线圈
	fmt.Println("\n4. 写单个线圈示例")
	mock.MockResponse = []byte{0x01, 0x05, 0x00, 0x0A, 0xFF, 0x00, 0xE8, 0x16} // 响应：写入成功

	err = client.WriteSingleCoil(10, true) // 向地址 10 写入值 ON
	if err != nil {
		log.Printf("写单个线圈失败: %v", err)
	} else {
		fmt.Println("写入成功")
	}

	// 模拟写多个寄存器
	fmt.Println("\n5. 写多个寄存器示例")
	mock.MockResponse = []byte{0x01, 0x10, 0x00, 0x01, 0x00, 0x02, 0x13, 0xCC} // 响应：写入成功

	err = client.WriteMultipleRegisters(1, []uint16{0x000A, 0x0102}) // 从地址 1 开始写入多个值
	if err != nil {
		log.Printf("写多个寄存器失败: %v", err)
	} else {
		fmt.Println("写入成功")
	}

	// 模拟写多个线圈
	fmt.Println("\n6. 写多个线圈示例")
	mock.MockResponse = []byte{0x01, 0x0F, 0x00, 0x14, 0x00, 0x05, 0x17, 0xB9} // 响应：写入成功

	err = client.WriteMultipleCoils(20, []bool{true, false, true, false, true}) // 从地址 20 开始写入多个值
	if err != nil {
		log.Printf("写多个线圈失败: %v", err)
	} else {
		fmt.Println("写入成功")
	}

	// 模拟异常响应
	fmt.Println("\n7. 处理异常响应示例")
	mock.MockResponse = []byte{0x01, 0x83, 0x02, 0xC0, 0xF1} // 异常响应：非法数据地址

	_, err = client.ReadHoldingRegisters(0, 10)
	if err != nil {
		log.Printf("预期的错误: %v", err)
	} else {
		fmt.Println("应该捕获到异常但没有")
	}
}
