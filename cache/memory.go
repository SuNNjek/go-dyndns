package cache

import "net"

type memoryCache struct {
	value net.IP
}

func (m *memoryCache) GetLastIp() (net.IP, error) {
	return m.value, nil
}

func (m *memoryCache) SetLastIp(ip net.IP) error {
	m.value = ip
	return nil
}
