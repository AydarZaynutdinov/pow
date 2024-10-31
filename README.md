# Proof of Work service

## Overview

[Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work) is a form of cryptographic proof in which one party (the
prover) proves to others (the verifiers) that a certain amount of a specific computational effort has been expended.
Using this mechanism, the service could be sure that the client has done a certain amount of work/computation, which
authorizes the service to send the required data to the client.

### Protocol

This PoW scheme supports various protocols, including the challenge-response protocol, used in this service.

### Algorithm

The PoW scheme uses leading-zero hash algorithm - hash calculations where a client attempts to find a solution that,
when combined with the challenge string, produces a hash with a specific number of leading zeros (e.g., SHA-256(
challenge + solution)). This difficulty can be adjusted by increasing the required number of leading zeros in the hash,
making it more difficult for the client to solve while still remaining easy for the server to verify.

#### Pros of chosen algorithm:

- Security: Hash-based algorithms (like SHA-256) are secure and widely recognized for cryptographic use. By basing the
  proof on leading zeroes, the algorithm remains simple to verify but computationally demanding to solve, especially as
  difficulty increases.
- Adjustable Difficulty: The difficulty can be increased or decreased simply by adjusting the number of leading zeroes
  required in the hash. This flexibility is useful for adapting the PoW to different load and security requirements,
  making it more scalable than algorithms with fixed difficulty.
- Efficient Verification: Hash-based solutions are easy for the server to verify since calculating the hash for one
  attempt is minimal computational work. The client bears the primary burden of proof, making it ideal for a scenario
  like DoS protection where server resources need to be conserved.

There are a lot of alternative [algorithms](https://en.wikipedia.org/wiki/Proof_of_work#List_of_proof-of-work_functions)
but chosen algorithm effectively balances ease of implementation, scalability, security, and efficiency, making it a
strong choice for this PoW system.

### Service

This service has 2 HTTP endpoints for handling Proof of Work.
First endpoint returns a `challenge` and `difficulty`.
Using second one client should send a `solution` that meets the condition

When a client requests a `challenge`, it is stored in the cache (using Redis) with a specified TTL. This allows the
server to validate the `solution` against a stored `challenge` and ensures the client’s response is for a legitimate,
server-generated `challenge`.

After sending correct solution client receives one of random quotes (quotes is placed in the config file)

Redis is used to be able to run this service as several pod app with load-balancing mechanism between server and
clients and allows pods work independently.

### Client

The project also includes a simple client that sends requests to the service every 30 seconds.

## API

API for the service could be found in the [Openapi file](/api/pow-openapi.yaml).

## How to run

### Docker

You can easily run the service and client with a single `docker-compose` command:

```docker
docker-compose up --build
```

This command starts the [docker-compose file](./docker-compose.yml) which includes `redis`, `server` and `client`.

### Locally

If you want to run it locally, you need to:

- Install and run [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/install-stack/docker/)
- Run the [server](./cmd/pow-server/main.go)
- Run the [client](./cmd/pow-server/main.go) (optional)

For the server, you can specify the config file path using the `-config` flag, or it will use the
default [config file](/config/config.dev.yaml).

## How to test

There is a [Makefile](./Makefile) that contains `all-test-one-command`. This command
- runs docker-compose file (to run Redis)
- build and run service
- runs unit and integration tests

## Maintenance

The service collects several metrics during its operation. These metrics are
accessible [there](http://127.0.0.1:8081/metrics):

| Name                 | Type        | Labels             | Description                    | 
|----------------------|-------------|--------------------|--------------------------------|
| http_requests_total  | `counter`   | `path`, `code`     | Total number of HTTP requests  |
| http_response_times  | `histogram` | `path`, `code`     | Response times in ms           |
| cache_requests_total | `counter`   | `method`, `status` | Total number of cache requests |
| cache_response_times | `histogram` | `method`, `status` | Response times in ms           |
| errors_total         | `counter`   | `level`            | Total number of errors         |

## Structure

```
|-- pow
   |-- api                      // OpenAPI documentation
   |-- cmd
      |-- pow-client            // client entrypoint
      |-- pow-server            // service entrypoint
   |-- config                   // configuration files
   |-- internal
      |-- cache                 // cache logic
      |-- config                // configuration entities
      |-- controller
         |-- api                // main service handlers
         |-- maintenance        // maintenance service handlers
      |-- infrastructure
         |-- server             // code to run the HTTP server
      |-- logger                // logging
      |-- metrics               // service metrics
      |-- service
         |-- challenge          // main service logic
```