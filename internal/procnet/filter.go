package procnet

var validNInterfaces = map[string]bool{
	"enp4s0": true, // Ethernet
	"wlp3s0": true, // WiFi
}

type FilterFunc func(nif string) (ignore bool)

func ValidNIsFunc(nif string) bool {
	return validNInterfaces[nif]
}
