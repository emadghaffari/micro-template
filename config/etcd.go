package config

type ETCD struct {
	Endpoints []string `json:"endpoints" yaml:"etcd.endpoints"`
	WatchList []string `json:"watch_list"`
	Username  string   `json:"username" yaml:"etcd.username"`
	Password  string   `json:"password" yaml:"etcd.password"`
}
