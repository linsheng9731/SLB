
# SLB (Simple Load Balancer) ver 0.1.0

It's a Simples Load Balancer, just a little project to achieve some kind of performance.

## Features
 * Manage configurations in runtime without downtime
 * High availability (improving with time the speed)
 * Support to WebSockets
 * Monitoring the internal state (improving)
 * Really easy to configure, just a little JSON file

## Next features
 * HTTP/2 support
 * Cache 
 * HTTPS support
 
 If you have any suggestion don't hesitate to open an issue, pull requests are welcome too.

## Install

To install type:

```
go get github.com/linsheng9731/SLB
```

Don't forget to create your configuration file `config.json` at the same directory of project and run it.


## Usage
Type `slb -h` for the command line help


After the configuration file completed you can type only `slb` to start SLB with verbose mode, that command will log the output from SLB in console. That will print something like that:

```
```

## Configuration options

* general:
	* maxProcs: Number of processors used by Go runtime (default: Number of CPUS)
	* workerPoolSize: Number of workers for processing request (default: 10)
	* gracefulShutdown: Wait for the last connection closed, before shutdown (default: true)
	* websocket: Ready for respond websocket connections (default: true)
	* rpchost: Address to expose the internal state (default: 127.0.0.1)
	* rpcport: Port to expose the internal state (default: 42555)
	
* frontends:
	* name: Just a identifier to your front server (required)
	* host: Host address that serves the HTTP front (required)
	* port: Port address that serves the HTTP front (required)
	* route: Route to receive the traffic (required)
	* timeout: How long can wait for the result (ms) from the backend (default: 30000ms)

* backends:
	* name: Just a identifier (required)
	* address: Address (URL) for your backend (required)
	* hearbeat: Addres to send Head request to test if it's ok (required)
	* hbmethod: Method used in request to check the heartbeat (default: HEAD)
	* inactiveAfter: Consider the backend inactive after the number of checks (default: 3)
	* activeAfter: COnsider the backend active after the number of checks (default: 1)
	* heartbeatTime: The interval to send a "ping" (default: 30000ms)
	* retryTime: The interval to send a "ping" after the first failed "ping" (default: 5000ms)
	
### Example (config.json)

```
{
    "general": {
        "maxProcs": 4,
        "workerPoolSize": 10,
    },
    
    "frontends" : [
        {
            "name" : "Front1",
            "host" : "127.0.0.1",
            "port" : 9000,
            "route" : "/",
            "timeout" : 5000,
            
            "backends" : [
                {
                    "name" : "Back1",
                    "address" : "http://127.0.0.1:9001",
                    "heartbeat" : "http://127.0.0.1:9001",
                    "inactiveAfter" : 3,
                    "heartbeatTime" : 5000,
                    "retryTime" : 5000
                },{
                    "name" : "Back2",
                    "address" : "http://127.0.0.1:9002",
                    "heartbeat" : "http://127.0.0.1:9002",
                    "inactiveAfter" : 3,
                    "heartbeatTime" : 5000,
                    "retryTime" : 5000
                }
            ]
        }
    ]
}
```


## LICENSE
Copyright 2016-2017, Buoyant Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
these files except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.