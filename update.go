package main

/* Planet Updater - easily update your Elos Planet lamps.
Copyright (C) 2015 Francesco Zanini <francesco@zanini.me>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.*/

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

type updateProgress struct {
	state   int
	message string
	pkts    int
}

func updatePlanet(destination string, port int, firmware []byte, c chan updateProgress) error {
	c <- updateProgress{0, "", 0}
	addr, err := net.ResolveIPAddr("ip", destination)
	if err != nil {
		return err
	}
	c <- updateProgress{1, fmt.Sprintf("%s:%d", addr.String(), port), 0}
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: addr.IP, Port: port})
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)

	c <- updateProgress{2, "", 0}
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

	c <- updateProgress{3, "", 0}
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

	c <- updateProgress{4, "", 0}
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

	c <- updateProgress{5, "", len(pktsContent)}
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
		c <- updateProgress{6, "", 0}
	}

	c <- updateProgress{6, "", 1}
	conn.Write([]byte(fmt.Sprintf("\x02PLANETPACKET01%04d%s\x03", len(pktsContent)+1, strings.ToUpper(hex.EncodeToString([]byte(uploadCompletedSignature))))))
	c <- updateProgress{7, "", 0}
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
	c <- updateProgress{8, "", 0}
	close(c)
	return nil
}

func preparePackets(bin []byte) []string {
	binLen := len(bin)
	numPackets := binLen / 32
	if binLen%32 != 0 {
		numPackets++
	}
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
