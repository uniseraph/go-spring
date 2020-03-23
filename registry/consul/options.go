package consul

import (
	"context"
	consul "github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/registry"
)

// Connect specifies services should be registered as Consul Connect services
func Connect() registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_connect", true)
	}
}

func Config(c *consul.Config) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_config", c)
	}
}

// AllowStale sets whether any Consul server (non-leader) can service
// a read. This allows for lower latency and higher throughput
// at the cost of potentially stale data.
// Works similar to Consul DNS Config option [1].
// Defaults to true.
//
// [1] https://www.consul.io/docs/agent/options.html#allow_stale
//
func AllowStale(v bool) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_allow_stale", v)
	}
}

// QueryOptions specifies the QueryOptions to be used when calling
// Consul. See `Consul API` for more information [1].
//
// [1] https://godoc.org/github.com/hashicorp/consul/api#QueryOptions
//
func QueryOptions(q *consul.QueryOptions) registry.Option {
	return func(o *registry.Options) {
		if q == nil {
			return
		}
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_query_options", q)
	}
}



//
// HttpCheck will tell the service provider to check the service address
// and port every `t` interval. It will enabled only if `t` is greater than 0.
// See `TCP + Interval` for more information [1].
//
// [1] https://www.consul.io/docs/agent/checks.html
//http check
//{
//  "check": {
//    "id": "api",
//    "name": "HTTP API on port 5000",
//    "http": "https://localhost:5000/health",
//    "tls_skip_verify": false,
//    "method": "POST",
//    "header": {"Content-Type": "application/json"},
//    "body": "{\"method\":\"health\"}",
//    "interval": "10s",
//    "timeout": "1s"
//  }
//}
//grpc check
//{
//  "check": {
//    "id": "mem-util",
//    "name": "Service health status",
//    "grpc": "127.0.0.1:12345/my_service",
//    "grpc_use_tls": true,
//    "interval": "10s"
//  }
//}
//




func AgentServiceCheck( check *consul.AgentServiceCheck ) registry.Option {
	return func(o *registry.Options) {
		if check==nil {
			return
		}
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_agent_service_check", check)
	}
}
