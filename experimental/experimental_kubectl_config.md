---
sidebarTitle: experimental kubectl config
sidebarOrder: 1
---

## experimental kubectl config

Modify kubeconfig files

### Synopsis

Modify kubeconfig files using subcommands like "kubectl config set current-context my-context".

 The loading order follows these rules:

  1.  If the --kubeconfig flag is set, then only that file is loaded. The flag may only be set once and no merging takes place.
  2.  If $KUBECONFIG environment variable is set, then it is used as a list of paths (normal path delimiting rules for your system). These paths are merged. When a value is modified, it is modified in the file that defines the stanza. When a value is created, it is created in the first file that exists. If no files in the chain exist, then it creates the last file in the list.
  3.  Otherwise, ${HOME}/.kube/config is used and no merging takes place.

```
qcadmin experimental kubectl config SUBCOMMAND
```

### Options

```
  -h, --help                help for config
      --kubeconfig string   use a particular kubeconfig file
```

### Options inherited from parent commands

```
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/home/runner/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --config string                  The qcadmin config file to use
      --context string                 The name of the kubeconfig context to use
      --debug                          Prints the stack trace if an error occurs
      --disable-compression            If true, opt-out of response compression for all requests to the server
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --match-server-version           Require server version to match client version
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --profile string                 Name of profile to capture. One of (none|cpu|heap|goroutine|threadcreate|block|mutex) (default "none")
      --profile-output string          Name of the file to write the profile to (default "profile.pprof")
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --silent                         Run in silent mode and prevents any qcadmin log output except panics & fatals
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
      --warnings-as-errors             Treat warnings received from the server as errors and exit with a non-zero exit code
```

### SEE ALSO

* [qcadmin experimental kubectl](experimental_kubectl.md)	 - Kubectl controls the Kubernetes cluster manager
* [experimental kubectl config current-context](experimental_kubectl_config_current-context.md)	 - Display the current-context
* [experimental kubectl config delete-cluster](experimental_kubectl_config_delete-cluster.md)	 - Delete the specified cluster from the kubeconfig
* [experimental kubectl config delete-context](experimental_kubectl_config_delete-context.md)	 - Delete the specified context from the kubeconfig
* [experimental kubectl config delete-user](experimental_kubectl_config_delete-user.md)	 - Delete the specified user from the kubeconfig
* [experimental kubectl config get-clusters](experimental_kubectl_config_get-clusters.md)	 - Display clusters defined in the kubeconfig
* [experimental kubectl config get-contexts](experimental_kubectl_config_get-contexts.md)	 - Describe one or many contexts
* [experimental kubectl config get-users](experimental_kubectl_config_get-users.md)	 - Display users defined in the kubeconfig
* [experimental kubectl config rename-context](experimental_kubectl_config_rename-context.md)	 - Rename a context from the kubeconfig file
* [experimental kubectl config set](experimental_kubectl_config_set.md)	 - Set an individual value in a kubeconfig file
* [experimental kubectl config set-cluster](experimental_kubectl_config_set-cluster.md)	 - Set a cluster entry in kubeconfig
* [experimental kubectl config set-context](experimental_kubectl_config_set-context.md)	 - Set a context entry in kubeconfig
* [experimental kubectl config set-credentials](experimental_kubectl_config_set-credentials.md)	 - Set a user entry in kubeconfig
* [experimental kubectl config unset](experimental_kubectl_config_unset.md)	 - Unset an individual value in a kubeconfig file
* [experimental kubectl config use-context](experimental_kubectl_config_use-context.md)	 - Set the current-context in a kubeconfig file
* [experimental kubectl config view](experimental_kubectl_config_view.md)	 - Display merged kubeconfig settings or a specified kubeconfig file

::: tip
>Auto generated by spf13/cobra on 2025-Jan-17
:::