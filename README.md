# Label Watcher

This is primarily an example usage of the [Kubernetes Go client](https://github.com/kubernetes/client-go) and demo of a few language features for a tech talk.

There are two binaries. They both perform the same basic function: watch a kubernetes cluster for changes and output the results.

- `lwcli`: Outputs the results on the terminal as they are received.
- `lwserver`: Serves a log of the results over HTTPS.


## Building, Testing, and Running

Run `make` to build both binaries and generate self-signed certs for TLS. There are individual targets for testing (`make test`), running (`make cli-run` and `make server-run`), and creating certs (`make cert`).

Lwserver accepts a few optional flags. Run `./lwserver -h` after building to see them.


## Connecting to Kubernetes

The client is configured to run outside of your cluster. For this demo, we only support the `~/.kube/config` and use the current context. More sophisticated kubeconfig options aren't handled.
