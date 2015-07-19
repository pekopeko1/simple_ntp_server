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
    //sb := strings.Repeat("\x00", 48)
    var sb [48]byte
    for i := range sb {
	sb[i] = "\x00"
    }
    for {
        handleClient(conn, sb)
    }
}
func handleClient(conn *net.UDPConn, sb *string) {
    var buf [256]byte
    size, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
        return
    }
    if size < 48 {
	return
    }
    sec, microsec = gettimeofday()
    vn := buf[0:][0:3] & 0x38000000
    sb[0:3] = 0x040106F0 | vn // flag
    sb[4] = 0 // delay
    sb[8:11] = 0x00000010 // dispersion
    sb[12:15] = "LOCL" // Ref ID
    sb[16:19] = sec + 0x83AA7E80
    sb[20:23] = 0
    sb[24:31] = buf[0:][40:47]
    sb[32:35] = sec + 0x83AA7E80
    sb[36:39] = (microsec / 500000) * 0x80000000
    sec, microsec = gettimeofday()
    sb[40:43] = sec + 0x83AA7E80
    sb[44:47] = (microsec / 500000) * 0x80000000

    conn.WriteToUDP([]byte(sb), addr)
}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
        os.Exit(1)
    }
}

func gettimeofday() {
    now = time.Nanoseconds()
    sec := now / 1e9
    microsec := now / 1e6
    return sec, microsec
}
