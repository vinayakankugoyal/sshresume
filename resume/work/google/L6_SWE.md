# Staff Software Engineer

## **Google** | Kubernetes Security Team | *Nov 2024 - Present*

As the technical lead I architected the roadmap and provided leadership for a 7 member software and security engineering team focused on hardening Google Kubernetes Engine (GKE) and OSS Kubernetes.

Proposed, designed and led a team to build a token broker service for GKE system workloads to exchange pod service account tokens for down-scoped IAM service account tokens for metrics and logging. This allowed us to decouple GKE system workload permissions from customer workload permissions and the GKE node service account. Giving us the flexibility to add new capabilities to system workloads without requiring any actions from customers.

Proposed, designed and implemented [KEP-4633](https://kep.k8s.io/4633) in OSS Kubernetes. This KEP allowed admins to explicitly configure which endpoints support anonymous authentication. This prevents cluster exploitation even if a user mistakenly created RBAC bindings for system:anonymous. To raise awareness about this security misconfiguration and about this KEP I gave a [talk](https://www.youtube.com/watch?v=PbZbojx4kVM) at KubeCon NA.

Designed and implemented [KEP-2862](https://kep.k8s.io/2862) in OSS Kubernetes. This KEP breaks the often misused nodes/proxy permissions into more fine-grained permissions. Allowing users to write RBAC policies that only grant the permissions required by their workloads, in a way that was backwards compatible. Before this KEP permissions to access the /healthz endpoint and /exec/ required the same nodes/proxy subresource permission.

Proposed, designed and implemented [KEP-5040](https://kep.k8s.io/5040) in OSS Kubernetes. This KEP removed support for the built-in gitRepo volume driver which had been unmaintained. Recently vulnerabilities were found which meant Kubernetes and Cloud Providers would spend precious resources patching it. I worked with SIG-STORAGE and SIG-API to drive consensus in removing this volume driver from Kubernetes.
