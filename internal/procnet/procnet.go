package procnet

type Stats struct {
	NInterface string // Network Interface
	RxBytes    uint64 // Recieved bytes
	TxBytes    uint64 // Transmitted bytes
}

func Diff(s1, s2 *Stats) (uint64, uint64) {
	return (s1.RxBytes - s2.RxBytes), (s1.TxBytes - s2.TxBytes)
}
