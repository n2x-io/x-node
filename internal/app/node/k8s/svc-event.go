package k8s

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"n2x.dev/x-lib/pkg/errors"
	"n2x.dev/x-lib/pkg/ipnet"
	"n2x.dev/x-lib/pkg/xlog"
	"n2x.dev/x-node/internal/app/node/mnet"
	"n2x.dev/x-node/internal/app/node/mnet/proxy"
)

// manageSvcEvent is the business logic of the controller.
func (c *controller) manageSvcEvent(s *v1.Service, evt eventType) error {
	ns := s.ObjectMeta.Namespace
	svcName := s.ObjectMeta.Name
	svcIP := s.Spec.ClusterIP

	if len(ns) == 0 || len(svcName) == 0 || len(svcIP) == 0 {
		return nil
	}

	if len(s.Spec.Ports) == 0 {
		return nil
	}

	endpointID := fmt.Sprintf("k8s:%s:%s", ns, svcName)

	k8sSvcCfg := &k8sSvcAnnotationsCfg{valid: false}

	if s.ObjectMeta.Annotations != nil {
		k8sSvcCfg = parseAnnotations(s)
	}

	// dnsName := k8sSvcCfg.dnsName
	dnsName := fmt.Sprintf("%s.%s", k8sSvcCfg.dnsName, ns)

	var vIP string
	var err error

	switch evt {
	case eventAdd:
		if k8sSvcCfg.valid {
			vIP, err = mnet.LocalNode().AddNetworkEndpoint(endpointID, dnsName)
			if err != nil {
				xlog.Errorf("Unable to add k8s network endpoint: %v", err)
				return errors.Wrapf(err, "[%v] function netp2p.AddNetworkEndpoint()", errors.Trace())
			}
		}
	case eventUpdate:
		if err = mnet.LocalNode().RemoveNetworkEndpoint(endpointID); err != nil {
			xlog.Errorf("Unable to remove k8s network endpoint: %v", err)
			return errors.Wrapf(err, "[%v] function netp2p.RemoveNetworkEndpoint()", errors.Trace())
		}

		if k8sSvcCfg.valid {
			vIP, err = mnet.LocalNode().AddNetworkEndpoint(endpointID, dnsName)
			if err != nil {
				xlog.Errorf("Unable to add k8s network endpoint: %v", err)
				return errors.Wrapf(err, "[%v] function netp2p.AddNetworkEndpoint()", errors.Trace())
			}
		}
	case eventDelete:
		if err = mnet.LocalNode().RemoveNetworkEndpoint(endpointID); err != nil {
			xlog.Errorf("Unable to remove k8s network endpoint: %v", err)
			return errors.Wrapf(err, "[%v] function netp2p.RemoveNetworkEndpoint()", errors.Trace())
		}
	}

	for _, port := range s.Spec.Ports {
		pName := port.Name
		pProto := ipnet.IPProtocol(string(port.Protocol))
		pPort := port.Port

		switch evt {
		case eventAdd:
			if k8sSvcCfg.valid {
				xlog.Debugf("Adding k8s service %s/%s (ClusterIP: %s): port %s (%v/%d)", ns, svcName, svcIP, pName, pProto, pPort)
				proxy.SetPort(proxy.ServiceTypeKubernetes, ns, svcName, svcIP, vIP, pName, pProto, pPort, ipnet.AddressFamilyIPv4)
			}
		case eventUpdate:
			xlog.Debugf("Updating k8s service %s/%s (ClusterIP: %s): port %s (%v/%d)", ns, svcName, svcIP, pName, pProto, pPort)
			proxy.DeletePort(ns, svcName, pName)
			if k8sSvcCfg.valid {
				proxy.SetPort(proxy.ServiceTypeKubernetes, ns, svcName, svcIP, vIP, pName, pProto, pPort, ipnet.AddressFamilyIPv4)
			}
		case eventDelete:
			xlog.Debugf("Deleting k8s service %s/%s (ClusterIP: %s): port %s (%v/%d)", ns, svcName, svcIP, pName, pProto, pPort)
			proxy.DeletePort(ns, svcName, pName)
		}
	}

	proxy.FwdSvc(ns, svcName)

	return nil
}
