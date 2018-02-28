## Kubernetes for Docker Host DNS
This is the repo for https://hub.docker.com/r/thrawn01/k8-docker-host-dns

## What do?
This container is designed to be run as a kubernetes job or cron job and update
the docker VM host with the appropriate DNS entries such that pods using
`hostNetwork: true` can resolve kubernetes service DNS names.

## How do?
Starts the container as a one time job
```bash
kubectl create -f https://raw.githubusercontent.com/thrawn01/k8-docker-host-dns/master/job.yaml
```

Starts the container as a cron job that will run once every minute. Although
the container runs every minute it only updates the /etc/resolv.conf file if
the pod is unable to resolve the 'kubernetes' service name
```bash
kubectl create -f https://raw.githubusercontent.com/thrawn01/k8-docker-host-dns/master/with-cron.yaml
```

## Rational
See https://github.com/docker/for-mac/issues/2646
