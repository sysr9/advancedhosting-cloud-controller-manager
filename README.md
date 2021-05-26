# Kubernetes Cloud Controller Manager for Advanced Hosting

## How to install

1. Every kubelet must run with `--cloud-provider=external`. This is to ensure that the kubelet is aware that it must be initialized by the cloud controller manager before it is scheduled any work. 
2. Set the following environment variables to the cluster:
- `AH_API_TOKEN` - Client access token.
- `AH_API_URL` (Optional) - Base api url. `"https://api.websa.com"` by default.
- `AH_CLUSTER_PRIVATE_NETWORK_ID` - ID of the private network to which the cluster is assigned.
- `AH_CLUSTER_DATACENTER_ID` - ID of the datacenter to which the cluster is assigned.

3. Build Cloud controller manager image:
```
docker build . --build-arg TAG=<IMAGE_TAG>
```

4. Deploy CCM. Manifest **example**:
```
---
  ---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: advancedhosting-cloud-controller-manager
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: advancedhosting-cloud-controller-manager
  template:
    metadata:
      labels:
        app: advancedhosting-cloud-controller-manager
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      dnsPolicy: Default
      hostNetwork: true
      serviceAccountName: cloud-controller-manager
      tolerations:
        - key: "node.cloudprovider.kubernetes.io/uninitialized"
          value: "true"
          effect: "NoSchedule"
        - key: "CriticalAddonsOnly"
          operator: "Exists"
        - key: "node-role.kubernetes.io/master"
          effect: NoSchedule
      containers:
      - image: <LINK_TO_REGISTRY>
        name: advancedhosting-cloud-controller-manager
        args:
          - --cloud-provider=advancedhosting
          - --leader-elect=true
          - --allow-untagged-cloud
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-controller-manager
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  name: system:cloud-controller-manager
rules:

- apiGroups:
  - coordination.k8s.io
  resources:
  - leases
  verbs:
  - get
  - watch
  - list
  - create
  - update
  - delete
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
  - update
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - '*'
- apiGroups:
  - ""
  resources:
  - nodes/status
  verbs:
  - patch
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - services/status
  verbs:
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - serviceaccounts
  verbs:
  - create
- apiGroups:
  - ""
  resources:
  - persistentvolumes
  verbs:
  - get
  - list
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - endpoints
  verbs:
  - create
  - get
  - list
  - watch
  - update
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:cloud-controller-manager
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:cloud-controller-manager
subjects:
- kind: ServiceAccount
  name: cloud-controller-manager
  namespace: kube-system
```