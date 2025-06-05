# API Deployment Override

## 1. Create Kustomize config
```yml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../shared/base/api/deployment

# Renames metadata.name, result: "processor-api"
namePrefix: processor-

# Centralizes images configuration
images:
  - name: charmingruby/podummy
    newName: charmingruby/pipoprocessor
    newTag: latest

# General changes on a new file to me more precise
patches:
  - path: deployment-patch.yml
    target:
      kind: Deployment
      name: api
```

## 2. Create deployment patch
```yml
- op: replace
  path: /metadata/namespace
  value: refinery
- op: replace
  path: /spec/selector/matchLabels/app
  value: processor-api
- op: replace
  path: /spec/template/metadata/labels/app
  value: processor-api
- op: replace
  path: /spec/template/spec/containers/0/name
  value: processor-api
- op: replace
  path: /spec/template/spec/containers/0/ports/0/containerPort
  value: 3001
```