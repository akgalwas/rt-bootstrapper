# Architectural Decisions

## Technical Design

Several architectural decisions were made during the Kyma architecture meeting and the implementation phase. These decisions were primarily driven by technical constraints and the need for timely solutions.

### High Level Design

![High Level Architecture](./assets/high-level-arch.drawio.svg)

**Components:**
 
 1. **KIM (Kyma Infrastructure Manager):** Responsible to deploy the webhook and shared resources to SKRs.
 2. **API Server:** The Kubernetes API server calls the manipulation webhook to intercept the Pod manifest before it gets applied.
 3. **RT Bootstrapper:** Modifies Pod manifests and applies landscape specific adjustments (e.g. adding pull-secret or rewriting image-registry host-names etc.).
 4. **Workload:** The manipulated workload is adjusted to the lanscape specific setup can use the shared resources (e.g. pull-secrets, cluster-trust-bundles etc.).
 

## Applied Manipulations

Depending on the landscape, the webhook has to apply different kind of pod manipulations.

The webhook configuration defines per landscape the set of relevant manipulations. 

The following sections describes how webhook has to behave and how the pod modifications will be applied.

### Manipulation is Limited to Pods
The webhook only manipulates Pod resources. Other resources, such as StatefulSets, DaemonSets, and Deployments, are ignored. This is required to avoid conflicts between Kyma Lifecycle Manager (KLM) and Kyma Infrastructure Manager (KIM). KLM regularly processes the resources it deployed (for example, Deployments of operators). If the webhook were to modify these deployments, the KLM would revert the modifications regularly, and both processes would "fight" against each other. To avoid such a situation, we agreed that KLM will never deploy Pods, but high-level resources like Deployments, DaemonSets, StatefulSets, etc. The drawback of this decision is that the deployed Pod can include different values compared to its definition within a Deployment, StatefulSet, DaemonSet, etc., which may be confusing for engineers or developers reviewing a Pod definition in Kubernetes who are unaware of the webhook's existence and its adjustments.

### Non-Blocking Webhook
The admission webhook must be configured as a non-blocking processing step for API-server requests. This means that the API server continues processing the request when the webhook cannot be invoked. This decision ensures that the API server continues to process requests even when the webhook is temporarily unavailable. The decision introduces the risk that Pods get scheduled without being manipulated.

### Detection of Non-Manipulated Resources Is Not Part of the Webhook
We agreed that the webhook is exclusively responsible for manipulating the manifest of Pods during their creation phase. If a Pod gets scheduled without being processed by the webhook (for example, when the webhook is temporarily down), the Pod might miss critical adjustments and, in the worst case, may not start up properly. To address this issue, a housekeeping process implemented outside of the webhook regularly scans all Pods for any missing manipulations. If such Pods are identified, the housekeeping process restarts them (during the re-creation, the webhook is invoked, and the manipulations are applied).

### Opt-In Approach
We agreed that only Pods are processed by the webhook, if one of the following conditions is fulfilled:

1. The configuration of the Webhook defines for the namespace a list of mandatory manipulations. This ensures that any Pod in the Kyma managed namespaces will be processed.
2. The namespace is annotated to received particular manipulations.
3. The Pod itself is annotated to receive manipulations.

This enables also customer to opt-in into this modification mechanism by annotating either the customer owned namespace or the Pod manifests accordingly.

### Webhook Configuration
The webhook retrieves a default configuration that defines the number of manipulations it must apply to all pod Pod within particular namespaces. Customers or other workloads cannot modify this configuration. 

The configuration is per default only considering Kyma managed namespaces (e.g. `kyma-system`, `istio-system` etc.) so that conflicts with customer owned namespaces are avoided.

### Applied Manipulations
The webhook supports multiple manipulations. The default configuration, managed by KIM, determines which manipulation is used.

### Pull Secret Synchronisation
Private container registries require a pull secret. The Kyma bacjend ensures that the latest pull secret becomes available within the `kyma-system` namespace. The name of the pull secret is static and does not change over time, allowing other components to use it as a unique identifier. This pull secret must be replicated across all namespaces. This is required because Kyma workloads, such as Istio sidecars or serverless, can be deployed in any namespace, and pull secrets are namespace-scoped. The webhook includes a dedicated controller that ensures the secret is available and synchronized in all namespaces of the Kyma runtime.
