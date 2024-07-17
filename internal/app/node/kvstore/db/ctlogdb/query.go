package ctlogdb

import (
	"time"

	"n2x.dev/x-api-go/grpc/resources/nstore/netdb"
	"n2x.dev/x-lib/pkg/errors"
)

func (ndb *ctlogDB) Query(r *netdb.ConntrackLogRequest) (*netdb.ConntrackLogResponse, error) {
	ctlr := &netdb.ConntrackLogResponse{
		AccountID: r.Request.AccountID,
		TenantID:  r.Request.TenantID,
		NodeID:    r.Request.NodeID,
		QueryID:   r.Request.QueryID,
		CtLog:     make([]*netdb.ConntrackLogEntry, 0),
		Timestamp: time.Now().UnixMilli(),
	}

	el, err := ndb.Scan()
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function ndb.Scan()", errors.Trace())
	}

	ctlr.CtLog = el

	return ctlr, nil
}
