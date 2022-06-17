https://kyverno.io/








kubectl get clusterpolicy

```yaml
[root@RqKubeDev03 kbuser]# kubectl get clusterpolicy add-default-resources -o yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: add-default-resources
spec:
  background: false
  failurePolicy: Fail
  rules:
  - exclude:
      resources: {}
    generate:
      clone: {}
    match:
      any:
      - resources:
          kinds:
          - Pod
      resources: {}
    mutate:
      patchStrategicMerge:
        spec:
          containers:
          - (name): '*'
            resources:
              limits:
                +(cpu): 1000m
                +(memory): 2Gi
              requests:
                +(cpu): 1000m
                +(memory): 2Gi
    name: add-default-requests
    preconditions:
      all:
      - key: '{{request.operation}}'
        operator: In
        value:
        - CREATE
        - UPDATE
      - key: '{{request.object.spec.runtimeClassName}}'
        operator: Equals
        value: kata
    validate: {}

```