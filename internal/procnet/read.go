package procnet

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var validNInterfaces = map[string]bool{
	"enp4s0": true, // Ethernet
	"wlp3s0": true, // WiFi
}

func validNIsFunc(nif string) bool {
	return validNInterfaces[nif]
}

func ReadTotalDataUsage() (Data, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return Data{}, fmt.Errorf("error reading /proc/net/dev: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Skip the two header lines
	scanner.Scan()
	scanner.Scan()

	var totalRx, totalTx uint64

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// /proc/net/dev lines usually have 17 fields
		if len(fields) < 10 {
			continue
		}

		// Clean the interface name (e.g., "eth0:")
		name := strings.TrimSuffix(fields[0], ":")

		if !validNIsFunc(name) {
			continue
		}

		// fields[1] is Received Bytes, fields[9] is Transmitted Bytes
		rx, _ := strconv.ParseUint(fields[1], 10, 64)
		tx, _ := strconv.ParseUint(fields[9], 10, 64)

		totalRx += rx
		totalTx += tx
	}

	if err := scanner.Err(); err != nil {
		return Data{}, fmt.Errorf("error reading /proc/net/dev: %w", err)
	}

	return Data{totalRx, totalTx}, nil
}
