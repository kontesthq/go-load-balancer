package loadbalancer

type IRule interface {
	ChooseServer(lb LoadBalancer) Server
}
