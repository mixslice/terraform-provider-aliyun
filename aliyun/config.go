package aliyun

type Config struct {
	APIKey     string
	Endpoint   string
	Timeout    int
	MaxRetries int
}

type Machine struct {
	Name string
	CPUs int
	RAM  int
}

func (m *Machine) Id() string {
	return "id-" + m.Name + "!"
}

func (c *Config) CreateMachine(m *Machine) error {
	return nil
}
