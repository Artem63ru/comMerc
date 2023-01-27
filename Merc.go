package main

import (
	"encoding/hex"
	"io"
)
import "fmt"
import "log"
import "github.com/tarm/serial"

func main() {
	config := &serial.Config{
		Name:        "COM2",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	int_value := 46835603
	hex_value := fmt.Sprintf("%08x", int_value)
	fmt.Printf("Decimal: %d,\n Hexa: %s", int_value, hex_value)
	hex_value1, err := hex.DecodeString(hex_value)
	fmt.Println(hex_value1)
	if err != nil {
		panic(err)
	}
	hex_value1 = append(hex_value1, 0x61)
	//crc16(hex_value1)

	_, err = stream.Write(crc16(hex_value1))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tx: ", hex.EncodeToString(crc16(hex_value1)))

	var st []byte
	for i := 0; i < 100; i++ {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {

			//	buf = buf[:n]
			if n > 0 {
				//	fmt.Println("Rx: ", hex.EncodeToString(buf[:n]))
				//	st = st + hex.EncodeToString(buf[:n])
				st = append(st, buf[:n]...)
			}

		}
	}
	fmt.Println("Rx: ", hex.EncodeToString(st))
}
