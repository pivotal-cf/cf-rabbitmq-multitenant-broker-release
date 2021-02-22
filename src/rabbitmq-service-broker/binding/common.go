package binding

func (b Builder) firstHostname() string {
	return b.Hostnames[0]
}

func (b Builder) amqpScheme(tls bool) string {
	if tls {
		return "amqps"
	}
	return "amqp"
}
