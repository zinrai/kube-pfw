# kube-pfw: kubectl port forward wrapper

`kube-pfw` (Kubernetes Port Forward Wrapper) is a command-line interface (CLI) tool designed to simplify the process of port forwarding for Kubernetes services. It provides an interactive way to select services and ports, making it easier to use than the standard `kubectl port-forward` command, especially for services with multiple ports.

## Features

- Lists all services in a specified namespace
- Displays available ports for each service
- Allows interactive selection of services and ports
- Supports services with multiple ports
- Executes `kubectl port-forward` command with selected options

## Installation

Build the tool:

```
$ go build
```

## Usage

Run the tool by specifying the namespace:

```
./kube-pfw <namespace>
```

For example:

```
./kube-pfw default
```

Follow the interactive prompts to:
1. Select a service
2. Choose a port (if the service has multiple ports)
3. Specify the local port to forward to

The tool will then execute the appropriate `kubectl port-forward` command.

## Examples

Let's look at the services in our cluster:

```
$ kubectl get svc
NAME             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                    AGE
kubernetes       ClusterIP   10.96.0.1       <none>        443/TCP                    2d
nginx-service1   ClusterIP   10.96.253.72    <none>        80/TCP,8000/TCP,8001/TCP   2m53s
nginx-service2   ClusterIP   10.96.234.132   <none>        8010/TCP                   2m53s
```

### Service with Multiple Ports

Here's an example of using `kube-pfw` with a service that has multiple ports:

```
$ ./kube-pfw default
* service:
  1. kubernetes ( port 443 )
  2. nginx-service1 ( port 80 , 8000 , 8001 )
  3. nginx-service2 ( port 8010 )
Enter the number: 2
* nginx-service1:
  1. 80
  2. 8000
  3. 8001
Enter the number: 2
* Local Port: 8000
Exec Command: /usr/local/bin/kubectl port-forward service/nginx-service1 8000:8000 -n default
Forwarding from 127.0.0.1:8000 -> 80
Forwarding from [::1]:8000 -> 80
Handling connection for 8000
^C
```

In this example, we selected `nginx-service1`, which has multiple ports. We then chose port 8000 and forwarded it to local port 8000.

### Service with a Single Port

Here's an example of using `kube-pfw` with a service that has a single port:

```
$ ./kube-pfw default
* service:
  1. kubernetes ( port 443 )
  2. nginx-service1 ( port 80 , 8000 , 8001 )
  3. nginx-service2 ( port 8010 )
Enter the number: 3
* Local Port: 8010
Exec Command: /usr/local/bin/kubectl port-forward service/nginx-service2 8010:8010 -n default
Forwarding from 127.0.0.1:8010 -> 80
Forwarding from [::1]:8010 -> 80
^C
```

In this example, we selected `nginx-service2`, which has only one port (8010). The tool automatically selected this port and prompted us for the local port to use.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
