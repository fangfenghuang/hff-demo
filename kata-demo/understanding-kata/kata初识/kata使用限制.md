https://github.com/kata-containers/kata-containers/blob/main/docs/Limitations.md

[TOC]
kata在硬件隔离的虚拟机中运行容器，每个vm都有独立的内核，由于这种高程度的隔离，某些容器功能无法启用，或者通过vm隐式启用。

OCI 规范定义了运行时必须支持的最低规范，以便与 Docker 等容器管理器进行互操作。如果运行时不支持 OCI 规范的某些方面，则根据定义它是一个限制。

但是，runc并不完全符合 OCI 规范本身。

以下是社区列出的限制
https://github.com/pulls?q=label%3Alimitation+org%3Akata-containers+is%3Aopen


# 一些列出的限制
## CLI命令：不支持docker和Podman
  不支持docker --runtime指定kata运行时

## runtime命令：不支持checkpoint、restore、events 、update 

## 网络？未列出

## 资源管理
对于基于 VM 的系统，将 cgroup、CPU、内存和存储等资源约束应用于工作负载并不总是那么简单。
### docker 的--cpu
### vcpus

## 架构限制
### 存储限制subPaths

### 主机资源共享:securityContext privileged
https://github.com/kata-containers/kata-containers/blob/main/docs/how-to/privileged.md


# kata容器限制级别
## 通过guest kernel参数设置限制kata vm，如sysctl参数等

## 限制kata container

## 通过主机级约束限制hypervisor；设置hypervisor参数


