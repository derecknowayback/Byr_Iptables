package iptables

import (
	"github.com/coreos/go-iptables/iptables"
	"net"
)

const (
	TABLE     = "NAT"
	DnatChain = "PREROUTING"
	SnatChain = "POSTROUTING"
)

// iptables -t nat -A PREROUTING -d 10.11.0.0/16 -j DNAT --to-destination 192.168.0.0/16
// iptables -t nat -A POSTROUTING -s 192.168.0.0/16 -j SNAT --to-source 10.11.0.0/16

var Iptables *iptables.IPTables

// InitIptables 初始化iptables
func InitIptables() error {
	ip4t, err := iptables.New() // ipv4版本, timeout 0
	Iptables = ip4t
	return err
}

func Dnat(ip1, ip2 net.IP) error {
	err := Iptables.Append(TABLE, DnatChain, getDnatArgs(ip1, ip2)...)
	return err
}

func deleteDnat(table, chain string, ip1, ip2 net.IP) error {
	err := Iptables.DeleteIfExists(table, chain, getDnatArgs(ip1, ip2)...)
	return err
}

// -d 10.11.0.0/16 -j DNAT --to-destination 192.168.0.0/16
func getDnatArgs(ip1, ip2 net.IP) []string {
	args := []string{
		"-d",
		ip1.String(),
		"-j",
		"DNAT",
		"--to-destination",
		ip2.String(),
	}
	return args
}

func Snat(ip1, ip2 net.IP) error {
	err := Iptables.Append(TABLE, SnatChain, getSnatArgs(ip1, ip2)...)
	return err
}

func deleteSnat(table, chain string, ip1, ip2 net.IP) error {
	err := Iptables.DeleteIfExists(table, chain, getSnatArgs(ip1, ip2)...)
	return err
}

// -s 192.168.0.0/16 -j SNAT --to-source 10.11.0.0/16
func getSnatArgs(ip1, ip2 net.IP) []string {
	args := []string{
		"-s",
		ip1.String(),
		"-j",
		"SNAT",
		"--to-source",
		ip2.String(),
	}
	return args
}