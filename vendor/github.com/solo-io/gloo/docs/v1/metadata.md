<a name="top"></a>

## Contents
  - [Metadata](#v1.Metadata)



<a name="metadata"></a>
<p align="right"><a href="#top">Top</a></p>




<a name="v1.Metadata"></a>

### Metadata
Metadata contains general properties of config resources useful to clients and the gloo control plane for purposes of versioning, annotating, and namespacing resources.


```yaml
resource_version: string
namespace: string
annotations: map<string,string>

```
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_version | string |  | ResourceVersion keeps track of the resource version of a config resource. This mechanism is used by [gloo-storage](https://github.com/solo-io/gloo/pkg/storage) to ensure safety with concurrent writes/updates to a resource in storage. |
| namespace | string |  | Namespace is used for the namespacing of resources. Currently unused by gloo internally. |
| annotations | map&lt;string,string&gt; |  | Annotations allow clients to tag resources for special use cases. gloo ignores annotations but preserved them on read/write from/to storage. |





 

 

 

