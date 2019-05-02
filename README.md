# cloud-run-controller
A controller to populate kubernetes service from GCP Cloud Run

## Usage

### Service creation

```yaml
apiVersion: cloud-run-controller.barpilot.io/v1alpha1
kind: Service
metadata:
  name: example-service
  annotation:
    removeOnDelete: "true"
spec:
  project: myProject
  location: us-central1
  service:
    apiVersion: serving.knative.dev/v1alpha1
    kind: Service
    metadata:
      name: service-nginx
    spec:
      runLatest:
        configuration:
          revisionTemplate:
            spec:
              container:
                image: nginx
                resources:
                  limits:
                    memory: 512Mi
              timeoutSeconds: 300
              containerConcurrency: 80
```

### k8s service sync

```yaml
apiVersion: cloud-run-controller.barpilot.io/v1alpha1
kind: RunServices
metadata:
  name: example-runservices
spec:
  project: myProject
```

this will create (and maintain)
```
$ kubectl get svc
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP                             PORT(S)    AGE
[...]
myruninstance          ExternalName   <none>          myruninstance-tn4meegzeq-uc.a.run.app   <none>     5m32s
[...]
```

## Status

ALPHA
