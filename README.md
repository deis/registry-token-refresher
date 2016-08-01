# Registry Token Refresher
[![Build Status](https://travis-ci.org/deis/registry-token-refresher.svg?branch=master)](https://travis-ci.org/deis/registry-token-refresher)
[![Docker Repository on Quay](https://quay.io/repository/deisci/registry-token-refresher/status "Docker Repository on Quay")](https://quay.io/repository/deisci/registry-token-refresher)

Deis (pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS) that adds a developer-friendly layer to any [Kubernetes](http://kubernetes.io) cluster, making it easy to deploy and manage applications on your own servers.

For more information about the Deis Workflow, please visit the main project page at https://github.com/deis/workflow.

We welcome your input! If you have feedback, please [submit an issue][issues]. If you'd like to participate in development, please read the "Development" section below and [submit a pull request][prs].

# About
The Registry Token Refresher service creates the [imagePullSecret][imagePullSecrets] and updates it at regular interval of time for each app(namespace) created by Deis Workflow. The secrets are used by [dockerbuilder][dockerbuilder] and [controller][controller] for authenticating with the private registry.

This service is run only when using Amazon's [ECR][ecr] or Google's [GCR][gcr] as they provide short lived tokens for authentication.

# License

Copyright 2013, 2014, 2015, 2016 Engine Yard, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

[issues]: https://github.com/deis/workflow/issues
[prs]: https://github.com/deis/workflow/pulls
[imagePullSecrets]: http://kubernetes.io/docs/user-guide/images/#specifying-imagepullsecrets-on-a-pod
[dockerbuilder]: https://github.com/deis/dockerbuilder
[controller]: https://github.com/deis/controller
[ecr]: http://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_GetStarted.html
[gcr]: https://cloud.google.com/container-registry/
