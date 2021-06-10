package ConsulClient

import (
	"errors"
	"fmt"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	defaultPort = "8500"
	ConsulAdress = "127.0.0.1:8500"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errAddrMisMatch = errors.New("consul resolver: invalied uri")

	errEndsWithColon = errors.New("consul resolver: missing port after port-separator colon")

	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")

	//单例模式
	builderInstance = &consulBuilder{}
)

func Init() {
	fmt.Printf("calling consul init\n")
	resolver.Register(CacheBuilder())
}

type consulBuilder struct {
}

type consulResolver struct {
	address              string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	Ch                   chan int
}

func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

func CacheBuilder() resolver.Builder {
	return builderInstance
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	host, port, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	if err != nil {
		fmt.Println("parse err")
		return nil, err
	}
	fmt.Println(fmt.Sprintf("consul service ==> host:%s, port%s, name:%s",host, port, name))
	cr := &consulResolver{
		address:              fmt.Sprintf("%s%s", host, port),
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		Ch:					  make(chan int, 0),
	}
	go cr.watcher()
	return cr, nil

}

func (cr *consulResolver) watcher() {
	fmt.Printf("calling [%s] consul watcher\n", cr.name)
	config := api.DefaultConfig()
	config.Address = cr.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}
	t := time.NewTicker(2000 * time.Millisecond)
	defer func() {
		fmt.Println("defer done")
	}()
	for {
		select {
		case <-t.C:
			//fmt.Println("定时")
		case <-cr.Ch:
			//fmt.Println("ch call")
		}
		//api添加了 lastIndex   consul api中并不兼容附带lastIndex的查询
		services, _, err := client.Health().Service(cr.name, "", true, &api.QueryOptions{})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		}

		newAddrs := make([]resolver.Address, 0)
		for _, service := range services {
			addr := net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))
			newAddrs = append(newAddrs, resolver.Address{
				Addr: addr,
				//type：不能是grpclb，grpclb在处理链接时会删除最后一个链接地址，不用设置即可 详见=> balancer_conn_wrappers => updateClientConnState
				ServerName:service.Service.Service,
			})
		}
		//cr.cc.NewAddress(newAddrs)
		//cr.cc.NewServiceConfig(cr.name)
		cr.cc.UpdateState(resolver.State{Addresses:newAddrs})
	}

}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	cr.Ch <- 1
}

func (cr *consulResolver) Close() {
}

func parseTarget(target string) (host, port, name string, err error) {

	if target == "" {
		return "", "", "", errMissingAddr
	}

	if !regexConsul.MatchString(target) {
		return "", "", "", errAddrMisMatch
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	if port == "" {
		port = defaultPort
	}
	return host, port, name, nil
}

func GetConsulHost() string {
	return ConsulAdress
}

type GrpcClient struct {
	Conn 		*grpc.ClientConn
	RpcTarget   string
	Name   		string
}

func (s *GrpcClient)RunGrpcClient(){
	conn, err := grpc.Dial(s.RpcTarget, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Conn = conn
	fmt.Println("grpc client start success")
}
func (s *GrpcClient)RunConsulClient(){
	//初始化 resolver 实例
	Init()
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s/%s", "consul", GetConsulHost(), s.Name),
		//不能block => blockkingPicker打开，在调用轮询时picker_wrapper => picker时若block则不进行robin操作直接返回失败
		//grpc.WithBlock(),
		grpc.WithInsecure(),
		//指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
		//grpc.WithBalancerName(roundrobin.Name),
		//grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
	)

	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
	s.Conn = conn
	fmt.Println(fmt.Sprintf("gRpc consul client [%s] start success", s.Name))
}
