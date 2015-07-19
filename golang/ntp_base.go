package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":123"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	sb := make([]byte, 48)
	for i := range sb {
		sb[i] = 0 //"\x00"
	}
	for {
		handleClient(conn, sb)
	}
}
func handleClient(conn *net.UDPConn, sb []byte) {
	var buf [256]byte
	size, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	if size < 48 {
		return
	}
	sec, microsec := gettimeofday()
	vn := byteArrayToUint32(buf[0:][0:3]) & 0x38000000
	copy(sb[0:3], uint32ToByteArray(0x040106F0|vn)) // flag
	sb[4] = 0                                       // delay
	copy(sb[8:11], uint32ToByteArray(0x00000010))   // dispersion
	copy(sb[12:15], []byte("LOCL"))                 // Ref ID
	copy(sb[16:19], uint32ToByteArray(sec+0x83AA7E80))
	copy(sb[20:23], []byte{0, 0, 0, 0})
	copy(sb[24:31], buf[0:][40:47])
	copy(sb[32:35], uint32ToByteArray(sec+0x83AA7E80))
	copy(sb[36:39], uint32ToByteArray((microsec/500000)*0x80000000))
	sec, microsec = gettimeofday()
	copy(sb[40:43], uint32ToByteArray(sec+0x83AA7E80))
	copy(sb[44:47], uint32ToByteArray((microsec/500000)*0x80000000))

	conn.WriteToUDP([]byte(sb), addr)
}

func uint32ToByteArray(num uint32) (ret []byte) {
	ret[0] = byte((0xff000000 | num) >> 24)
	ret[1] = byte((0x00ff0000 | num) >> 16)
	ret[2] = byte((0x0000ff00 | num) >> 8)
	ret[3] = byte(0x000000ff | num)
	return ret
}

func byteArrayToUint32(ba []byte) (ret uint32) {
	if len(ba) != 4 {
		panic(0)
	}
	for i := 0; i < 4; i++ {
		ret += uint32(ba[i]) << uint32(24-(8*i))
	}
	return ret
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func gettimeofday() (sec uint32, microsec uint32) {
	now := time.Now().UnixNano()
	sec = uint32(now / 1e9)
	microsec = uint32(now / 1e6)
	return sec, microsec
}
