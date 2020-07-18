<p align="right">
 <img src="https://github.com/hellgate75/go-network/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://api.travis-ci.com/hellgate75/go-network.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-network">Check last build on Travis-CI</a>
 </p>

<p align="center">
<image width="150" height="146" src="images/network.png"></image>&nbsp;
<image width="260" height="410" src="images/golang-logo.png">
&nbsp;<image width="150" height="150" src="images/library.png"></image>
</p><br/>
<br/>

# go-network
Go Network Library

## Library content

This library contains following modules.

Packages:
* [Api library](/api) - Api server and client library


### Api library

This module manages the API Server and client.

* [Model](/model/api.go) - Api Server and Client model
* [Server](/api/server.go) - Api Server Implementation
* [Client](/api/client.go) - Api Client Implementation
* [ApiActionBuilder](/api/builders/apiactionbuilder.go) - ApiAction Builder Component
* [ApiCallHandlerBuilder](/api/builders/apicallhandlerbuilder.go) - ApiCallHandler Builder Component
* [ClientConfigBuilder](/api/builders/clientconfigbuilder.go) - ClientConfig Builder Component
* [ServerConfigBuilder](/api/builders/serverconfigbuilder.go) - ServerConfig Builder Component


## DevOps

Build procedures are reported in following sections.



### Build the project

Build command :

```
go build github.com/hellgate75/go-network/...
```



### Test the project

Build command :

```
go test github.com/hellgate75/go-network/...
```


Enjoy the experience.


## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the following email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)
 

