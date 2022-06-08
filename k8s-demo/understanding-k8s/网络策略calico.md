```yaml
---
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: global-default
spec:
  egress:
  - action: Allow
    destination: {}
    source: {}
  ingress:
  - action: Allow
    destination: {}
    source: {}
  order: 9999
  types:
  - Ingress
  - Egress
---
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: global-egress-deny
spec:
  egress:
  - action: Deny
    destination:
      selector: area == 'dummy'
    source: {}
  namespaceSelector: (! has(privileged-namespace))
  order: 502
  types:
  - Egress
---
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: global-egress-allow
spec:
  egress:
  - action: Allow
    destination:
      nets:
      - 127.0.0.1/32
    source:
      namespaceSelector: name == "dummy"
  order: 501
  types:
  - Egress
---
apiVersion: projectcalico.org/v3
kind: GlobalNetworkPolicy
metadata:
  name: global-privileged-namespace
spec:
  egress:
  - action: Allow
    destination: {}
    source: {}
  ingress:
  - action: Allow
    destination: {}
    source: {}
  namespaceSelector: has(privileged-namespace)
  order: 500
  types:
  - Ingress
  - Egress
```