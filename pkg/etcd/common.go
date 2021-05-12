package etcd

import "os"

func GetIport() string {
	etcdIport := os.Getenv("etcdiport")
	if etcdIport == "" {
		etcdIport = "127.0.0.1:12379"
	}
	return etcdIport
}
