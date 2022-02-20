package pkg

type CduleConfig struct {
	Cduletype        string `yaml:"cduletype"`
	Dburl            string `yaml:"dburl"`
	Cduleconsistency string `yaml:"cduleconsistency"`
	WorkerHostIP     string `yaml:"workerhostip"` // underscore creates the problem for e.f. worker_host_ip, so should be avoided
	WorkerPort       string `yaml:"workerport"`
}
