package k8s

import (
	v1 "k8s.io/api/core/v1"
	"n2x.dev/x-node/internal/app/node/mnet"
)

type k8sSvcAnnotationsCfg struct {
	// enabled bool
	dnsName string
	valid   bool
}

func parseAnnotations(s *v1.Service) *k8sSvcAnnotationsCfg {
	// var svcEnabled bool
	// var err error

	// if enabled, ok := s.ObjectMeta.Annotations["n2x.io/enabled"]; !ok {
	// 	return nil
	// } else {
	// 	if svcEnabled, err = strconv.ParseBool(enabled); err != nil {
	// 		return nil
	// 	}
	// }

	n := mnet.LocalNode().Node()

	cfgAccountID := n.AccountID
	cfgTenantID := n.TenantID
	cfgNetID := n.Cfg.NetID
	cfgSubnetID := n.Cfg.SubnetID

	valid := true

	accountID, ok := s.ObjectMeta.Annotations["n2x.io/account"]
	if !ok {
		valid = false
	}
	if accountID != cfgAccountID {
		valid = false
	}

	tenantID, ok := s.ObjectMeta.Annotations["n2x.io/tenant"]
	if !ok {
		valid = false
	}
	if tenantID != cfgTenantID {
		valid = false
	}

	netID, ok := s.ObjectMeta.Annotations["n2x.io/network"]
	if !ok {
		valid = false
	}
	if netID != cfgNetID {
		valid = false
	}

	subnetID, ok := s.ObjectMeta.Annotations["n2x.io/subnet"]
	if !ok {
		valid = false
	}
	if subnetID != cfgSubnetID {
		valid = false
	}

	dnsName, ok := s.ObjectMeta.Annotations["n2x.io/dnsName"]
	if !ok {
		dnsName = s.ObjectMeta.Name
	}

	return &k8sSvcAnnotationsCfg{
		// enabled: svcEnabled,
		dnsName: dnsName,
		valid:   valid,
	}
}
