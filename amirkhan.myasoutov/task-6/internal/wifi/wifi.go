package wifi

import (
	"fmt"
	"github.com/mdlayher/wifi"
)

type WiFiScanner interface {
	Interfaces() ([]*wifi.Interface, error)
}

type NetworkService struct {
	scanner WiFiScanner
}

func NewNetworkService(s WiFiScanner) NetworkService {
	return NetworkService{scanner: s}
}

func (ns NetworkService) GetInterfaceNames() ([]string, error) {
	ifaces, err := ns.scanner.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	
	names := make([]string, 0, len(ifaces))
	for _, i := range ifaces {
		names = append(names, i.Name)
	}
	return names, nil
}