package procnet

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadStats(filterFunc FilterFunc) ([]Stats, error) {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var stats []Stats
	scanner := bufio.NewScanner(f)

	// Skip the two header lines
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		// Network Interface name is "eth0:", strip the colon
		name := strings.TrimSuffix(fields[0], ":")

		if filterFunc != nil {
			if validNI := filterFunc(name); !validNI {
				continue
			}
		}

		var rx, tx uint64
		fmt.Scanf(fields[1], "%d", &rx) // Recieved bytes
		fmt.Scanf(fields[9], "%d", &tx) // Transmitted bytes

		stats = append(stats, Stats{name, rx, tx})
	}
	return stats, nil
}
