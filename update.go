package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

var (
	errStreamClosed    = errors.New("Stream closed unexpectedly")
	errInvalidResponse = errors.New("Invalid response from WiFish")
)

const (
	uploadCompletedSignature = "0123456789012345678901234EQUADRO"
)

func updatePlanet(destination string, port int, firmware []byte) error {
	addr, err := net.ResolveIPAddr("ip", destination)
	if err != nil {
		return err
	}
	fmt.Println("Connecting", addr.String())
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: addr.IP, Port: port})
	reader := bufio.NewReader(conn)

	fmt.Println("Handshake")
	helloBuff := make([]byte, 7, 7)
	_, err = reader.Read(helloBuff)
	if err != nil {
		if err == io.EOF {
			return errStreamClosed
		}
		return err
	}
	if bytes.Compare([]byte("*HELLO*"), helloBuff) != 0 {
		fmt.Println("Error", helloBuff)
		return errInvalidResponse
	}

	fmt.Println("Request download mode")
	conn.Write([]byte("\x02PROGRAMPLANET01\x03"))
	data, _, err := reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return errStreamClosed
		}
		return err
	}
	if bytes.Compare([]byte("PLANETGOINDWL01"), data) != 0 {
		return errInvalidResponse
	}

	fmt.Println("Headers")
	conn.Write([]byte("\x02SENDHEADER012000000\x03"))
	data, _, err = reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return errStreamClosed
		}
		return err
	}
	if bytes.Compare([]byte("PLANETPACKETOK01"), data) != 0 {
		// what about closing the socket?
		return errInvalidResponse
	}

	pktsContent := preparePackets(firmware)

	fmt.Println("Firmware upload")
	for i, d := range pktsContent {
		conn.Write([]byte(fmt.Sprintf("\x02PLANETPACKET01%04d%s\x03", i+1, d)))
		data, _, err = reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return errStreamClosed
			}
			return err
		}
		if bytes.Compare([]byte("PLANETPACKETOK01"), data) != 0 {
			// what about closing the socket?
			return errInvalidResponse
		}
		fmt.Println("Uploading", 100*i/len(pktsContent))
	}

	fmt.Println("Sending completed")
	conn.Write([]byte(fmt.Sprintf("\x02PLANETPACKET01%04d%s\x03", len(pktsContent)+1, strings.ToUpper(hex.EncodeToString([]byte(uploadCompletedSignature))))))
	data, _, err = reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return errStreamClosed
		}
		return err
	}
	if bytes.Compare([]byte("PLANETPACKETOK01"), data) != 0 {
		// what about closing the socket?
		return errInvalidResponse
	}
	conn.Close()
	fmt.Println("Completed")
	return nil
}

func preparePackets(bin []byte) []string {
	binLen := len(bin)
	numPackets := binLen / 32
	if binLen%32 != 0 {
		numPackets++
	}
	fmt.Println(binLen)
	fmt.Println(numPackets)
	pkts := make([]string, numPackets)
	writtenBytes := 0
	for i := range pkts {
		if writtenBytes < binLen-31 {
			pkts[i] = strings.ToUpper(hex.EncodeToString(bin[i*32 : (i+1)*32]))
		} else {
			prefix := strings.ToUpper(hex.EncodeToString(bin[i*32:]))
			pkts[i] = fmt.Sprint(prefix, strings.Repeat("U", 32-len(prefix)))
		}
		writtenBytes += 32
	}
	return pkts
}
