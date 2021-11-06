package utils

import "github.com/hashicorp/consul/api"

type ConsulClient interface {
	Session() *api.Session
}

type ConsulKV interface {
	Acquire(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error)
	Release(p *api.KVPair, q *api.WriteOptions) (bool, *api.WriteMeta, error)
}
type SessionFactory interface {
	Create(se *api.SessionEntry, q *api.WriteOptions) (string, *api.WriteMeta, error)
}
