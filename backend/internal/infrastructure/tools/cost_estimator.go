package tools

import (
	"fmt"
	"math"
	"strings"

	"github.com/Wendell-Dorta/SystemCrafter-Bot/backend/internal/domain/entity"
)

// EstimateCloudCostParams contains fine-grained parameters for real-world cloud cost calculation.
type EstimateCloudCostParams struct {
	Provider            string  `json:"provider"`            // AWS, GCP, AZURE, ORACLE (OCI)
	WorkloadTier        string  `json:"workloadTier"`        // Startup, Growth, Enterprise
	ComputeType         string  `json:"computeType"`         // serverless, containers (fargate/cloudrun/container apps), kubernetes (eks/gke/aks/oke), vms
	ComputeArchitecture string  `json:"computeArchitecture"` // arm (graviton/ampere) or x86
	DatabaseType        string  `json:"databaseType"`        // postgresql, mysql, mongodb, dynamo/cosmos/nosql
	DatabaseRedundancy  string  `json:"databaseRedundancy"`  // single-az, multi-az, multi-region
	CacheEnabled        bool    `json:"cacheEnabled"`        // redis/memcached
	StorageGB           int     `json:"storageGB"`           // GBs of object/block storage
	MonthlyReqsMil      float64 `json:"monthlyReqsMil"`      // Millions of monthly requests
	DataEgressGB        int     `json:"dataEgressGB"`        // Outbound internet traffic in GB
}

// EstimateCloudCosts calculates a realistic, benchmarked FinOps cost breakdown across AWS, GCP, Azure, and OCI.
func EstimateCloudCosts(params EstimateCloudCostParams) entity.CloudCostEstimate {
	provider := normalizeProvider(params.Provider)
	tier := normalizeTier(params.WorkloadTier)

	// Baseline calculations
	reqs := params.MonthlyReqsMil
	if reqs <= 0 {
		switch tier {
		case "Startup":
			reqs = 2.0
		case "Enterprise":
			reqs = 80.0
		default: // Growth
			reqs = 15.0
		}
	}

	storageGB := params.StorageGB
	if storageGB <= 0 {
		switch tier {
		case "Startup":
			storageGB = 100
		case "Enterprise":
			storageGB = 2500
		default:
			storageGB = 500
		}
	}

	egressGB := params.DataEgressGB
	if egressGB <= 0 {
		switch tier {
		case "Startup":
			egressGB = 150
		case "Enterprise":
			egressGB = 5000
		default:
			egressGB = 1000
		}
	}

	// Calculate primary provider cost
	breakdown, calcs, total := calculateProviderCosts(provider, tier, params, reqs, storageGB, egressGB)

	// Calculate comparative prices for all major clouds (AWS, GCP, Azure, Oracle OCI)
	comparative := make(map[string]float64)
	allProviders := []string{"AWS", "GCP", "AZURE", "ORACLE"}
	for _, p := range allProviders {
		_, _, pTotal := calculateProviderCosts(p, tier, params, reqs, storageGB, egressGB)
		comparative[p] = math.Round(pTotal*100) / 100
	}

	// Tailored FinOps optimization recommendations
	optimizations := generateOptimizations(provider, tier, params)

	return entity.CloudCostEstimate{
		Provider:           provider,
		WorkloadTier:       tier,
		MonthlyTotalUSD:    math.Round(total*100) / 100,
		BreakdownUSD:       breakdown,
		Calculations:       calcs,
		ComparativePricing: comparative,
		Optimizations:      optimizations,
	}
}

func normalizeProvider(p string) string {
	norm := strings.ToUpper(strings.TrimSpace(p))
	switch {
	case strings.Contains(norm, "ORACLE") || strings.Contains(norm, "OCI"):
		return "ORACLE"
	case strings.Contains(norm, "GCP") || strings.Contains(norm, "GOOGLE"):
		return "GCP"
	case strings.Contains(norm, "AZURE") || strings.Contains(norm, "MICROSOFT"):
		return "AZURE"
	default:
		return "AWS"
	}
}

func normalizeTier(t string) string {
	norm := strings.Title(strings.ToLower(strings.TrimSpace(t)))
	if norm != "Startup" && norm != "Enterprise" {
		return "Growth"
	}
	return norm
}

func calculateProviderCosts(
	provider, tier string,
	params EstimateCloudCostParams,
	reqs float64,
	storageGB, egressGB int,
) (map[string]float64, map[string]string, float64) {
	breakdown := make(map[string]float64)
	calcs := make(map[string]string)
	computeType := strings.ToLower(strings.TrimSpace(params.ComputeType))
	isARM := strings.Contains(strings.ToLower(params.ComputeArchitecture), "arm") ||
		strings.Contains(strings.ToLower(params.ComputeArchitecture), "graviton") ||
		strings.Contains(strings.ToLower(params.ComputeArchitecture), "ampere")

	// 1. COMPUTE CALCULATION
	var computeCost float64
	switch {
	case strings.Contains(computeType, "serverless") || strings.Contains(computeType, "lambda") ||
		strings.Contains(computeType, "cloudrun") || strings.Contains(computeType, "container app"):
		switch provider {
		case "ORACLE":
			computeCost = reqs * 2.10 // OCI Functions (Free tier discount)
			calcs["Compute"] = fmt.Sprintf("OCI Functions: %.1fM invocações (~128MB RAM/100ms)", reqs)
		case "GCP":
			computeCost = reqs * 2.80 // Cloud Run (Concurrency bundling)
			calcs["Compute"] = fmt.Sprintf("GCP Cloud Run: %.1fM requisições com auto-scale a zero", reqs)
		case "AZURE":
			computeCost = reqs * 3.10 // Azure Container Apps
			calcs["Compute"] = fmt.Sprintf("Azure Container Apps: %.1fM execuções gerenciadas", reqs)
		default: // AWS
			computeCost = reqs * 3.40 // Lambda + API Gateway HTTP
			calcs["Compute"] = fmt.Sprintf("AWS Lambda + HTTP API Gateway: %.1fM reqs/mês", reqs)
		}

	case strings.Contains(computeType, "k8s") || strings.Contains(computeType, "kubernetes") ||
		strings.Contains(computeType, "eks") || strings.Contains(computeType, "gke") ||
		strings.Contains(computeType, "aks") || strings.Contains(computeType, "oke"):
		switch provider {
		case "ORACLE":
			// OKE has no cluster control plane fee!
			if tier == "Enterprise" {
				computeCost = 380.0
				calcs["Compute"] = "OCI OKE (Cluster Gratuito) + 6x Nós Ampere A1 ARM Multi-AD"
			} else if tier == "Startup" {
				computeCost = 45.0
				calcs["Compute"] = "OCI OKE (Sem taxa de cluster) + 2x Nós E4 Flex"
			} else {
				computeCost = 160.0
				calcs["Compute"] = "OCI OKE (Sem taxa de cluster) + 3x Nós Ampere A1 (4 OCPU, 24GB)"
			}
		case "GCP":
			clusterFee := 73.0
			if tier == "Enterprise" {
				computeCost = clusterFee + 580.0
				calcs["Compute"] = "GKE Enterprise ($73 cluster fee) + 6x nós e2-standard-4"
			} else if tier == "Startup" {
				computeCost = 65.0 // Free zonal cluster fee
				calcs["Compute"] = "GKE Zonal (Tier Gratuito de Cluster) + 2x nós e2-medium"
			} else {
				computeCost = clusterFee + 220.0
				calcs["Compute"] = "GKE Standard ($73 fee) + 3x nós e2-standard-2 (Spot/On-Demand)"
			}
		case "AZURE":
			if tier == "Enterprise" {
				computeCost = 73.0 + 610.0
				calcs["Compute"] = "Azure AKS Standard + 6x nós D4s_v5 Multi-AZ"
			} else if tier == "Startup" {
				computeCost = 75.0
				calcs["Compute"] = "Azure AKS Free Tier + 2x nós B2ms Burstable"
			} else {
				computeCost = 73.0 + 240.0
				calcs["Compute"] = "Azure AKS Standard + 3x nós D2s_v5"
			}
		default: // AWS
			eksClusterFee := 73.0
			if tier == "Enterprise" {
				nodeCost := 590.0
				if isARM {
					nodeCost = 480.0
				}
				computeCost = eksClusterFee + nodeCost
				calcs["Compute"] = fmt.Sprintf("AWS EKS ($73 cluster) + 6x nós %s Multi-AZ", tern(isARM, "c7g.xlarge (ARM)", "m6i.xlarge (x86)"))
			} else if tier == "Startup" {
				computeCost = eksClusterFee + 70.0
				calcs["Compute"] = "AWS EKS ($73 cluster) + 2x nós t4g.medium (Graviton)"
			} else {
				nodeCost := 240.0
				if isARM {
					nodeCost = 195.0
				}
				computeCost = eksClusterFee + nodeCost
				calcs["Compute"] = fmt.Sprintf("AWS EKS ($73 cluster) + 3x nós %s", tern(isARM, "m7g.large (ARM)", "t3.xlarge"))
			}
		}

	default: // VMs (EC2, Compute Engine, Azure VMs, OCI Compute)
		switch provider {
		case "ORACLE":
			if tier == "Enterprise" {
				computeCost = 280.0
				calcs["Compute"] = "OCI Compute E4 Flex / Ampere A1 com Auto-Scaling"
			} else if tier == "Startup" {
				computeCost = 25.0
				calcs["Compute"] = "OCI Compute Always Free + 1x VM.Standard.E4.Flex"
			} else {
				computeCost = 95.0
				calcs["Compute"] = "OCI Compute: 2x instâncias E4 Flex balanceadas"
			}
		case "GCP":
			if tier == "Enterprise" {
				computeCost = 410.0
				calcs["Compute"] = "GCP Compute Engine (MIG Auto-Scaler com n2-standard-4)"
			} else if tier == "Startup" {
				computeCost = 48.0
				calcs["Compute"] = "GCP Compute Engine: 2x instâncias e2-small"
			} else {
				computeCost = 145.0
				calcs["Compute"] = "GCP Compute Engine: 2x instâncias e2-standard-2"
			}
		case "AZURE":
			if tier == "Enterprise" {
				computeCost = 430.0
				calcs["Compute"] = "Azure Virtual Machine Scale Sets (VMSS D4s_v5)"
			} else if tier == "Startup" {
				computeCost = 52.0
				calcs["Compute"] = "Azure VMs: 2x B2s Burstable"
			} else {
				computeCost = 155.0
				calcs["Compute"] = "Azure VMs: 2x D2s_v5 balanceadas"
			}
		default: // AWS
			if tier == "Enterprise" {
				computeCost = ternFloat(isARM, 360.0, 440.0)
				calcs["Compute"] = fmt.Sprintf("AWS EC2 Auto Scaling Group (%s)", tern(isARM, "c7g.xlarge Graviton", "c6i.xlarge"))
			} else if tier == "Startup" {
				computeCost = 42.0
				calcs["Compute"] = "AWS EC2: 2x t4g.small (Graviton2)"
			} else {
				computeCost = ternFloat(isARM, 125.0, 155.0)
				calcs["Compute"] = fmt.Sprintf("AWS EC2: 2x %s balanceadas", tern(isARM, "t4g.xlarge (ARM)", "t3.xlarge"))
			}
		}
	}
	breakdown["Compute"] = computeCost

	// 2. DATABASE CALCULATION
	var dbCost float64
	dbType := params.DatabaseType
	if dbType == "" {
		dbType = "PostgreSQL"
	}
	isMultiAZ := strings.Contains(strings.ToLower(params.DatabaseRedundancy), "multi") || tier == "Enterprise"
	if strings.Contains(strings.ToLower(params.DatabaseRedundancy), "single") {
		isMultiAZ = false
	}

	switch provider {
	case "ORACLE":
		if tier == "Enterprise" {
			dbCost = 280.0
			calcs["Database"] = fmt.Sprintf("OCI Autonomous DB / MySQL HeatWave (Multi-AD HA com 1TB Storage)")
		} else if tier == "Startup" && !isMultiAZ {
			dbCost = 28.0
			calcs["Database"] = fmt.Sprintf("OCI Managed %s (Single Node)", dbType)
		} else {
			dbCost = ternFloat(isMultiAZ, 90.0, 50.0)
			calcs["Database"] = fmt.Sprintf("OCI Managed %s (%s com 256GB NVMe)", dbType, tern(isMultiAZ, "Multi-AD HA", "Single Node"))
		}
	case "GCP":
		if tier == "Enterprise" {
			dbCost = 390.0
			calcs["Database"] = fmt.Sprintf("GCP Cloud SQL %s (High Availability Multi-Zone + Read Replica)", dbType)
		} else if tier == "Startup" && !isMultiAZ {
			dbCost = 38.0
			calcs["Database"] = fmt.Sprintf("GCP Cloud SQL %s (Single-Zone, db-custom-1-3840)", dbType)
		} else {
			dbCost = ternFloat(isMultiAZ, 145.0, 78.0)
			calcs["Database"] = fmt.Sprintf("GCP Cloud SQL %s (%s, db-custom-2-7680)", dbType, tern(isMultiAZ, "Regional HA Multi-Zone", "Zonal Standalone"))
		}
	case "AZURE":
		if tier == "Enterprise" {
			dbCost = 410.0
			calcs["Database"] = fmt.Sprintf("Azure Database for %s (Flexible Server, Zone-Redundant HA)", dbType)
		} else if tier == "Startup" && !isMultiAZ {
			dbCost = 40.0
			calcs["Database"] = fmt.Sprintf("Azure Database for %s (Burstable B2s, 32GB)", dbType)
		} else {
			dbCost = ternFloat(isMultiAZ, 150.0, 80.0)
			calcs["Database"] = fmt.Sprintf("Azure Database for %s (General Purpose D2ds_v5, %s)", dbType, tern(isMultiAZ, "Zone-Redundant HA", "Single Zone"))
		}
	default: // AWS
		if tier == "Enterprise" {
			dbCost = 420.0
			calcs["Database"] = fmt.Sprintf("AWS Aurora %s (Multi-AZ Cluster + 1x Auto-scaling Read Replica)", dbType)
		} else if tier == "Startup" && !isMultiAZ {
			dbCost = 38.0
			calcs["Database"] = fmt.Sprintf("AWS RDS %s (db.t4g.medium Single-AZ, 50GB gp3)", dbType)
		} else {
			dbCost = ternFloat(isMultiAZ, 148.0, 75.0)
			calcs["Database"] = fmt.Sprintf("AWS RDS %s (db.t4g.xlarge %s com Failover)", dbType, tern(isMultiAZ, "Multi-AZ", "Single-AZ"))
		}
	}
	breakdown["Database"] = dbCost

	// 3. CACHE & IN-MEMORY
	if params.CacheEnabled || tier == "Growth" || tier == "Enterprise" {
		var cacheCost float64
		switch provider {
		case "ORACLE":
			cacheCost = ternFloat(tier == "Enterprise", 95.0, 35.0)
			calcs["Cache"] = fmt.Sprintf("OCI Cache with Redis (%s)", tern(tier == "Enterprise", "Cluster 3 nós", "Dedicado 2GB"))
		case "GCP":
			cacheCost = ternFloat(tier == "Enterprise", 140.0, 48.0)
			calcs["Cache"] = fmt.Sprintf("GCP Memorystore for Redis (%s)", tern(tier == "Enterprise", "Standard HA 5GB", "Basic Tier 2GB"))
		case "AZURE":
			cacheCost = ternFloat(tier == "Enterprise", 155.0, 52.0)
			calcs["Cache"] = fmt.Sprintf("Azure Cache for Redis (%s)", tern(tier == "Enterprise", "Premium P1", "Standard C1"))
		default: // AWS
			cacheCost = ternFloat(tier == "Enterprise", 145.0, 46.0)
			calcs["Cache"] = fmt.Sprintf("AWS ElastiCache Redis (%s)", tern(tier == "Enterprise", "Cluster Multi-Node Multi-AZ", "cache.t4g.medium"))
		}
		breakdown["Cache"] = cacheCost
	}

	// 4. STORAGE (OBJECT & DISK)
	var storageCost float64
	switch provider {
	case "ORACLE":
		storageCost = float64(storageGB) * 0.015 // OCI Object Storage $0.0255 w/ generous free tier
		calcs["Storage"] = fmt.Sprintf("OCI Object & Block Storage: %d GB", storageGB)
	case "GCP":
		storageCost = float64(storageGB) * 0.020
		calcs["Storage"] = fmt.Sprintf("GCP Cloud Storage Standard: %d GB", storageGB)
	case "AZURE":
		storageCost = float64(storageGB) * 0.018
		calcs["Storage"] = fmt.Sprintf("Azure Blob Hot Tier: %d GB", storageGB)
	default: // AWS
		storageCost = float64(storageGB) * 0.023
		calcs["Storage"] = fmt.Sprintf("AWS S3 Standard: %d GB", storageGB)
	}
	breakdown["Storage"] = math.Round(storageCost*100) / 100

	// 5. NETWORKING & EGRESS
	var networkCost float64
	switch provider {
	case "ORACLE":
		// OCI includes 10 TB free egress per month!
		networkCost = 15.0 // Flexible Load Balancer only
		calcs["Networking_CDN"] = fmt.Sprintf("OCI Flexible Load Balancer (Egress gratuito até 10TB/mês)")
	case "GCP":
		egressCost := float64(egressGB) * 0.085
		lbCost := 18.0
		networkCost = egressCost + lbCost
		calcs["Networking_CDN"] = fmt.Sprintf("GCP Cloud Armor/CDN + Cloud LB + %d GB Egress", egressGB)
	case "AZURE":
		egressCost := float64(egressGB) * 0.087
		lbCost := 22.0
		networkCost = egressCost + lbCost
		calcs["Networking_CDN"] = fmt.Sprintf("Azure Application Gateway + %d GB Tráfego de Saída", egressGB)
	default: // AWS
		egressCost := float64(egressGB) * 0.09
		albCost := 22.50
		networkCost = egressCost + albCost
		calcs["Networking_CDN"] = fmt.Sprintf("AWS ALB + CloudFront CDN + %d GB Data Transfer Out", egressGB)
	}
	breakdown["Networking_CDN"] = math.Round(networkCost*100) / 100

	total := 0.0
	for _, val := range breakdown {
		total += val
	}

	return breakdown, calcs, total
}

func generateOptimizations(provider, tier string, params EstimateCloudCostParams) []string {
	var opts []string

	if provider == "AWS" {
		opts = append(opts,
			"Migrar instâncias x86 para processadores AWS Graviton3/4 (ARM) para reduzir custos em até 20% com 25% mais performance.",
			"Contratar Compute Savings Plans de 1 ou 3 anos com adiantamento parcial para economia de até 45% em EKS/EC2/Fargate.",
			"Configurar S3 Intelligent-Tiering para mover automaticamente dados frios para Archive Instant Access.",
		)
	} else if provider == "ORACLE" {
		opts = append(opts,
			"Aproveitar o OCI Always Free Tier (4 OCPUs ARM Ampere A1 + 24GB RAM + 10TB de Egress grátis por mês) para ambientes de homologação/staging.",
			"Utilizar OCI Container Engine (OKE) sem taxa de gerenciamento de cluster ($0/mês vs $73/mês na AWS/GCP).",
			"Configurar Autonomous Database Auto-Scaling para desligar OCPUs ociosas fora do horário comercial.",
		)
	} else if provider == "GCP" {
		opts = append(opts,
			"Habilitar Committed Use Discounts (CUDs) para Cloud SQL e Compute Engine com até 37% de desconto.",
			"Usar Cloud Run com alocação de CPU sob demanda para zerar cobrança em horários de baixo tráfego.",
			"Utilizar GKE Autopilot para pagar estritamente pelos Pods alocados sem desperdício de nós ociosos.",
		)
	} else if provider == "AZURE" {
		opts = append(opts,
			"Contratar Azure Reserved Virtual Machine Instances de 1 ano para economia de até 40% em VMs e AKS.",
			"Utilizar Azure Hybrid Benefit se a empresa já possuir licenças existentes de Windows Server/SQL Server.",
			"Adotar Azure Container Apps com KEDA para auto-scale baseado em fila/eventos a custo zero quando inativo.",
		)
	}

	opts = append(opts,
		fmt.Sprintf("Comparativo Multi-Cloud: OCI costuma ser a opção mais barata em tráfego e compute ARM, enquanto %s oferece maior ecossistema de serviços gerenciados.", provider),
	)

	return opts
}

func tern(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func ternFloat(cond bool, a, b float64) float64 {
	if cond {
		return a
	}
	return b
}
