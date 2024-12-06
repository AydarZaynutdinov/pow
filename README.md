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

## How to run

### Docker

You can easily run the service and client with a single `docker-compose` command:

```docker
docker-compose up --build
```

or via Makefile command

```makefile
make up
```

This command starts the [docker-compose file](./docker-compose.yml) that includes `server` and `client`.

### Locally

If you want to run it locally, you need to:

- Run the [server](./cmd/server/main.go)
- Run the [client](./cmd/server/main.go) (optional)

For the server, you can specify the config file path using the `SERVER_CONFIG_FILE` env, or it will use the default
values.

## How to test

There is a Makefile command `test` to run unit tests