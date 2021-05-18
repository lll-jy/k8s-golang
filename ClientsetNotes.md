# K8s Go client API

## Connect to k8s

### Imports

```go
"k8s.io/client-go/kubernetes"
"k8s.io/client-go/tools/clientcmd"
```

### Connection function

```go
func connectToK8s() *kubernetes.Clientset {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	configPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		log.Panicln("failed to create K8s config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("Failed to create K8s clientset")
	}

	return clientset
}
```

## Clientset content API

* Format: `clientset.<V1_field>().<content>(args)`. For example, the first row in the following table means `clientset.
AppsV1().ControllerRevisions(namespace string)`. If "has namespace" column is false, the arguments is empty.
* Only important ones are noted down here.

| V1 field           | Content                       | has namespace |
|--------------------|-------------------------------|---------------|
| `AppsV1`           | `ControllerRevisions`         | true          | 
| `AppsV1`           | `DaemonSets`                  | true          |
| `AppsV1`           | `Deployments`                 | true          |
| `AppsV1`           | `ReplicaSets`                 | true          |
| `AppsV1`           | `StatefulSets`                | true          |
| `AuthenticationV1` | `TokenReviews`                | false         |
| `AuthorizationV1`  | `LocalSubjectAccessReviews`   | true          |
| `AuthorizationV1`  | `SelfSubjectAccessReviews`    | false         |
| `AuthorizationV1`  | `SelfSubjectRulesReviews`     | false         |
| `AuthorizationV1`  | `SubjectAccessReviews`        | false         |
| `AutoscalingV1`    | `HorizontalPodAutoscalers`    | true          |
| `BatchV1`          | `Jobs`                        | true          |
| `CertificatesV1`   | `CertificatesSigningRequests` | false         |
| `CoordinationV1`   | `Leases`                      | true          |
| `CoreV1`           | `ConfigMaps`                  | true          |
| `CoreV1`           | `ComponentStatuses`           | false         |
| `CoreV1`           | `Endpoints`                   | true          |
| `CoreV1`           | `Events`                      | true          |
| `CoreV1`           | `LimitRanges`                 | true          |
| `CoreV1`           | `Namespaces`                  | false         |
| `CoreV1`           | `Nodes`                       | false         |
| `CoreV1`           | `PersistentVolumes`           | false         |
| `CoreV1`           | `PersistentVolumesClaims`     | true          |
| `CoreV1`           | `Pods`                        | true          |
| `CoreV1`           | `PodTemplates`                | true          |
| `CoreV1`           | `ReplicationControllers`      | true          |
| `CoreV1`           | `ResourceQuotas`              | true          |
| `CoreV1`           | `Secrets`                     | true          |
| `CoreV1`           | `ServiceAccounts`             | true          |
| `CoreV1`           | `Services`                    | true          |
| `DiscoveryV1`      | `EndpointSlices`              | true          |
| `EventsV1`         | `Events`                      | true          |
| `NetworkingV1`     | `IngressClasses`              | false         |
| `NetworkingV1`     | `Ingresses`                   | true          |
| `NetworkingV1`     | `NetworkPolicies`             | true          |
| `NodeV1`           | `RuntimeClasses`              | false         |
| `PolicyV1`         | `PodDisruptionBudgets`        | true          |
| `RbacV1`           | `ClusterRoleBindings`         | false         |
| `RbacV1`           | `ClusterRoles`                | false         |
| `RbacV1`           | `RoleBindings`                | true          |
| `RbacV1`           | `Roles`                       | true          |
| `SchedulingV1`     | `PriorityClasses`             | false         |
| `StorageV1`        | `CSIDrivers`                  | false         |
| `StorageV1`        | `CSINodes`                    | false         |
| `StorageV1`        | `StorageClasses`              | false         |
| `StorageV1`        | `VolumeAttachments`           | false         |
