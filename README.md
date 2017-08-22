
|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Deis Workflow will soon no longer be maintained.<br />Please [read the announcement](https://deis.com/blog/2017/deis-workflow-final-release/) for more detail. |
|---:|---|
| 09/07/2017 | Deis Workflow [v2.18][] final release before entering maintenance mode |
| 03/01/2018 | End of Workflow maintenance: critical patches no longer merged |

# Registry Token Refresher
[![Build Status](https://ci.deis.io/job/registry-token-refresher/badge/icon)](https://ci.deis.io/job/registry-token-refresher)
[![Docker Repository on Quay](https://quay.io/repository/deisci/registry-token-refresher/status "Docker Repository on Quay")](https://quay.io/repository/deisci/registry-token-refresher)

Deis (pronounced DAY-iss) Workflow is an open source Platform as a Service (PaaS) that adds a developer-friendly layer to any [Kubernetes](http://kubernetes.io) cluster, making it easy to deploy and manage applications on your own servers.

For more information about the Deis Workflow, please visit the main project page at https://github.com/deis/workflow.

We welcome your input! If you have feedback, please [submit an issue][issues]. If you'd like to participate in development, please read the "Development" section below and [submit a pull request][prs].

# About
The Registry Token Refresher service creates the [imagePullSecret][imagePullSecrets] and updates it at regular interval of time for each app(namespace) created by Deis Workflow. The secrets are used by [dockerbuilder][dockerbuilder] and [controller][controller] for authenticating with the private registry.

This service is run only when using Amazon's [ECR][ecr] or Google's [GCR][gcr] as they provide short lived tokens for authentication.

[issues]: https://github.com/deis/workflow/issues
[prs]: https://github.com/deis/workflow/pulls
[imagePullSecrets]: http://kubernetes.io/docs/user-guide/images/#specifying-imagepullsecrets-on-a-pod
[dockerbuilder]: https://github.com/deis/dockerbuilder
[controller]: https://github.com/deis/controller
[ecr]: http://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_GetStarted.html
[gcr]: https://cloud.google.com/container-registry/
[v2.18]: https://github.com/deis/workflow/releases/tag/v2.18.0
