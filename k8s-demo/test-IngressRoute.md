apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  annotations:
    kubernetes.io/ingress.class: traefik
  name: ingressroute-example
spec:
  routes:
    - kind: Rule
      match: PathPrefix(`/api/kubeapiproxy`)
      services:
        - kind: Service
          name: hfftest
          namespace: default
          port: 9999
