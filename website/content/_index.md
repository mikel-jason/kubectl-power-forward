---
title: kubectl power-forward
---

{{< blocks/cover title="Port-forwarding on steriods!" image_anchor="top" byline="Photo by Clark Van Der Beken on Unsplash" >}}
<a class="btn btn-lg btn-primary me-3 mb-4" href="/docs/">
  Getting Started <i class="fas fa-arrow-alt-circle-right ms-2"></i>
</a>
<a class="btn btn-lg btn-secondary me-3 mb-4" href="https://github.com/mikel-jason/kubectl-power-forward" target="_blank">
  Source Code <i class="fab fa-github ms-2 "></i>
</a>
{{< /blocks/cover >}}

{{% blocks/lead color="dark" %}}
Take you port-forwarding to the next level! Automatic reconnecting when your pod has been descheduled or using real-world domains, we got you covered!
{{% /blocks/lead %}}

{{% blocks/section color="primary" type="row" %}}
{{% blocks/feature icon="fa-lightbulb" title="Services as primary targets" %}}
By choosing to port-forward to services (multiple!), your connections are automatically reconnected once the workloads behind a service changes. You don't have to manually reconnect.
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-github" title="Use custom domains locally" %}}
Point domains of your choise to you local machine and use your real-world addresses without deploying DNS records, a load balancer etc.
{{% /blocks/feature %}}

{{% blocks/feature icon="fab fa-github" title="Real-life ports with HTTP reverse proxy" %}}
On top of SOCKS, you can use a local reverse proxy to serve all your domains from the same local port.
{{% /blocks/feature %}}



{{% /blocks/section %}}
