package pkg

// CduleConfig cdule configuration
type CduleConfig struct {
	Cduletype        string `yaml:"cduletype"`
	Dburl            string `yaml:"dburl"` // underscore creates the problem for e.f. db_url, so should be avoided
	Cduleconsistency string `yaml:"cduleconsistency"`
}
