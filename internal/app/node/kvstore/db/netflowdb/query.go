package netflowdb

import (
	"sort"
	"time"

	"n2x.dev/x-api-go/grpc/resources/nstore/netdb"
	"n2x.dev/x-lib/pkg/errors"
)

func (nfdb *netflowDB) Query(r *netdb.TrafficMetricsRequest) (*netdb.TrafficMetricsResponse, error) {
	tmr := &netdb.TrafficMetricsResponse{
		AccountID:       r.Request.AccountID,
		TenantID:        r.Request.TenantID,
		NodeID:          r.Request.NodeID,
		QueryID:         r.Request.QueryID,
		ByProtocol:      nil,
		ByL5Port:        nil,
		ByDirection:     nil,
		ByAddressFamily: nil,
		TopTalkers:      nil,
		Timestamp:       time.Now().UnixMilli(),
	}

	nfl, err := nfdb.Scan()
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function nfdb.Scan()", errors.Trace())
	}

	switch r.Type {
	case netdb.TrafficQueryType_TRAFFIC_BY_PROTOCOL:
		tmr.ByProtocol = getTrafficByProtocol(nfl)
	case netdb.TrafficQueryType_TRAFFIC_BY_L5_PORT:
		tmr.ByL5Port = getTrafficByL5Port(nfl)
	case netdb.TrafficQueryType_TRAFFIC_BY_DIRECTION:
		tmr.ByDirection = getTrafficByDirection(nfl)
	case netdb.TrafficQueryType_TRAFFIC_BY_ADDRESS_FAMILY:
		tmr.ByAddressFamily = getTrafficByAddressFamily(nfl)
	case netdb.TrafficQueryType_TRAFFIC_TOP_TALKERS:
		tmr.TopTalkers = getTrafficTopTalkers(nfl)
	}

	return tmr, nil
}

func getTrafficByProtocol(nfl []*netdb.NetFlowEntry) []*netdb.TrafficByProtocol {
	tbp := make([]*netdb.TrafficByProtocol, 0)

	for _, nfe := range nfl {
		if nfe.Flow == nil || nfe.Flow.Connection == nil {
			continue
		}

		e := &netdb.TrafficByProtocol{
			Timestamp: nfe.Timestamp,
			Traffic:   make(map[string]*netdb.TrafficCounter, 0),
		}

		for _, proto := range protocols() {
			e.Traffic[proto.String()] = &netdb.TrafficCounter{}
		}

		e.Traffic[nfe.Flow.Connection.Proto.String()] = nfe.Traffic

		tbp = append(tbp, e)
	}

	tbp = aggregateByProtocol(tbp)

	sort.Slice(tbp, func(i, j int) bool {
		return tbp[i].Timestamp < tbp[j].Timestamp
	})

	return tbp
}

func getTrafficByL5Port(nfl []*netdb.NetFlowEntry) []*netdb.TrafficByL5Port {
	tbl5p := make([]*netdb.TrafficByL5Port, 0)

	for _, nfe := range nfl {
		if nfe.Flow == nil || nfe.Flow.Connection == nil {
			continue
		}

		e := &netdb.TrafficByL5Port{
			Timestamp: nfe.Timestamp,
			Traffic:   make(map[string]*netdb.TrafficCounter, 0),
		}

		for _, l5port := range l5ports() {
			e.Traffic[l5port.String()] = &netdb.TrafficCounter{}
		}

		e.Traffic[getL5Port(nfe.Flow.Connection).String()] = nfe.Traffic

		tbl5p = append(tbl5p, e)
	}

	tbl5p = aggregateByL5Port(tbl5p)

	sort.Slice(tbl5p, func(i, j int) bool {
		return tbl5p[i].Timestamp < tbl5p[j].Timestamp
	})

	return tbl5p
}

func getTrafficByDirection(nfl []*netdb.NetFlowEntry) []*netdb.TrafficByDirection {
	tbd := make([]*netdb.TrafficByDirection, 0)

	for _, nfe := range nfl {
		if nfe.Flow == nil || nfe.Flow.Connection == nil {
			continue
		}

		e := &netdb.TrafficByDirection{
			Timestamp: nfe.Timestamp,
			Traffic:   make(map[string]*netdb.TrafficCounter, 0),
		}

		for _, dir := range directions() {
			e.Traffic[dir.String()] = &netdb.TrafficCounter{}
		}

		e.Traffic[nfe.Flow.Direction.String()] = nfe.Traffic

		tbd = append(tbd, e)
	}

	tbd = aggregateByDirection(tbd)

	sort.Slice(tbd, func(i, j int) bool {
		return tbd[i].Timestamp < tbd[j].Timestamp
	})

	return tbd
}

func getTrafficByAddressFamily(nfl []*netdb.NetFlowEntry) []*netdb.TrafficByAddressFamily {
	tbaf := make([]*netdb.TrafficByAddressFamily, 0)

	for _, nfe := range nfl {
		if nfe.Flow == nil || nfe.Flow.Connection == nil {
			continue
		}

		e := &netdb.TrafficByAddressFamily{
			Timestamp: nfe.Timestamp,
			Traffic:   make(map[string]*netdb.TrafficCounter, 0),
		}

		for _, af := range addressFamilies() {
			e.Traffic[af.String()] = &netdb.TrafficCounter{}
		}

		e.Traffic[nfe.Flow.Connection.AF.String()] = nfe.Traffic

		tbaf = append(tbaf, e)
	}

	tbaf = aggregateByAddressFamily(tbaf)

	sort.Slice(tbaf, func(i, j int) bool {
		return tbaf[i].Timestamp < tbaf[j].Timestamp
	})

	return tbaf
}

func getTrafficTopTalkers(nfl []*netdb.NetFlowEntry) *netdb.TopTalkers {
	srcTalkersMap := make(map[string]uint64, 0) // map[addr]bytes
	dstTalkersMap := make(map[string]uint64, 0) // map[addr]bytes

	for _, nfe := range nfl {
		if nfe.Flow == nil || nfe.Flow.Connection == nil {
			continue
		}

		switch nfe.Flow.Direction {
		case netdb.ConnectionDirection_INCOMING:
			srcTalkersMap[nfe.Flow.Connection.SrcIP] += nfe.Traffic.Bytes
		case netdb.ConnectionDirection_OUTGOING:
			dstTalkersMap[nfe.Flow.Connection.DstIP] += nfe.Traffic.Bytes
		}
	}

	return &netdb.TopTalkers{
		Src: getTopTalkers(srcTalkersMap),
		Dst: getTopTalkers(dstTalkersMap),
	}
}

func getTopTalkers(talkersMap map[string]uint64) []*netdb.Talker {
	talkers := make([]*netdb.Talker, 0)

	for addr, bytes := range talkersMap {
		talkers = append(talkers, &netdb.Talker{
			Addr:  addr,
			Bytes: bytes,
		})
	}

	sort.Slice(talkers, func(i, j int) bool {
		return talkers[i].Bytes > talkers[j].Bytes
	})

	n := 10
	if len(talkers) < n {
		n = len(talkers)
	}

	topTalkers := make([]*netdb.Talker, n)

	for i := 0; i < n; i++ {
		topTalkers[i] = talkers[i]
	}

	return topTalkers
}

func getL5Port(c *netdb.Connection) netdb.L5Port {
	if c.Proto == netdb.Protocol_TCP && (c.DstPort == 80 || c.SrcPort == 80) {
		return netdb.L5Port_HTTP
	}
	if c.Proto == netdb.Protocol_TCP && (c.DstPort == 8080 || c.SrcPort == 8080) {
		return netdb.L5Port_HTTP
	}
	if c.DstPort == 443 || c.SrcPort == 443 {
		return netdb.L5Port_HTTPS
	}
	if c.Proto == netdb.Protocol_TCP && (c.DstPort == 22 || c.SrcPort == 22) {
		return netdb.L5Port_SSH
	}
	if c.Proto == netdb.Protocol_TCP && (c.DstPort == 3389 || c.SrcPort == 3389) {
		return netdb.L5Port_RDP
	}
	if c.Proto == netdb.Protocol_TCP && (c.DstPort == 445 || c.SrcPort == 445) {
		return netdb.L5Port_SMB
	}
	if c.Proto == netdb.Protocol_TCP && (c.DstPort == 139 || c.SrcPort == 139) {
		return netdb.L5Port_SMB
	}
	// if c.DstPort == 53 || c.SrcPort == 53 {
	// 	return netdb.L5Port_DNS
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 25 || c.SrcPort == 25) {
	// 	return netdb.L5Port_SMTP
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 465 || c.SrcPort == 465) {
	// 	return netdb.L5Port_SMTPS
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 587 || c.SrcPort == 587) {
	// 	return netdb.L5Port_MAIL_SUBMISSION
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 143 || c.SrcPort == 143) {
	// 	return netdb.L5Port_IMAP
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 993 || c.SrcPort == 993) {
	// 	return netdb.L5Port_IMAPS
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 110 || c.SrcPort == 110) {
	// 	return netdb.L5Port_POP3
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 995 || c.SrcPort == 995) {
	// 	return netdb.L5Port_POP3S
	// }
	// if c.Proto == netdb.Protocol_UDP && (c.DstPort == 123 || c.SrcPort == 123) {
	// 	return netdb.L5Port_NTP
	// }
	// if c.DstPort == 161 || c.SrcPort == 161 {
	// 	return netdb.L5Port_SNMP
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 179 || c.SrcPort == 179) {
	// 	return netdb.L5Port_BGP
	// }
	// if c.DstPort == 389 || c.SrcPort == 389 {
	// 	return netdb.L5Port_LDAP
	// }
	// if c.DstPort == 636 || c.SrcPort == 636 {
	// 	return netdb.L5Port_LDAPS
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 3306 || c.SrcPort == 3306) {
	// 	return netdb.L5Port_MYSQL
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 5432 || c.SrcPort == 5432) {
	// 	return netdb.L5Port_POSTGRESQL
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 1433 || c.SrcPort == 1433) {
	// 	return netdb.L5Port_MSSQL
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 6379 || c.SrcPort == 6379) {
	// 	return netdb.L5Port_REDIS
	// }
	// if c.DstPort == 2049 || c.SrcPort == 2049 {
	// 	return netdb.L5Port_NFS
	// }
	// if c.DstPort == 5060 || c.SrcPort == 5060 {
	// 	return netdb.L5Port_SIP
	// }
	// if c.DstPort == 5061 || c.SrcPort == 5061 {
	// 	return netdb.L5Port_SIPTLS
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 5672 || c.SrcPort == 5672) {
	// 	return netdb.L5Port_AMQP
	// }
	// if c.Proto == netdb.Protocol_TCP && (c.DstPort == 5671 || c.SrcPort == 5671) {
	// 	return netdb.L5Port_AMQPS
	// }

	return netdb.L5Port_OTHER_L5PORT
}
