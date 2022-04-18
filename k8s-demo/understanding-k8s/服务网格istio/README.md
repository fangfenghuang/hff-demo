# 服务网格Service Mesh
服务网格(Service Mesh)用来描述组成这些应用程序的微服务网络以及它们之间的交互。

# istio
Istio是最初由IBM，Google和Lyft开发的服务网格的开源实现。它可以透明地分层到分布式应用程序上，并提供服务网格的所有优点，例如流量管理，安全性和可观察性。

Istio的工作原理是以Sidcar的形式将Envoy的扩展版本作为代理布署到每个微服务中


熔断: client 太热情, server 扛不住, 需要对超出能力外的请求进行熔断, 或者说限流
重试: client 没有收到 server 的正常响应, 需要进行重试
链路监控: client / server 很多, 调用链复杂, 如何有效追踪和监控


