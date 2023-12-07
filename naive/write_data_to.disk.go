package main

import (
	"fmt"
	"os"
	"time"
	"bytes"
)

func assert(b bool) {
	if !b {
		panic("assert")
	}
}
const BUFFER_SIZE = 4096

func readNBytesdata(fileName string, numberOfBytesRead int) []byte{
	fileDescriptor, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fileDescriptor.Close()

	data := make([]byte, 0, numberOfBytesRead)
	var buffer = make([]byte, BUFFER_SIZE)
	for len(data) < numberOfBytesRead {
		read, err := fileDescriptor.Read(buffer)
		if err != nil {
			panic(err)
		}
		data = append(data, buffer[:read]...)
	}
	assert(len(data) == numberOfBytesRead)
	return data

}

func benchmark(name string, data []byte, fn func(*os.File)) {
	fmt.Printf("%s", name)
	fd, err := os.OpenFile("out.bin", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	beginTime := time.Now()
	fn(fd)
	currTime := time.Now().Sub(beginTime).Seconds()
	fmt.Printf(",%f, %f\n", currTime, float64(len(data))/currTime)

	if err := fd.Close(); err != nil {
		panic(err)
	}

	assert(bytes.Equal(readNBytesdata("out.bin", len(data)), data))
}
func main() {
	size := 104857600
	data := readNBytesdata("/dev/random", size)
	const RUNS = 10
	for i := 0 ; i < RUNS; i++ {
		benchmark("blocking", data, func(f *os.File){
			for i := 0 ; i < len(data); i += BUFFER_SIZE {
				size := min(BUFFER_SIZE, len(data)-i)
				writeBytes, err := f.Write(data[i:i+size])
				if err != nil {
					panic(err)
				}
				assert(writeBytes == BUFFER_SIZE)
			}

		})
	}
	
}
