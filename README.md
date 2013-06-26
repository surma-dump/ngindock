`ngindock` is a small helper tool to automatically generate [nginx]
configuration filles according to your currently running [Docker]
containers.

## Installation

	$ go get github.com/surma/ngindock

Or as binary downloads:

* [Linux amd64](http://filedump.surmair.de/binaries/ngindock/linux_amd64/ngindock)
* [Linux 386](http://filedump.surmair.de/binaries/ngindock/linux_386/ngindock)
* [Linux arm](http://filedump.surmair.de/binaries/ngindock/linux_arm/ngindock)
* [Darwin amd64](http://filedump.surmair.de/binaries/ngindock/darwin_amd64/ngindock)
* [Darwin 386](http://filedump.surmair.de/binaries/ngindock/darwin_386/ngindock)

## What does it do?

`ngindock` uses the Docker API to obtain a list of all existing
containers and discards all containers from that list that are not
running or don't have port 80 exposed.

The data of the remaining containers is passed into a template
and can be used to render arbitrary [Go template] files. If no
template is specified, a minimalistic [nginx] configuration is rendered to
`/etc/nginx/conf.d/docker.conf` to make the containers reachable from
the outside under their respective hostname (so make sure you start
your containers with a `-h` flag value).

Afterwards, [nginx] is forced to reload its configuration.

## Usage

	$./ ngindock -h
	Usage: ngindock [global options]

	Global options:
	        -H, --docker      Address of docker daemon (default: localhost:4243)
	        -t, --template    Template to render
	        -o, --output      File to render to (default: /etc/nginx/conf.d/docker.conf)
	            --dont-reload Dont make nginx reload its configuration
	        -h, --help        Show this help

## Example

	$ docker run -d -h 'first-instance.surmair.de' surma/lamp
	$ docker run -d -h 'second-instance.surmair.de' surma/lamp
	$ ./ngindock
	$ cat /etc/nginx/conf.d/docker.conf
	server {
	        listen 80;
	        server_name second-instance.surmair.de;
	        proxy_set_header Host second-instance.surmair.de;


	        location / {
	                proxy_pass http://localhost:49217;
	        }

	}

	server {
	        listen 80;
	        server_name first-instance.surmair.de;
	        proxy_set_header Host first-instance.surmair.de;


	        location / {
	                proxy_pass http://localhost:49215;
	        }

	}

---
Version 1.0.0

[nginx]: http://nginx.org/
[Docker]: http://www.docker.io/
[Go template]: http://golang.org/pkg/text/template/
