package main

import (
	"encoding/binary"
	"encoding/hex"
	//_ "encoding/hex"
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

	//buf := make([]byte, 1024)
	//	02 CA A7 93 2F FE D7                               .К§“/юЧ

	buf := []byte{0x02, 0xCA, 0xA7, 0x93, 0x63, 0xFF, 0x22}

	int_value := 46835603
	hex_value := fmt.Sprintf("%X", int_value)
	fmt.Printf("Hex value of %d is = %s\n", int_value, hex_value)
	bs := make([]byte, 8)
	binary.AppendVarint(bs, 46835603)
	fmt.Println("Tx: ", bs)

	//crc16(buf)
	//for {
	_, err = stream.Write(crc16(buf))
	if err != nil {
		log.Fatal(err)
	}
	//	s := string(buf[:n])
	fmt.Println("Tx: ", hex.EncodeToString(crc16(buf)))

	//}
	//buf1 := make([]byte, 1024)
	var st string
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
				fmt.Println("Rx: ", hex.EncodeToString(buf[:n]))
				st = st + hex.EncodeToString(buf[:n])
			}

		}
	}
	fmt.Println("Rx: ", st)
}
