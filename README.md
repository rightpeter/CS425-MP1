# CS425-MP1

## How to run this project

In /CS425-MP1 folder:

## Build

Use

	go build -o server.bin ./server && go build -o client.bin ./client

to build excecutable files.

## Run

Server:

	./server.bin -p <PORT_NUMBER> -ip <IP_ADDRESS> -c <CONFIG_FILE_PATH>

	// IP flag is currently optional

Client:

	./client.bin -i -grep "GREP_PATTERNS"


# Test

	1) [run from local computer] Run script to generate random log files with known patterns and send to vms
		python pre_test_script.py

	2) ssh into vms and run server with test config file:
		./server.bin -c /tmp/mp1.test.config.json 
	
	3) [run from any vm] Run unit tests
		python test_script.py
