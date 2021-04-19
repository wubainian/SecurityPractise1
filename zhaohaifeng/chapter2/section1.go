package chapter2

import (
	"net"
	"fmt"
	"time"
	"github.com/malfunkt/iprange"
	"log"
	"strings"
	"strconv"
	"os"
)

func Connect(ip string, port int)(net.Conn, error){
	//tcp 全连接端口扫描
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 2 * time.Second)
	defer func(){
		if conn != nil {
			conn.Close()
		}
	}()
	return conn, err
}

func IprangeTest(){
	lists, err := iprange.ParseList("10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24")
	if err != nil {
		log.Printf("error : %s", err)
	}
	log.Printf("lists = %+v", lists)
	rng := lists.Expand()
	log.Printf("rng = %v", rng)
}

func GetIpList(ips string)([]net.IP, error){
	//封装iprange处理
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	list := addressList.Expand()
	return list, nil
}

func GetPorts(selection string)([]int, error){
	//多端口处理支持，-
	ports := []int{}
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment : '%s'", r)
			}
			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number '%s'", parts[0])
			}
			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number '%s'", parts[1])
			}
			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range '%d - %d'", p1, p2)
			}

			for i:=p1; i<p2; i++ {
				ports = append(ports, i)
			}
		}else{
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number : '%s'", r)
			}else{
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}

func Scan1(){
	if len(os.Args) == 3 {
		ipList := os.Args[1]
		portList := os.Args[2]
		ips, _ := GetIpList(ipList)
		ports, _ := GetPorts(portList)
		for _, ip := range ips {
			for _, port := range ports {
				_, err := Connect(ip.String(), port)
				if err != nil {
					continue
				}
				fmt.Printf("ip:%v, port:%v is open\n", ip, port)
			}
		}
	}
}