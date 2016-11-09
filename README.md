# summar

Summar is a utility for summarizing nginx access logs. Every five seconds new entries
are summarized into a statsd-compatible message and appended to an output file.

## Requirements
* Docker
* go (~1.73) with a working go workspace

## Installation

To fetch and build:
```
go get github.com/hankjacobs/summar
```

Then change to summar's working directory:
```
cd $GOPATH/src/github.com/hankjacobs/summar
```

## Usage
### Locally
Normally, summar looks at `/var/log/nginx/access.log` for input and outputs metrics to `/var/log/stats.log`.
For local testing, we can change this using the `-in` and `-out` flags

To run using simulated data, run summar using:

```
$GOPATH/bin/summar -in access.log -out stats.log
```

then start tail in another terminal:

```
tail -f stats.log
```

and copy the test data line by line into access.log:

```
while read line; do echo $line >> access.log; done < testdata/sample50x.log
```
### Docker

Docker is used to simulate a nginx web server serving real traffic.

The docker image sets up nginx with routes that produce 200, 300, 400, and 500 response status codes. A request generator is configured to randomly invoke those paths every 10 milliseconds. Summar is then used to summarize `/var/log/nginx/access.log` into `/var/log/stats.log`. `tail -f /var/log/stats.log` is used to tail the metrics so that they are visible to a console that is running the docker image in foreground mode.

Build the docker image using:
```
docker build -t summar .
```

Run the docker image (in foreground mode) using:
```
docker run -it --rm --name summar_dsc_test summar
```

Every five seconds, you should see the contents of `/var/log/stats.log` output to your console.
