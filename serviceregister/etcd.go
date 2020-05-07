package serviceregister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/naming"
	gnaming "google.golang.org/grpc/naming"
	"os"
	"time"
)

func GetIport() string {
	etcdIport := os.Getenv("etcdiport")
	if etcdIport == "" {
		etcdIport = "127.0.0.1:12379"
	}
	return etcdIport
}

// RegisterToEtcd 向etcd注册参数
//  - project 项目名
//  - service 服务名
//  - proto 配置文件
//  - iport ip和端口
func RegisterToEtcd(project string, service string, proto string, iport string, content string) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{GetIport()},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r := &naming.GRPCResolver{Client: cli}
	//metaData := map[string]interface{}{
	//	"proto": proto,
	//}
	err = r.Update(context.TODO(), fmt.Sprintf("%s-%s", project, service), gnaming.Update{Op: gnaming.Add, Addr: iport})
	if err != nil {
		fmt.Println(err.Error())
	}
	data := map[string]string{
		"iport":   iport,
		"proto":   proto,
		"content": content,
	}
	d, _ := json.Marshal(data)
	cli.Put(ctx, project+"-"+service, string(d))
	return nil
}
