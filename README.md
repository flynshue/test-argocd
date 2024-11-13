# Getting Started
This will set up a local kind cluster with argocd for testing/developing argocd apps

All steps assumes that you're in the `hack/` dir

```bash
cd hack
```

## Delete existing kind clusters
Verify that there are no running clusters and delete them if they exist

```bash
$ kind get clusters
kind
```

```bash
$ kind delete cluster
Deleting cluster "kind" ...
Deleted nodes: ["kind-control-plane"]
```

## Create kind cluster with ingress nginx
```bash
./kind.sh
```

Validate ingress nginx installation
```bash
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

Install example ingress app to test ingress
```bash
kubectl apply -f example-ingress.yaml
```

```bash
$ kubectl -n example-ingress get pod,svc,ingress
NAME          READY   STATUS    RESTARTS   AGE
pod/bar-app   1/1     Running   0          24s
pod/foo-app   1/1     Running   0          24s

NAME                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/bar-service   ClusterIP   10.96.196.105   <none>        8080/TCP   24s
service/foo-service   ClusterIP   10.96.112.155   <none>        8080/TCP   24s

NAME                                        CLASS    HOSTS   ADDRESS   PORTS   AGE
ingress.networking.k8s.io/example-ingress   <none>   *                 80      24s
```

```bash
$ curl http://localhost/foo/hostname
```

**Note: Need to enable ssl-passthrough on ingress-nginx**

## Install ArgoCD
```bash
helm repo add argo-cd https://argoproj.github.io/argo-helm
```

```bash
helm upgrade argocd argo-cd/argo-cd \
--install \
--namespace argocd \
--create-namespace \
-f ./argocd-values.yaml \
--atomic
```

Add argocd hostname to `/etc/hosts`

```bash
sudo vim /etc/hosts
127.0.0.1 argocd.example.com
```

Logging into argocd
**Note: since this is for local testing, the user/password creds are set to admin/argocdadmin in the helm values**

```bash
argocd login argocd.example.com --username admin --password argocdadmin
```

## Bootstrap argocd

```bash
argocd app create apps \
    --dest-namespace argocd \
    --dest-server https://kubernetes.default.svc \
    --repo https://github.com/argoproj/argocd-example-apps.git \
    --path apps
argocd app sync apps
