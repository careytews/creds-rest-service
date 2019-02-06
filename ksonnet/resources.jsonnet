
//
// Definition for creds-rest-service resources on Kubernetes.
//

// Import KSonnet library.
local k = import "ksonnet.beta.2/k.libsonnet";

// Short-cuts to various objects in the KSonnet library.
local depl = k.extensions.v1beta1.deployment;
local container = depl.mixin.spec.template.spec.containersType;
local mount = container.volumeMountsType;
local volume = depl.mixin.spec.template.spec.volumesType;
local resources = container.resourcesType;
local env = container.envType;
local gceDisk = volume.mixin.gcePersistentDisk;
local secretDisk = volume.mixin.secret;
local svc = k.core.v1.service;
local svcPort = svc.mixin.spec.portsType;
local svcLabels = svc.mixin.metadata.labels;

local credsRestService(config) = {

    local version = import "version.jsonnet",

    name: "creds-rest-service",
    images: [
	config.containerBase + "/creds-rest-service:" + version
    ],

    // Environment
    envs:: [
        env.new("GOOGLE_APPLICATION_CREDENTIALS", "/key/private.json"),
        env.new("GOOGLE_CLOUD_PROJECT", config.project),
        env.new("CREDENTIALS_BUCKET", "%s-credentials" % [config.project]),
    ],

    // Volume mount points
    volumeMounts:: [
        mount.new("keys", "/key") + mount.readOnly(true)
    ],

    // Container definition.
    containers:: [

        container.new("creds-rest-service", self.images[0]) +
            container.env(self.envs) +
            container.volumeMounts(self.volumeMounts) +
            container.mixin.resources.limits({
                memory: "64M", cpu: "1.0"
            }) +
            container.mixin.resources.requests({
                memory: "64M", cpu: "0.05"
            })

    ],

    // Volumes - this invokes a secret containing the cert/key
    volumes:: [
        volume.name("keys") + secretDisk.secretName("credential-svc-keys")
    ],

    // Deployment definition.  id is the node ID.
    deployments:: [
        depl.new("creds-rest-service", 1, self.containers,
                 {app: "creds-rest-service", component: "frontend"}) +
            depl.mixin.spec.template.spec.volumes(self.volumes) +
            depl.mixin.metadata.namespace(config.namespace)
    ],

    // Ports declared on the service.
    servicePorts:: [
        svcPort.newNamed("http", 8080, 8080) + svcPort.protocol("TCP")
    ],

    services:: [
        // One service
        svc.new("creds-rest-service", {app: "creds-rest-service"}, self.servicePorts) +

        // Label
        svcLabels({app: "creds-rest-service", component: "frontend"}) +

        svc.mixin.metadata.namespace(config.namespace)

    ],

    // Function which returns resource definitions - deployments and services.
    resources:
		if config.options.includeCredentialSvc then
		 	 self.deployments + self.services
		else []
};

// Return the function which creates resources.
[credsRestService]
