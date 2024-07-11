# Broker

This broker microservice is the singular API gateway for [@YanSystems](https://github.com/YanSystems). It proxies all requests from the client to the corresponding microservice.

## Running locally

This microservice leverages Kubernetes' DNS to proxy requests to the corresponding microservices by their service names. Therefore, you can only test and/or work with it in a kubernetes cluster where the other services are defined. To learn how to do this, please read the `README.md` in the [deployment repository](https://github.com/YanSystems/k8s).

## License

This broker microservice is [MIT licensed.](https://github.com/YanSystems/broker/blob/main/LICENSE)