package gopoke

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum>>16 > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}

func stripIPHeaderIFPresent(data []byte) []byte {
	if len(data) > 0 && (data[0]>>4) == 4 {
		ihl := int(data[0]&0x0f) * 4
		if len(data) >= ihl {
			return data[ihl:]
		}
	}
	return data
}

type PingResult struct {
	Bytes int
	Type  uint8
}

func ExtractHost(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "://") {
		return rawURL, nil
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func Ping(host string) (PingResult, error) {
	var result PingResult
	conn, err := net.Dial("ip4:icmp", host)
	if err != nil {
		return result, fmt.Errorf("接続エラー（root権限が必要です）：%w", err)
	}
	defer conn.Close()

	id := uint16(os.Getpid() & 0xffff)
	seq := uint16(1)

	msg := make([]byte, 8)
	msg[0] = 8 // Type: Echo Request
	msg[1] = 0
	binary.BigEndian.PutUint16(msg[4:6], id)
	binary.BigEndian.PutUint16(msg[6:8], seq)
	binary.BigEndian.PutUint16(msg[2:4], checksum(msg))

	if _, err := conn.Write(msg); err != nil {
		return result, fmt.Errorf("送信エラー：%w", err)
	}

	reply := make([]byte, 1500)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(reply)
	if err != nil {
		return result, fmt.Errorf("受信エラー：%w", err)
	}

	icmpReply := stripIPHeaderIFPresent(reply[:n])
	if len(icmpReply) < 1 {
		return result, fmt.Errorf("ICMP応答が不正です")
	}
	if icmpReply[0] != 0 {
		return result, fmt.Errorf("ICMP応答が不正です: Type=%d", icmpReply[0])
	}

	result = PingResult{
		Bytes: len(icmpReply),
		Type:  icmpReply[0],
	}
	return result, nil
}
