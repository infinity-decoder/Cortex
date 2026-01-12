<a name="readme-top"></a>

# 👨‍💻Author: INFINITY DECODER

<p align="center">
  <a href="https://github.com/infinity-decoder/Cortex">
    <img src="https://raw.githubusercontent.com/infinity-decoder/Cortex/master/frontend/public/cortex_SaaS.png" alt="Cortex Logo" width="120" height="120">
  </a>
  <h1 align="center">Cortex</h1>
  <p align="center">
    <strong>Container-Aware Attack Surface Monitoring SaaS for SMBs</strong>
    <br />
    <a href="https://github.com/infinity-decoder/Cortex/issues">Report Bug</a>
    ·
    <a href="https://github.com/infinity-decoder/Cortex/issues">Request Feature</a>
    <br />
    <a href="https://infinitydecoder.net"><strong>Visit Website »</strong></a>
    <br />
    <br />
    <a href="#features">View Features</a>
    ·
    <a href="#getting-started">Getting Started</a>
    ·
    <a href="#tech-stack">Tech Stack</a>
  </p>

---

## 🚀 Mission

Build a continuous, external, container-focused attack surface monitoring platform that detects exposed containers, Kubernetes services, Docker APIs, misconfigured ingress controllers, and container-related risks for SMBs before attackers do.

Cortex is designed to be **External (no agents)**, **Continuous (scheduled scans)**, and **Container-aware**, providing actionable remediation steps for high-risk exposures.

## ✨ Core Features

### 1. Asset Discovery Engine
- **DNS Enumeration**: Advanced subdomain discovery using intelligent wordlists.
- **TLS Certificate Analysis**: Extracting Subject Alternative Names (SANs) to uncover hidden assets.
- **Service Mapping**: Identifying public-facing infrastructure across root domains.

### 2. Container & Orchestration Detection
- **Docker Remote API**: Detection of exposed :2375/:2376 ports.
- **Kubernetes API & Kubelet**: Monitoring for exposed :6443 and :10250 endpoints.
- **Container Registries**: Identifying unauthenticated registries (/v2/_catalog).
- **Ingress Controllers**: Fingerprinting NGINX, Traefik, Istio, and more.

### 3. Exposure Classification Engine
- **Rules-Based Risk Scoring**: Classifies findings into Low, Medium, High, and Critical.
- **Threat Mapping**: Maps exposures to potential attack vectors like Data Exposure and RCE.
- **Safe Checks**: All detection is passive or read-only—no exploitation or brute forcing.

### 4. Continuous Monitoring
- **Delta Detection**: Identify when new services become public or risk levels change.
- **Scan History**: Maintain a comprehensive timeline of your attack surface evolution.

## 🛠 Tech Stack

- **Backend**: [Go](https://go.dev/) (Goroutine-based concurrency, Modular architecture)
- **Frontend**: [Next.js](https://nextjs.org/) (React, Tailwind CSS, Dashboard-focused UI)
- **Database**: [PostgreSQL](https://www.postgresql.org/) (Multi-tenant schema)
- **Deployment**: Optimized for single-VPS production with horizontal scalability potential.

## 🏁 Getting Started

### Prerequisites
- Go 1.23+
- Node.js 22+
- PostgreSQL

### Installation

1. **Backend Setup**
   ```bash
   cd backend
   go mod download
   go build -o bin/api.exe ./cmd/api
   ```

2. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## ⚖️ Legal & Safety
Cortex is built for ethical security monitoring.
- Only scan user-verified assets (DNS TXT verification required).
- Rate-limited scanning to respect network etiquette.
- Absolutely no exploitation, brute forcing, or credential testing.

---

<p align="right">(<a href="#readme-top">back to top</a>)</p>
