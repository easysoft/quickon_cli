## qcadmin cluster join

join cluster

```
qcadmin cluster join [flags]
```

### Options

```
  -h, --help                 help for join
      --master stringArray   master ip list, e.g: 192.168.0.1:22
      --password string      ssh password
      --pkfile string        ssh private key, if not set, will use password
      --pkpass string        ssh private key password
  -u, --username string      ssh user (default "root")
      --worker stringArray   worker ip list, e.g: 192.168.0.1:22
```

### Options inherited from parent commands

```
      --config string   The qcadmin config file to use
      --debug           Prints the stack trace if an error occurs
      --silent          Run in silent mode and prevents any qcadmin log output except panics & fatals
```

### SEE ALSO

* [qcadmin cluster](qcadmin_cluster.md)	 - Cluster commands

###### Auto generated by spf13/cobra on 6-Nov-2023