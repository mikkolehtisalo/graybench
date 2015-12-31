# graybench
Graylog2 HTTP/Gelf benchmark tool. Assumes TLS is enabled. 

## Usage

```
Usage of ./graybench:
  -ca string
    	file with ca certificate chain (default "cert.pem")
  -customlen int
    	length of random custom message (default 20)
  -events int
    	events per thread (default 100000)
  -fulllen int
    	length of random full message (default 200)
  -shortlen int
    	length of random short message (default 20)
  -target string
    	target HTTP Gelf service (default "https://graylog.local:12201/gelf")
  -threads int
    	amount of threads (default 10)
```

## Example

Graylog 1.3 Virtual Machine Appliance, Linux KVM, 4* CPU (i7-4710HQ), 4096M, without major setting changes:

```
$ ./graybench
2015-12-31T18:15:42+02:00: Threads: 10, Events: 100000, CA certificate chain: cert.pem, Target: https://graylog.local:12201/gelf
2015-12-31T18:15:42+02:00: Launching threads
2015-12-31T18:18:12+02:00: Threads finished
2015-12-31T18:18:12+02:00: Total events: 1000000, Total time: 149s, EPS: 6673
```

Event source is benchmarkhost, and the events attempt to mimick realistic data amounts (tune with command line settings).

