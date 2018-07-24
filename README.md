# job-controller

golang script, dockerfile, and kubernetes deployment manifest for running a pod
which will spawn kubernetes jobs on a set interval.

## Acknowledgements

* Multi-stage docker build to keep golang docker image size down:
  https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

* k8s cluster role binding to allow the default serviceaccount to create jobs:
  https://github.com/fabric8io/fabric8/issues/6840
