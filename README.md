# cloud-run-controller
A controller to populate kubernetes service from GCP Cloud Run

## Usage

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
