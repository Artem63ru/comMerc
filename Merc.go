package main

import (
	"encoding/hex"
	"io"
	"strconv"
)
import "fmt"
import "log"
import "github.com/tarm/serial"

func convert(x int, y int, divider int, arr []byte) float32 {
	largeArray := make([]byte, 2)
	copy(largeArray, arr[x:y])
	i, err := strconv.Atoi(hex.EncodeToString(largeArray))
	if err != nil {
		return 0
		panic(err)
	} else {
		return float32(i) / float32(divider)
	}
}

func send_to(s_n int, stream serial.Port, comm byte) []byte {
	hex_value := fmt.Sprintf("%08x", s_n)
	fmt.Printf("Decimal: %d,\n Hexa: %s", s_n, hex_value)
	hex_value1, err := hex.DecodeString(hex_value)
	fmt.Println(hex_value1)
	if err != nil {
		panic(err)
	}
	hex_value1 = append(hex_value1, comm)
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
			if n > 0 {
				st = append(st, buf[:n]...)
			}

		}
	}
	return st
}

func main() {
	config := &serial.Config{
		Name:        "COM2",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}
	var v float32
	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	int_value := 46835603
	//hex_value := fmt.Sprintf("%08x", int_value)
	//fmt.Printf("Decimal: %d,\n Hexa: %s", int_value, hex_value)
	//hex_value1, err := hex.DecodeString(hex_value)
	//fmt.Println(hex_value1)
	//if err != nil {
	//	panic(err)
	//}
	//hex_value1 = append(hex_value1, 0x63)
	//_, err = stream.Write(crc16(hex_value1))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Tx: ", hex.EncodeToString(crc16(hex_value1)))
	//var v float32
	//var st []byte
	//for i := 0; i < 100; i++ {
	//	buf := make([]byte, 1024)
	//	n, err := stream.Read(buf)
	//	if err != nil {
	//		if err != io.EOF {
	//			fmt.Println("Error reading from serial port: ", err)
	//		}
	//	} else {
	//		if n > 0 {
	//			st = append(st, buf[:n]...)
	//		}
	//
	//	}
	//}
	st := send_to(int_value, *stream, 0x63)
	fmt.Println("Rx: ", hex.EncodeToString(st))
	st1 := crc16(st[:len(st)-2])
	if hex.EncodeToString(st[len(st)-2:]) == hex.EncodeToString(st1[len(st1)-2:]) {
		fmt.Println("crc ОК")
		if st[4] == 99 {
			v = convert(5, 8, 10, st)
			i := convert(9, 10, 100, st)
			p := convert(11, 13, 1000, st)
			fmt.Println("Напряжение", v)
			fmt.Println("Ток", i)
			fmt.Println("Мощность", p)
		} else {
			fmt.Println("Не тот ответ")
		}

	}
}
