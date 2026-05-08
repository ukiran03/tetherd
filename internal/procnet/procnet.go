package procnet

import "fmt"

type Data struct {
	RxBytes uint64 // Recieved bytes
	TxBytes uint64 // Transmitted bytes
}

func (d Data) GreaterThan(other Data) bool {
	return (d.RxBytes + d.TxBytes) > (other.RxBytes + other.TxBytes)
}

func (d Data) GreaterOrEq(other Data) bool {
	return d.GreaterThan(other) ||
		(d.RxBytes+d.TxBytes) == (other.RxBytes+other.TxBytes)
}

// Delta calculates (d - previous).
// It returns 0 if previous is larger than d (to prevent uint64 wrap-around).
func (d Data) Delta(previous Data) Data {
	var res Data
	if d.RxBytes > previous.RxBytes {
		res.RxBytes = d.RxBytes - previous.RxBytes
	}
	if d.TxBytes > previous.TxBytes {
		res.TxBytes = d.TxBytes - previous.TxBytes
	}
	return res
}

// Add sums two Data structs together.
func (d Data) Add(other Data) Data {
	return Data{
		RxBytes: d.RxBytes + other.RxBytes,
		TxBytes: d.TxBytes + other.TxBytes,
	}
}

// String provides a human-readable version for logging.
func (d Data) String() string {
	return fmt.Sprintf("RX: %d, TX: %d", d.RxBytes, d.TxBytes)
}

func (d Data) PrintHuman() string {
	return fmt.Sprintf("Total Usage: %v", humanSize(d.RxBytes+d.TxBytes))
}

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func humanSize(size uint64) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fG", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fM", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fK", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}
