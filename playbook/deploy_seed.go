package playbook

import (
	"fmt"

	"gitee.com/wisecloud/wise-deploy/database"
)

type DeploySeed2 struct {
	Registry     *Component
	Etcd         *Component
	MySQL        *Component
	LoadBalancer *LoadBalancer
	K8sMaster    *Component
	K8sNode      *Component
	WiseCloud    *Component
}

func (d *DeploySeed2) RegistryEndpoint() string {
	ip := d.LoadBalancer.getEndpoint("registry")
	if ip != "" {
		return ip
	}

	if len(d.Registry.Hosts) > 0 {
		return d.Registry.Hosts[0].IP
	}

	return ""
}

func (d *DeploySeed2) MySQLEndpoint() string {
	ip := d.LoadBalancer.getEndpoint("mysql")
	if ip != "" {
		return ip
	}

	if len(d.MySQL.Hosts) > 0 {
		return d.MySQL.Hosts[0].IP
	}

	return ""
}

func (d *DeploySeed2) K8sEndpoint() string {
	ip := d.LoadBalancer.getEndpoint("k8s")
	if ip != "" {
		return ip
	}

	if len(d.K8sMaster.Hosts) > 0 {
		return d.K8sMaster.Hosts[0].IP
	}

	return ""
}

type Component struct {
	Property map[string]interface{}
	Hosts    []*database.Host
}

func NewDeploySeed(c *database.Cluster) *DeploySeed2 {
	hs := make(map[string]*database.Host)
	for _, h := range c.Hosts {
		hs[h.ID] = h
	}

	ds := &DeploySeed2{
		Registry:     &Component{},
		Etcd:         &Component{},
		MySQL:        &Component{},
		LoadBalancer: &LoadBalancer{},
		K8sMaster:    &Component{},
		K8sNode:      &Component{},
		WiseCloud:    &Component{},
	}

	for _, cp := range c.Components {
		switch cp.Name {
		case "etcd":
			setComponentHost(hs, cp, ds.Etcd)
		case "registry":
			setComponentHost(hs, cp, ds.Registry)
		case "mysql":
			setComponentHost(hs, cp, ds.MySQL)
		case "loadbalancer":
			setComponentHost(hs, cp, (*Component)(ds.LoadBalancer))
		case "k8smaster":
			setComponentHost(hs, cp, ds.K8sMaster)
		case "k8snode":
			setComponentHost(hs, cp, ds.K8sNode)
		case "wisecloud":
			setComponentHost(hs, cp, ds.WiseCloud)
		}
	}

	return ds
}

func setComponentHost(
	hs map[string]*database.Host,
	sourceCp *database.Component,
	destCp *Component,
) {
	destCp.Property = sourceCp.Property
	for _, h := range sourceCp.Hosts {
		th, ok := hs[h]
		if !ok {
			panic(fmt.Errorf("unexpected host: %s", h))
		}
		destCp.Hosts = append(destCp.Hosts, th)
	}
}
