
# 自动清理
自动清理已完成 Job （状态为 Complete 或 Failed）的另一种方式是使用由 TTL 控制器所提供 的 TTL 机制。 通过设置 Job 的 .spec.ttlSecondsAfterFinished 字段，可以让该控制器清理掉 已结束的资源。

TTL 控制器清理 Job 时，会级联式地删除 Job 对象。 换言之，它会删除所有依赖的对象，包括 Pod 及 Job 本身。 注意，当 Job 被删除时，系统会考虑其生命周期保障，例如其 Finalizers。
