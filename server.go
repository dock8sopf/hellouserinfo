package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "hellouserinfo/protofiles"
	"hellouserinfo/serviceregister"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

const (
	port = "0.0.0.0:50052"
)

type server struct{}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

// GetUserInfo 返回hello
//  - ctx 上下文参数
//  - in grpc请求的入参
func (s *server) GetUserInfo(ctx context.Context, in *pb.UserReq) (*pb.UserResp, error) {
	fmt.Printf("这里收到请求了: %#v", in)
	if in.Id == 1 {
		resp := &pb.UserResp{
			Id:      "1",
			Name:    "陆隐峰",
			Gender:  1,
			Address: "广东省东莞市",
		}
		return resp, nil
	} else {
		return nil, errors.New("user not found")
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserInfoServer(s, &server{})
	f, err := ioutil.ReadFile("protofiles/hellouserinfo.proto")
	if err != nil {
		fmt.Println("read fail", err)
	}
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	iport := ip.String() + ":" + strings.Split(port, ":")[1]
	registerErr := serviceregister.RegisterToEtcd("hellouserinfo", "UserInfo", "hellouserinfo.proto", iport, string(f))
	if registerErr != nil {
		log.Fatalf("failed to serve: etcd(%v)", err)
	}
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
