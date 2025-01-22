package parsers

import "net"

type IPv4Parser struct{}

func (IPv4Parser) Name() string {
    return "ipv4"
}

func (IPv4Parser) Parse(s string) (net.IP, error) {
    return net.ParseIP(s).To4(), nil
}
