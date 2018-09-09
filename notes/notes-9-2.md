# Note 9-2

## System design

Hard code ID and IP in config files. Every node config files like

	{
	  "current": { "id": 1, "ip": "192.168.0.1", "port": 8080 },
	  "nodes": [
	    {
		"id": 1,
		"ip": "192.168.0.1",
		"port": "8080"
	    },
	    {
		"id": 2,
		"ip": "192.168.0.2",
		"port": "8080"
	    }
	  ]
	}

Every node will have a deamon server process to handle the grep request from other nodes.
There will also be a command line tool for user to grep log from all alive nodes.
Below is how the network would look like.

<br>

![Semantic description of image](./structure.svg "Structure")

## File Structure

Below is the file structure of the repo.

	├── README.md
	├── client
	│   └── client.go
	├── config.json
	├── logs
	│   └── example.logs
	├── model
	│   └── model.go
	├── notes
	│   └── example.md
	├── server
	│   ├── grep
	│   │   └── grep.go
	│   └── server.go
	└── tests
