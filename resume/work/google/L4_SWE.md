# Software Engineer

## **Google** | Kubernetes Security Team | *Sep 2019 - Nov 2021*

Proposed the migration of GKE workloads to rootless. I owned this workstream end-to-end working with teams across the GKE org and migrated all GKE system workloads to rootless with a few well documented exemptions. This vastly reduced the risk of container escapes in GKE and its customers. I proposed [KEP-2568](https://kep.k8s.io/2568) and drove consensus in Kubernetes to support running the Kubernetes control-plane components as non-root in kubeadm, which is the canonical tool used to deploy Kubernetes clusters by millions. I also gave a [talk](https://www.youtube.com/watch?v=uouH9fsWVIE) at KubeCon about the importance and techniques of how to migrate container workloads to rootless.

Designed and implemented secure by default controls for GKE Autopilot clusters that automatically applied security controls on customer workloads like dropping CAP_NET_RAW and applying seccomp filters to disallow dangerous syscalls.

Designed and implemented CMEK for boot disk and etcd secrets encryption for GKE Autopilot clusters. This allowed customers to encrypt the node bootdisk and Kubernetes secrets stored in etcd with a key that was managed by customers putting them in control of their data.