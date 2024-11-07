---
sidebarTitle: experimental kubectl
sidebarOrder: 1
---

## experimental kubectl

Kubectl controls the Kubernetes cluster manager

### Synopsis

kubectl controls the Kubernetes cluster manager.

 Find more information at: https://kubernetes.io/docs/reference/kubectl/

```
qcadmin experimental kubectl [flags]
```

### Options

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/runner/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -h, --help                           help for kubectl
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --match-server-version           Require server version to match client version
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --profile string                 Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex) (default "none")
      --profile-output string          Name of the file to write the profile to (default "profile.pprof")
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
      --warnings-as-errors             Treat warnings received from the server as errors and exit with a non-zero exit code
```

### Options inherited from parent commands

```
      --config string   The qcadmin config file to use
      --debug           Prints the stack trace if an error occurs
      --silent          Run in silent mode and prevents any qcadmin log output except panics & fatals
```

### SEE ALSO

* [qcadmin experimental](experimental.md)	 - Experimental commands that may be modified or deprecated
* [experimental kubectl annotate](experimental_kubectl_annotate.md)	 - Update the annotations on a resource
* [experimental kubectl api-resources](experimental_kubectl_api-resources.md)	 - Print the supported API resources on the server
* [experimental kubectl api-versions](experimental_kubectl_api-versions.md)	 - Print the supported API versions on the server, in the form of "group/version"
* [experimental kubectl apply](experimental_kubectl_apply.md)	 - Apply a configuration to a resource by file name or stdin
* [experimental kubectl attach](experimental_kubectl_attach.md)	 - Attach to a running container
* [experimental kubectl auth](experimental_kubectl_auth.md)	 - Inspect authorization
* [experimental kubectl autoscale](experimental_kubectl_autoscale.md)	 - Auto-scale a deployment, replica set, stateful set, or replication controller
* [experimental kubectl certificate](experimental_kubectl_certificate.md)	 - Modify certificate resources
* [experimental kubectl cluster-info](experimental_kubectl_cluster-info.md)	 - Display cluster information
* [experimental kubectl completion](experimental_kubectl_completion.md)	 - Output shell completion code for the specified shell (bash, zsh, fish, or powershell)
* [experimental kubectl config](experimental_kubectl_config.md)	 - Modify kubeconfig files
* [experimental kubectl cordon](experimental_kubectl_cordon.md)	 - Mark node as unschedulable
* [experimental kubectl cp](experimental_kubectl_cp.md)	 - Copy files and directories to and from containers
* [experimental kubectl create](experimental_kubectl_create.md)	 - Create a resource from a file or from stdin
* [experimental kubectl debug](experimental_kubectl_debug.md)	 - Create debugging sessions for troubleshooting workloads and nodes
* [experimental kubectl delete](experimental_kubectl_delete.md)	 - Delete resources by file names, stdin, resources and names, or by resources and label selector
* [experimental kubectl describe](experimental_kubectl_describe.md)	 - Show details of a specific resource or group of resources
* [experimental kubectl diff](experimental_kubectl_diff.md)	 - Diff the live version against a would-be applied version
* [experimental kubectl drain](experimental_kubectl_drain.md)	 - Drain node in preparation for maintenance
* [experimental kubectl edit](experimental_kubectl_edit.md)	 - Edit a resource on the server
* [experimental kubectl events](experimental_kubectl_events.md)	 - List events
* [experimental kubectl exec](experimental_kubectl_exec.md)	 - Execute a command in a container
* [experimental kubectl explain](experimental_kubectl_explain.md)	 - Get documentation for a resource
* [experimental kubectl expose](experimental_kubectl_expose.md)	 - Take a replication controller, service, deployment or pod and expose it as a new Kubernetes service
* [experimental kubectl get](experimental_kubectl_get.md)	 - Display one or many resources
* [experimental kubectl kustomize](experimental_kubectl_kustomize.md)	 - Build a kustomization target from a directory or URL
* [experimental kubectl label](experimental_kubectl_label.md)	 - Update the labels on a resource
* [experimental kubectl logs](experimental_kubectl_logs.md)	 - Print the logs for a container in a pod
* [experimental kubectl options](experimental_kubectl_options.md)	 - Print the list of flags inherited by all commands
* [experimental kubectl patch](experimental_kubectl_patch.md)	 - Update fields of a resource
* [experimental kubectl plugin](experimental_kubectl_plugin.md)	 - Provides utilities for interacting with plugins
* [experimental kubectl port-forward](experimental_kubectl_port-forward.md)	 - Forward one or more local ports to a pod
* [experimental kubectl proxy](experimental_kubectl_proxy.md)	 - Run a proxy to the Kubernetes API server
* [experimental kubectl replace](experimental_kubectl_replace.md)	 - Replace a resource by file name or stdin
* [experimental kubectl rollout](experimental_kubectl_rollout.md)	 - Manage the rollout of a resource
* [experimental kubectl run](experimental_kubectl_run.md)	 - Run a particular image on the cluster
* [experimental kubectl scale](experimental_kubectl_scale.md)	 - Set a new size for a deployment, replica set, or replication controller
* [experimental kubectl set](experimental_kubectl_set.md)	 - Set specific features on objects
* [experimental kubectl taint](experimental_kubectl_taint.md)	 - Update the taints on one or more nodes
* [experimental kubectl top](experimental_kubectl_top.md)	 - Display resource (CPU/memory) usage
* [experimental kubectl uncordon](experimental_kubectl_uncordon.md)	 - Mark node as schedulable
* [experimental kubectl version](experimental_kubectl_version.md)	 - Print the client and server version information
* [experimental kubectl wait](experimental_kubectl_wait.md)	 - Experimental: Wait for a specific condition on one or many resources

::: tip
>Auto generated by spf13/cobra on 2024-Nov-7
:::