package config

type ETCD struct {
	Endpoints []string `json:"endpoints" yaml:"etcd.endpoints"`
	WatchList []string `json:"watch_list" yaml:"etcd.watchList"`
}
