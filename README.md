
# SLB (Simple Load Balancer) ver 0.1.0

It's a Simple Load Balancer, inspired by sslb(https://github.com/eduardonunesp/sslb), a Super Simple Loader Balancer.
SLB improved sslb , so it's no super simple anymore, just simple.

## Features
 * Http proxy
 * Rest API
 * Single binary with no other dependencies for easy deployment
 * Dynamic reloading without restart (SLB reload)
 * Really easy to configure, just a little JSON file
 * Support to WebSockets

## Install

To install:

```
go get github.com/linsheng9731/SLB
cd $GOPATH/src/github.com/linsheng9731/SLB
go install
```

Don't forget to create your configuration file `config.json` at the same directory of project and run it.

### Example (config.json)

```
{
  "general": {
    "maxProcs": 4,
    "workerPoolSize": 1
  },

  "frontends" : [
    {
      "name" : "Front1",
      "host" : "127.0.0.1",
      "port" : 9000,
      "route" : "/dir",
      "timeout" : 5000,
      "backends" : [
        {
          "name" : "Back1",
          "address" : "http://127.0.0.1:9001",
          "heartbeat" : "http://127.0.0.1:9001",
          "inactiveAfter" : 3,
          "heartbeatTime" : 5000,
          "retryTime" : 5000
        }
      ]
    },
    {
      "name" : "Front2",
      "host" : "127.0.0.1",
      "port" : 9003,
      "route" : "/",
      "timeout" : 5000,
      "backends" : [
        {
          "name" : "Back1",
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
## Usage
Type `slb -h` for the command line help


After the configuration file completed you can type only `slb` to start SLB with verbose mode, that command will log the output from SLB in console. That will print something like that:

```
2017/09/26 17:54:45 Start SLB (LbServer)
2017/09/26 17:54:45 Create worker pool with [1]
2017/09/26 17:54:45 Prepare to run server ...
2017/09/26 17:54:45 Setup and check configuration
2017/09/26 17:54:45 Setup ok ...
2017/09/26 17:54:45 Api server listen on 127.0.0.1:9292
2017/09/26 17:54:45 Start frontend http server [Front1] at [127.0.0.1:9000]
2017/09/26 17:54:45 Start frontend http server [Front2] at [127.0.0.1:9003]
2017/09/26 17:54:45 Backend active [Back1]
```

## Configuration options

* general:
	* maxProcs: Number of processors used by Go runtime (default: Number of CPUS)
	* workerPoolSize: Number of workers for processing request (default: 10)
	* gracefulShutdown: Wait for the last connection closed, before shutdown (default: true)
	* websocket: Ready for respond websocket connections (default: true)
	* rpchost: Address to expose the internal state (default: 127.0.0.1)
	* rpcport: Port to expose the internal state (default: 42555)
	* apihost: Http api address (default: 127.0.0.1)
	* apiport: Http api address port (default: 9292)

* frontends:
	* name: Just a identifier to your front server (required)
	* host: Host address that serves the HTTP front (required)
	* port: Port address that serves the HTTP front (required)
	* route: Route to receive the traffic (required)
	* timeout: How long can wait for the result (ms) from the backend (default: 30000ms)

* backends:
	* name: Just a identifier (required)
	* address: Address (URL) for your backend (required)
	* hearbeat: Address to send Head request to test if it's ok (required)
	* hbmethod: Method used in request to check the heartbeat (default: HEAD)
	* inactiveAfter: Consider the backend inactive after the number of checks (default: 3)
	* activeAfter: Consider the backend active after the number of checks (default: 1)
	* heartbeatTime: The interval to send a "ping" (default: 30000ms)
	* retryTime: The interval to send a "ping" after the first failed "ping" (default: 5000ms)

## Rest API
* http://apihost:apiport/health-check (deafult: http://127.0.0.1:9292/health-check)
* http://apihost:apiport/config (deafult: http://127.0.0.1:9292/config)
* http://apihost:apiport/status (deafult: http://127.0.0.1:9292/status)

# Road map
## v0.2
 * Configurations check
 * Internal status and metrics http api
 * WebUI

## v0.3
 * Cache
 * HTTP/2 support
 * HTTPS support

## v0.4
 * Integration marathon as backend

 If you have any suggestion don't hesitate to open an issue, pull requests are welcome too.


## LICENSE
Copyright (c) 2017, Lin Shengsheng
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of slb nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.