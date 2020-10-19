# Upgrading and migrating OpenEBS data plane components

## Contents

* [Introduction](#introduction)
* [Prerequisites](#prerequisites)
* [Step 1: Upgrade OpenEBS control plane components](#step-1-upgrade-openebs-control-plane-components)
* [Step 2: Upgrade OpenEBS data plane components](#step-2-upgrade-openebs-data-plane-components)
	* [Upgrade cStor pools](#upgrade-cstor-pools)
	* [Upgrade cStor volumes](#upgrade-cstor-volumes)
* [Step 3: Verify](#step-3-verify)
* [Summary](#summary)
* [Additional resources](#additional-resources)

## Introduction

This guide provides installation steps to upgrade OpenEBS dataplane components
from `v1.12.0` to `v2.2.0` and migrate from cStor Pools/Volumes to latest CSI based provisioning
(CSPC Pools/Volumes).

## Prerequisites

To upgrade the OpenEBS data plane components, we need the following:

* Lokomotive `v0.6.0`
    ```bash
    lokoctl version
    v0.6.0
    ```

* A Lokomotive cluster accessible via `kubectl` deployed on a [Equinix]
    metal(https://metal.equinix.com/).

* OpenEBS `v1.12.0` installed. You can check if OpenEBS is indeed in expected
    version:
    ```bash
    $ kubectl get pods -n openebs -n openebs.io/version=1.12.0
    ```

* Highly recommended to schedule a downtime for the applications consuming the
    OpenEBS PV and make sure to take a backup of the data before starting the
    below upgrade procedure.


## Steps

###Step 1: Upgrade OpenEBS control plane components

Lokomotive provides an easy way of upgrading the OpenEBS control plane.

Execute the following command to upgrade OpenEBS control plane:

```bash
lokoctl component apply openebs-operator
```

Doing so, terminates the OpenEBS resources associated with `v1.12.0` and new
resources are created with `v2.2.0`.

Verify all the pods are in `Running` state before proceeding:

```bash
$ kubectl get pods -n openebs -n openebs.io/version=2.2.0
```

### Step 2: Upgrade OpenEBS data plane components

OpenEBS control plane and data plane components work independently. Even if the
control plane components are upgraded to `v2.2.0`, the data plane components
continue to work with the older version `v1.12.0`.

#### Upgrade cStor pools

Get existing `StoragePoolClaims` by executing:

```bash
kubectl get spc
NAME                           AGE
cstor-pool-openebs-replica-3   2d21h26m
```

Create a Job to upgrade the existing cStor pools:

```bash
# upgrade-cstor-pools-1120-220.yaml

#This is an example YAML for upgrading cstor SPC.
#Some of the values below needs to be changed to
#match your openebs installation. The fields are
#indicated with VERIFY
---
apiVersion: batch/v1
kind: Job
metadata:
  #VERIFY that you have provided a unique name for this upgrade job.
  #The name can be any valid K8s string for name. This example uses
  #the following convention: cstor-spc-<flattened-from-to-versions>
  name: cstor-spc-1120220

  #VERIFY the value of namespace is same as the namespace where openebs components
  # are installed. You can verify using the command:
  # `kubectl get pods -n <openebs-namespace> -l openebs.io/component-name=maya-apiserver`
  # The above command should return status of the openebs-apiserver.
  namespace: openebs
spec:
  backoffLimit: 4
  template:
    spec:
      serviceAccountName: openebs-operator
      containers:
      - name:  upgrade
        args:
        - "cstor-spc"

        # --from-version is the current version of the pool
        - "--from-version=1.12.0"

        # --to-version is the version desired upgrade version
        - "--to-version=2.2.0"

        # Bulk upgrade is supported
        # To make use of it, please provide the list of SPCs
        # as mentioned below
        - "cstor-pool-openebs-replica-3"

        #Following are optional parameters
        #Log Level
        - "--v=4"
        #DO NOT CHANGE BELOW PARAMETERS
        env:
        - name: OPENEBS_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        tty: true

        # the image version should be same as the --to-version mentioned above
        # in the args of the job
        image: openebs/m-upgrade:2.2.0
        imagePullPolicy: Always
      restartPolicy: OnFailure
---
```

Apply the job to start the upgrade for cStor pools:

```bash
kubectl apply upgrade-cstor-pools-1120-220.yaml
```

Ensure the job runs to completion.


#### Upgrade cStor volumes

Get existing `CStorVolumes` by executing:

```bash
k get cstorvolumes --all-namespaces
NAMESPACE   NAME                                       STATUS    AGE     CAPACITY
openebs     pvc-5b7d8ee3-59b8-4a02-bd6d-a4e66fbecf9f   Healthy   2d21h   50Gi
openebs     pvc-e8f8b268-b81c-4201-841a-283f557c44a7   Healthy   2d21h   50Gi
```


Create a Kubernetes job to upgrade the existing cStor volumes:

```bash
# upgrade-cstor-vols-1120-220.yaml

#This is an example YAML for upgrading cstor volume.
#Some of the values below needs to be changed to
#match your openebs installation. The fields are
#indicated with VERIFY
---
apiVersion: batch/v1
kind: Job
metadata:
  #VERIFY that you have provided a unique name for this upgrade job.
  #The name can be any valid K8s string for name. This example uses
  #the following convention: cstor-vol-<flattened-from-to-versions>
  name: cstor-vol-1120220

  #VERIFY the value of namespace is same as the namespace where openebs components
  # are installed. You can verify using the command:
  # `kubectl get pods -n <openebs-namespace> -l openebs.io/component-name=maya-apiserver`
  # The above command should return status of the openebs-apiserver.
  namespace: openebs

spec:
  backoffLimit: 4
  template:
    spec:
      serviceAccountName: openebs-operator
      containers:
      - name:  upgrade
        args:
        - "cstor-volume"

        # --from-version is the current version of the volume
        - "--from-version=1.12.0"

        # --to-version is the version desired upgrade version
        - "--to-version=2.2.0"

        # Bulk upgrade is supported from 1.9
        # To make use of it, please provide the list of PVs
        # as mentioned below
        - "pvc-5b7d8ee3-59b8-4a02-bd6d-a4e66fbecf9f"
        - "pvc-e8f8b268-b81c-4201-841a-283f557c44a7"

        #Following are optional parameters
        #Log Level
        - "--v=4"
        #DO NOT CHANGE BELOW PARAMETERS
        env:
        - name: OPENEBS_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        tty: true

        # the image version should be same as the --to-version mentioned above
        # in the args of the job
        image: openebs/m-upgrade:2.2.0
        imagePullPolicy: Always
      restartPolicy: OnFailure
---
```

```bash
kubectl apply upgrade-cstor-vols-1120-220.yaml
```

Ensure the job runs to completion.

### Step 3: Verify

To check if the upgrade process was successful, execute:

```bash
$ kubectl get pods -n openebs -n openebs.io/version=2.2.0
NAME                                                              READY   STATUS    RESTARTS   AGE
cstor-pool-openebs-replica-3-6qxo-bd59f5554-d97l4                 3/3     Running   0          27m
cstor-pool-openebs-replica-3-ag9v-78cfbb64b4-vzcs9                3/3     Running   0          26m
cstor-pool-openebs-replica-3-gt7g-8ddd9b457-jqwz4                 3/3     Running   0          25m
openebs-operator-admission-server-6b5fb6dff5-6dc8d                1/1     Running   2          34m
openebs-operator-apiserver-5c467bc588-fsl58                       1/1     Running   0          34m
openebs-operator-localpv-provisioner-74d76d55b-s8zwv              1/1     Running   0          33m
openebs-operator-ndm-cfnxp                                        1/1     Running   0          34m
openebs-operator-ndm-nsqkg                                        1/1     Running   0          33m
openebs-operator-ndm-operator-758fdbc5f4-5qqxx                    1/1     Running   0          34m
openebs-operator-ndm-x44d8                                        1/1     Running   0          34m
openebs-operator-provisioner-59c6dc5dfc-q9k65                     1/1     Running   0          34m
openebs-operator-snapshot-operator-bf49c5dc6-xqqlq                2/2     Running   0          34m
pvc-5b7d8ee3-59b8-4a02-bd6d-a4e66fbecf9f-target-68984f69b7l42sl   3/3     Running   0          22m
pvc-e8f8b268-b81c-4201-841a-283f557c44a7-target-7484b968c4pzld6   3/3     Running   0          20m
```

To check if all the `StoragePoolClaims` and `CStorVolumes` have been upgraded,
execute:

```bash
$ kubectl get pods -n openebs -n openebs.io/version=1.12.0
```
No output should be displayed.

## Summary

This guide helps to upgrade the OpenEBS control plane and data plane components.

## Additional resources

[OpenEBS troubleshooting guide](https://docs.openebs.io/docs/next/troubleshooting.html).

For additional information regarding the upgrade steps, see [OpenEBS upgrade
documentation](https://github.com/openebs/openebs/blob/master/k8s/upgrades/README.md).

