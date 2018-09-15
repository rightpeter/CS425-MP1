package model

// RPCArgs Input Arguments for rpc
type RPCArgs struct {
	Commands []string
}

// RPCResult Result for rpc
type RPCResult struct {
	ClientID int
	Reply    string
	Alive    bool
	Error    error
}

type nodeConfig struct {
	ID      int    `json:"id"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	LogPath string `json:"log_path"`
}

// NodesConfig structure to unmarshal json config file {id: int, ip: string, port: int}
type NodesConfig struct {
	Current nodeConfig   `json:"current"`
	Nodes   []nodeConfig `json:"nodes"`
}
