<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/sreejeet/garagesale">&nbsp;<img src="https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg">
[![Go Report Card](https://goreportcard.com/badge/github.com/sreejeet/garagesale)](https://goreportcard.com/report/github.com/sreejeet/garagesale)

### (Work in progress)
<img alt="Image" src="https://i.imgur.com/5K6jBOC.png">

# Production ready RESTful API service.

Garagesale is a production ready RESTful API service running on docker. It is build without any framework to keep the service as light as possible. The final product should resemble [The Ultimate Service](https://github.com/ardanlabs/service)

This is a product of my work while training for Go based web services.  
The open source [training material](https://github.com/ardanlabs/service-training) is provided by [Ardan Labs](http://www.ardanlabs.com/).

## Tech stack
1. Go's [net/http](https://golang.org/pkg/net/http/)
2. [Docker](https://www.docker.com)
3. [PostgreSQL](https://www.postgresql.org/)

## A consolidated list of resources I found useful.
(Please raise an issue if you find broken links)
1. [The complete guide to Go net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
2. [So you want to expose Go on the Internet](https://blog.cloudflare.com/exposing-go-on-the-internet/)
3. [Parsing JSON files With Golang](https://tutorialedge.net/golang/parsing-json-with-golang/)
4. [Go database/sql tutorial](http://go-database-sql.org/)
5. [Docker Compose](https://docs.docker.com/compose/compose-file)
6. [Docker logs](https://docs.docker.com/config/containers/logging/)
7. [Package Oriented Design](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html)
8. [Package conf](https://github.com/ardanlabs/service-training/blob/master/06-configuration/internal/platform/conf/README.md)
9. [Error Handling in Go](https://medium.com/@hussachai/error-handling-in-go-a-quick-opinionated-guide-9199dd7c7f76)
10. [Package initialization and program execution order](https://yourbasic.org/golang/package-init-function-main-execution-order/)
11. [How to collect, standardize, and centralize Golang logs](https://www.datadoghq.com/blog/go-logging/)
12. [Error handling and Go](https://blog.golang.org/error-handling-and-go)
13. [Structs and Interfaces](https://www.golang-book.com/books/intro/9) (This is part of a book from 2012. Outdated, but a good read nonetheless.)
14. [Understanding the context package in golang](http://p.agnihotry.com/post/understanding_the_context_package_in_golang/)
15. [Go: Context and Cancellation by Propagation](https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c)
16. [PostgreSQL Joins](https://www.postgresqltutorial.com/postgresql-joins/)
17. [PostgreSQL SUM Function](https://www.postgresqltutorial.com/postgresql-sum-function/)
18. [Profiling Go Programs](https://blog.golang.org/pprof)
19. [Middleware (Advanced)](https://gowebexamples.com/advanced-middleware/)
20. [Creating a Middleware in Golang for JWT based Authentication](https://hackernoon.com/creating-a-middleware-in-golang-for-jwt-based-authentication-cx3f32z8)
21. [How to instrument Go code with custom expvar metrics](https://sysdig.com/blog/golang-expvar-custom-metrics/)
22. [Go App Monitoring: expvar, Prometheus and StatsD](https://www.opsdash.com/blog/golang-app-monitoring-statsd-expvar-prometheus.html)
23. [Expose application metrics with expvar](http://blog.ralch.com/tutorial/golang-metrics-with-expvar/)