# Upgrading Rook Ceph

## Contents

- [Introduction](#introduction)
- [Steps](#steps)
  - [Step 1: Make a note of existing image versions](#step-1-make-a-note-of-existing-image-versions)
  - [Step 2: Ensure autoscale is on](#step-2-ensure-autoscale-is-on)
  - [Step 3: Watch ceph status](#step-3-watch-ceph-status)
  - [Step 4: Watch pods in rook namespace](#step-4-watch-pods-in-rook-namespace)
  - [Step 5: Watch rook version update](#step-5-watch-rook-version-update)
  - [Step 6: Watch ceph version update](#step-6-watch-ceph-version-update)
  - [Step 7: Watch events in rook namespace](#step-7-watch-events-in-rook-namespace)
  - [Step 8: Ceph dashboard](#step-8-ceph-dashboard)
  - [Step 9: Grafana dashboard](#step-9-grafana-dashboard)
  - [Step 10: Perform updates](#step-10-perform-updates)
  - [Step 11: Verify that the CSI images are updated](#step-11-verify-that-the-csi-images-are-updated)
  - [Step 12: Final checks](#step-12-final-checks)
- [Additional resources](#additional-resources)

## Introduction

Rook Ceph is one of the storage providers of Lokomotive. With a distributed system as complex as
Ceph, the upgrade can be tricky. This document enlists steps on how to keep an eye on the upgrade
process.

## Steps

Following steps are inspired from [`rook`](https://rook.io/docs/rook/master/ceph-upgrade.html) docs.

### Step 1: Make a note of existing image versions

Execute the following command to list images of the running CSI components:

```bash
kubectl --namespace rook get pod -o \
  jsonpath='{range .items[*]}{range .spec.containers[*]}{.image}{"\n"}' \
  -l 'app in (csi-rbdplugin,csi-rbdplugin-provisioner,csi-cephfsplugin,csi-cephfsplugin-provisioner)' | \
  sort | uniq
```

### Step 2: Ensure autoscale is on

Exec into the toolbox pod as specified in [this
doc](rook-ceph-storage.md#enable-and-access-toolbox). Once you get shell access, run the following
command:

```console
# ceph osd pool autoscale-status | grep replicapool
POOL                     SIZE  TARGET SIZE  RATE  RAW CAPACITY   RATIO  TARGET RATIO  EFFECTIVE RATIO  BIAS  PG_NUM  NEW PG_NUM  AUTOSCALE
replicapool                0                 3.0         3241G  0.0000                                  1.0      32              on
```

Ensure that the `AUTOSCALE` column outputs `on` and not `warn`. If the output of the `AUTOSCALE`
column says `warn`, then run the following command:

```bash
ceph osd pool set replicapool pg_autoscale_mode on
```

### Step 3: Watch ceph status

Leave the following running in the toolbox pod:

```bash
watch ceph status
```

Ensure that the output says that `health:` is `HEALTH_OK`. Match the output such that everything
looks fine as explained in the [rook upgrade
docs](https://rook.io/docs/rook/master/ceph-upgrade.html#status-output).


### Step 4: Watch pods in rook namespace

Keep an eye on the pods status in another terminal window from the `rook` namespace. Leave the
following command running:

```bash
watch kubectl -n rook get pods -o wide
```

### Step 5: Watch rook version update

Run the following command in a new terminal window to keep an eye on the rook version update as it
is upgrades for all the sub-components:

```bash
watch --exec kubectl -n rook get deployments -l rook_cluster=rook -o \
    jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \trook-version="}{.metadata.labels.rook-version}{"\n"}{end}'
```

You should see that `rook-version` slowly changes to `v1.4.6`.

### Step 6: Watch ceph version update

Run the following command to keep an eye on the Ceph version update as the new pods come up in a new
terminal window:

```bash
watch --exec kubectl -n rook get deployments -l rook_cluster=rook -o \
    jsonpath='{range .items[*]}{.metadata.name}{"  \treq/upd/avl: "}{.spec.replicas}{"/"}{.status.updatedReplicas}{"/"}{.status.readyReplicas}{"  \tceph-version="}{.metadata.labels.ceph-version}{"\n"}{end}'
```

You should see that `ceph-version` slowly changes to `15.2.5`.

### Step 7: Watch events in rook namespace

In a new terminal leave the following command running, to keep track of the events happening in the
rook namespace:

```bash
kubectl -n rook get events -w
```

### Step 8: Ceph dashboard

Open the Ceph dashboard in a browser window. Instructions to access the dashboard can be found
[here](rook-ceph-storage.md#access-the-ceph-dashboard).

> **NOTE**: Accessing dashboard can be a hassle because while the components are upgrading you may
> lose access to it multiple times.

### Step 9: Grafana dashboard

Gain access to the Grafana dashboard as instructed
[here](monitoring-with-prometheus-operator.md#access-grafana). And keep an eye on the dashboard
named `Ceph - Cluster`.

> **NOTE**: The data in the Grafana dashboard will always be outdated compared to the `watch ceph
> status` running inside the toolbox pod.


### Step 10: Perform updates

```bash
kubectl apply -f https://raw.githubusercontent.com/kinvolk/lokomotive/master/assets/charts/components/rook/templates/resources.yaml
lokoctl component apply rook rook-ceph
```

### Step 11: Verify that the CSI images are updated

Run the same command as [Step 1](#step-1-make-a-note-of-existing-image-versions) to verify if the
images were updated:

```bash
kubectl --namespace rook get pod -o \
  jsonpath='{range .items[*]}{range .spec.containers[*]}{.image}{"\n"}' \
  -l 'app in (csi-rbdplugin,csi-rbdplugin-provisioner,csi-cephfsplugin,csi-cephfsplugin-provisioner)' | \
  sort | uniq
```

### Step 12: Final checks

Once everything is up to date then run following commands in the toolbox pod:

```bash
ceph osd status
ceph df
rados df
```

## Additional resources

- Rook Upgrade docs: https://rook.io/docs/rook/v1.4/ceph-upgrade.html.
