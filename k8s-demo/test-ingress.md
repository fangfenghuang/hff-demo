apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: nginx-example
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx-example
  rules:
  - http:
      paths:
      - path: /api/kubeapiproxy
        backend:
          serviceName: hfftest
          servicePort: 9999