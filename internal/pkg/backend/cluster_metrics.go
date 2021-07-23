package backend

func (p *Cluster) byLowConnections() *Client {
	client, _ := p.byLowConnectionsWithMetrics()
	return client
}

func (p *Cluster) byLowConnectionsWithMetrics() (*Client, float32) {
	var count = float32(9999999)
	var client *Client
	for _, s := range p.ro {
		var idleCount = float32(s.pool.Stat().AcquiredConns()) * 0.6
		if idleCount < count {
			count = idleCount
			client = s
		}
	}
	var idleCount = float32(p.rw.pool.Stat().AcquiredConns())
	if idleCount < count {
		count = idleCount
		client = p.rw
	}

	return client, count
}
