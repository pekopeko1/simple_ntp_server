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
	//println(size, addr)
	if err != nil {
		return
	}
	if size < 48 {
		return
	}
	sec, microsec := gettimeofday()
	vn := byteArrayToUint64(buf[0:][0:3]) & 0x38000000
	copy(sb[0:3], uint64ToByteArray(0x040106F0|vn)) // flag
	copy(sb[4:7], []byte{0, 0, 0, 0})               //delay
	copy(sb[8:11], uint64ToByteArray(0x00000010))   // dispersion
	copy(sb[12:15], []byte("LOCL"))                 // Ref ID
	copy(sb[16:19], uint64ToByteArray(sec+0x83AA7E80))
	copy(sb[20:23], []byte{0, 0, 0, 0})
	copy(sb[24:31], buf[0:][40:47])
	copy(sb[32:35], uint64ToByteArray(sec+0x83AA7E80))
	copy(sb[36:39], uint64ToByteArray((microsec/500000)*0x80000000))
	sec, microsec = gettimeofday()
	copy(sb[40:43], uint64ToByteArray(sec+0x83AA7E80))
	copy(sb[44:47], uint64ToByteArray((microsec/500000)*0x80000000))

	conn.WriteToUDP([]byte(sb), addr)
}

func uint64ToByteArray(num uint64) (ret []byte) {
	ret = make([]byte, 4)
	ret[0] = byte(0xFF000000 & num >> 24)
	ret[1] = byte(0x00FF0000 & num >> 16)
	ret[2] = byte(0x0000FF00 & num >> 8)
	ret[3] = byte(0x000000FF & num)
	//println(ret[0], ret[1], ret[2], ret[3])
	return ret
}

func byteArrayToUint64(ba []byte) (ret uint64) {
	switch len(ba) {
	case 4:
		//right way
	case 3:
		ba = append(ba, 0)
	default:
		panic(len(ba))
	}
	for i := 0; i < 4; i++ {
		ret += uint64(ba[i]) << uint64(24-(8*i))
	}
	return ret
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func gettimeofday() (sec uint64, microsec uint64) {
	now := time.Now().UnixNano()
	sec = uint64(now / 1e9)
	microsec = uint64(now / 1e6)
	return sec, microsec
}
