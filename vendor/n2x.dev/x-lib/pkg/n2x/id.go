package n2x

import (
	"fmt"
	"strings"

	"n2x.dev/x-api-go/grpc/resources/topology"
	"n2x.dev/x-lib/pkg/errors"
)

type ID interface {
	String() string
	Node() (*topology.NodeReq, error)
	AccountID() string
	TenantID() string
	NetID() string
	SubnetID() string
	NodeID() string

	isValid() error
}
type n2xid string

func GetID(nr *topology.NodeReq) (ID, error) {
	if nr == nil {
		return nil, fmt.Errorf("invalid data")
	}

	if len(nr.AccountID) == 0 {
		return nil, fmt.Errorf("missing accountID")
	}

	if len(nr.TenantID) == 0 {
		return nil, fmt.Errorf("missing tenantID")
	}

	// if len(nr.NetID) == 0 {
	// 	return nil, fmt.Errorf("missing netID")
	// }

	// if len(nr.SubnetID) == 0 {
	// 	return nil, fmt.Errorf("missing subnetID")
	// }

	if len(nr.NodeID) == 0 {
		return nil, fmt.Errorf("missing nodeID")
	}

	// return n2xid(fmt.Sprintf("%s:%s:%s:%s:%s",
	// 	nr.AccountID, nr.TenantID, nr.NetID, nr.SubnetID, nr.NodeID)), nil
	return n2xid(fmt.Sprintf("%s:%s:%s", nr.AccountID, nr.TenantID, nr.NodeID)), nil
}

func ParseID(id string) (ID, error) {
	if err := n2xid(id).isValid(); err != nil {
		return nil, errors.Wrapf(err, "[%v] function n2xid(id).isValid()", errors.Trace())
	}

	return n2xid(id), nil
}

func ParseCLIID(id string) (ID, error) {
	if err := n2xid(id).isValidCLIID(); err != nil {
		return nil, errors.Wrapf(err, "[%v] function n2xid(id).isValidCLIID()", errors.Trace())
	}

	return n2xid(id), nil
}

func (id n2xid) String() string {
	return string(id)
}

func (id n2xid) Node() (*topology.NodeReq, error) {
	if err := id.isValid(); err != nil {
		return nil, errors.Wrapf(err, "[%v] function id.isValid()", errors.Trace())
	}

	return &topology.NodeReq{
		AccountID: id.AccountID(),
		TenantID:  id.TenantID(),
		// NetID:     id.NetID(),
		// SubnetID:  id.SubnetID(),
		NodeID: id.NodeID(),
	}, nil
}

func (id n2xid) AccountID() string {
	return strings.Split(id.String(), ":")[0]
}

func (id n2xid) TenantID() string {
	return strings.Split(id.String(), ":")[1]
}

func (id n2xid) NetID() string {
	return strings.Split(id.String(), ":")[2]
}

func (id n2xid) SubnetID() string {
	return strings.Split(id.String(), ":")[3]
}

func (id n2xid) NodeID() string {
	return strings.Split(id.String(), ":")[4]
}

func (id n2xid) isValid() error {
	s := id.String()

	if len(s) == 0 {
		return fmt.Errorf("missing id")
	}

	if len(strings.Split(s, ":")) != 5 {
		return fmt.Errorf("invalid format")
	}

	if len(id.AccountID()) == 0 {
		return fmt.Errorf("missing accountID")
	}

	if len(id.TenantID()) == 0 {
		return fmt.Errorf("missing tenantID")
	}

	if len(id.NetID()) == 0 {
		return fmt.Errorf("missing netID")
	}

	if len(id.SubnetID()) == 0 {
		return fmt.Errorf("missing subnetID")
	}

	if len(id.NodeID()) == 0 {
		return fmt.Errorf("missing nodeID")
	}

	return nil
}

// CLI ID

func (id n2xid) cliIDPrefix() string {
	return strings.Split(id.String(), ":")[1]
}

func (id n2xid) cliIDHostID() string {
	return strings.Split(id.String(), ":")[2]
}

func (id n2xid) cliIDGID() string {
	return strings.Split(id.String(), ":")[3]
}

func (id n2xid) cliIDTimestamp() string {
	return strings.Split(id.String(), ":")[4]
}

func (id n2xid) isValidCLIID() error {
	s := id.String()

	if len(s) == 0 {
		return fmt.Errorf("missing id")
	}

	if len(strings.Split(s, ":")) != 5 {
		return fmt.Errorf("invalid format")
	}

	if len(id.AccountID()) == 0 {
		return fmt.Errorf("missing accountID")
	}

	if len(id.cliIDPrefix()) == 0 {
		return fmt.Errorf("missing cliID prefix")
	}

	if len(id.cliIDHostID()) == 0 {
		return fmt.Errorf("missing cliID hostID")
	}

	if len(id.cliIDGID()) == 0 {
		return fmt.Errorf("missing cliID GID")
	}

	if len(id.cliIDTimestamp()) == 0 {
		return fmt.Errorf("missing cliID timestamp")
	}

	return nil
}
