# Senior Software Engineer

## **Google** | Kubernetes Security Team | *Nov 2021 - Nov 2024*

Led a team of SWEs and SEs defining the roadmap and OKRs and identifying the top x risks and prioritizing work that mitigated the risk.

Spearheaded the design, implementation and deployment of the Kubernetes Security Validation service, a proprietary policy engine for scanning misconfigurations in Kubernetes workloads. This system proactively detects thousands of Kubernetes security misconfigurations across GKE and high-priority Google 1Ps, and prevents such workloads from being submitted; producing an inventory of the existing risk. Led a multi-year cross functional program to remediate the findings from the service and reduced them from thousands to a few well documented exemptions.

io_uring was the largest contributor to GKE kernel vulnerabilities and Google spent upwards of $1.8M in bug bounties on it. I worked with containerd maintainers and drove consensus to block io_using system calls in the default seccomp profile, a change that docker also made following our lead.

Designed an allowlist system for popular partner workloads to be onboarded to GKE Autopilot clusters. GKE Autopilot clusters apply restrictions on workloads that prevent popular logging, monitoring and security scanning tools like jFrog, CrowdStrike, Wiz etc. from running on these clusters. I built an allowlisting scheme that matched the incoming workload to the allowlist and allowed it to bypass certain restrictions when the match succeeded.
