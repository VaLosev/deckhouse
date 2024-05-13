{{- define "template-files-values"  }}
files:
 - templateName: static-pod-system-registry.yaml
   filePath: /manifests/static_pods/system-registry.yaml
 - templateName: docker-auth-config.yaml
   filePath: /manifests/auth_config/config.yaml
 - templateName: distribution-config.yaml
   filePath: /manifests/distribution_config/config.yaml
 - templateName: seaweedfs-filer.toml
   filePath: /manifests/seaweedfs_config/filer.toml
 - templateName: seaweedfs-master.toml
   filePath: /manifests/seaweedfs_config/master.toml
{{- end }}

{{- define "docker-auth-config.yaml" }}
server:
  addr: "${discovered_node_ip}:5051"
token:
  issuer: "Registry server"
  expiration: 900
  certificate: "/system_registry_pki/token.crt"
  key: "/system_registry_pki/token.key"
users:
  # Password is specified as a BCrypt hash. Use htpasswd -nB USERNAME to generate.
  "pusher":
    password: "\$2y\$05\$d9Ko2sN9YKSgeu9oxfPiAeopkPTaD65RWQiZtaZ2.hnNnLyFObRne"  # pusher
  "puller":
    password: "\$2y\$05\$wVbhDuuhL/TAVj4xMt3lbeCAYWxP1JJNZJdDS/Elk7Ohf7yhT5wNq"  # puller
acl:
  - match: { account: "pusher" }
    actions: [ "*" ]
    comment: "Pusher has full access to everything."
  - match: {account: "/.+/"}  # Match all accounts.
    actions: ["pull"]
    comment: "readonly access to all accounts"
  # Access is denied by default.
{{- end }}


{{- define "seaweedfs-filer.toml" }}
filer.options:
  recursive_delete: false # do we really need for registry?
etcd:
  enabled: true
  servers: "{{ $.Values.systemRegistry.internal.etcd.addresses | join "," }}"
  key_prefix: "seaweedfs_meta."
  tls_ca_file: "/kubernetes_pki/etcd/ca.crt"
  tls_client_crt_file: "/system_registry_pki/seaweedfs-etcd-client.crt"
  tls_client_key_file: "/system_registry_pki/seaweedfs-etcd-client.key"
{{- end }}


{{- define "seaweedfs-master.toml" }}
master.volume_growth:
  copy_1: 1
  copy_2: 2
  copy_3: 3
  copy_other: 1
{{- end }}


{{- define "distribution-config.yaml" }}
version: 0.1
log:
  level: info
storage:
  s3:
    accesskey: awsaccesskey
    secretkey: awssecretkey
    region: us-west-1
    regionendpoint: http://localhost:8333
    bucket: registry
    encrypt: false
    secure: false
    v4auth: true
    chunksize: 5242880
    rootdirectory: /
    multipartcopy:
      maxconcurrency: 100
      chunksize: 33554432
      thresholdsize: 33554432
  delete:
    enabled: true
  redirect:
    disable: true
  cache:
    blobdescriptor: inmemory
http:
  addr: ${discovered_node_ip}:5000
  #addr: 0.0.0.0:5000
  prefix: /
  secret: asecretforlocaldevelopment
  debug:
    addr: localhost:5001
    prometheus:
      enabled: true
      path: /metrics
proxy:
  username: "$UPSTREAM_REGISTRY_LOGIN"
  password: "$UPSTREAM_REGISTRY_PASSWORD"
  ttl: 72h
auth:
  token:
    realm: http://${discovered_node_ip}:5051/auth
    service: Docker registry
    issuer: Registry server
    rootcertbundle: /system_registry_pki/token.crt
    autoredirect: false
{{- end }}


{{- define "static-pod-system-registry.yaml" }}
apiVersion: v1
kind: Pod
metadata:
  labels:
    component: system-registry
    tier: control-plane
  name: system-registry
  namespace: d8-system
spec:
  dnsPolicy: ClusterFirst
  hostNetwork: true
  containers:
  - name: distribution
    image: "{{ $.Values.global.modulesImages.registry.base }}@{{ $.Values.global.modulesImages.digests.systemRegistry.dockerDistribution }}"
    imagePullPolicy: IfNotPresent
    args:
      - serve
      - /config/config.yaml
    volumeMounts:
      - mountPath: /config
        name: distribution-config-volume
      - mountPath: /system_registry_pki
        name: system-registry-pki-volume
  - name: auth
    image: "{{ $.Values.global.modulesImages.registry.base }}@{{ $.Values.global.modulesImages.digests.systemRegistry.dockerAuth }}"
    imagePullPolicy: IfNotPresent
    args:
      - -logtostderr
      - /config/config.yaml
    volumeMounts:
      - mountPath: /config
        name: auth-config-volume
      - mountPath: /system_registry_pki
        name: system-registry-pki-volume
  - name: seaweedfs
    image: "{{ $.Values.global.modulesImages.registry.base }}@{{ $.Values.global.modulesImages.digests.systemRegistry.seaweedfs }}"
    imagePullPolicy: IfNotPresent
    args:
      - -config_dir="/config"
      - -logtostderr=true
      - -v=0
      - server
      - -filer
      - -s3
      - -dir=/data
      - -volume.max=0
      - -master.volumeSizeLimitMB=1024
      - -metricsPort=9324
      - -volume.readMode=redirect
      - -s3.allowDeleteBucketNotEmpty=true
      - -master.defaultReplication=000
      - -volume.pprof
      - -filer.maxMB=16
      - -ip=${discovered_node_ip}
      - -master.peers=${discovered_node_ip}:9333
    env:
      - name: GOGC
        value: "20"
      - name: GOMEMLIMIT
        value: "500MiB"
    volumeMounts:
      - mountPath: /data
        name: seaweedfs-data-volume
      - mountPath: /config
        name: seaweedfs-config-volume
      - mountPath: /kubernetes_pki
        name: kubernetes-pki-volume
      - mountPath: /system_registry_pki
        name: system-registry-pki-volume
  priorityClassName: system-node-critical
  volumes:
  # PKI
  - name: kubernetes-pki-volume
    hostPath:
      path: /etc/kubernetes/pki
      type: Directory
  - name: system-registry-pki-volume
    hostPath:
      path: /etc/kubernetes/system-registry/pki
      type: Directory
  # Configs
  - name: auth-config-volume
    hostPath:
      path: /etc/kubernetes/system-registry/auth_config
      type: DirectoryOrCreate
  - name: seaweedfs-config-volume
    hostPath:
      path: /etc/kubernetes/system-registry/seaweedfs_config
      type: DirectoryOrCreate
  - name: distribution-config-volume
    hostPath:
      path: /etc/kubernetes/system-registry/distribution_config
      type: DirectoryOrCreate
  # Data
  - name: seaweedfs-data-volume
    hostPath:
      path: /opt/deckhouse/system-registry/seaweedfs_data
      type: DirectoryOrCreate
  # - name: tmp
  #   emptyDir: {}
{{- end }}