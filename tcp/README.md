<p align="right">
 <img src="https://github.com/hellgate75/go-network/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://tcp.travis-ci.com/hellgate75/go-network.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-network">Check last build on Travis-CI</a>
 </p>

<p align="center">
<image width="150" height="146" src="../images/network.png"></image>&nbsp;
<image width="260" height="410" src="../images/golang-logo.png">
&nbsp;<image width="150" height="150" src="../images/library.png"></image>
</p><br/>
<br/>

# go-network

Go Network Library


## Tcp library

This module manages the Tcp Server and client.

* [Model](/model/tcp.go) - Tcp Server and Client model
* [tcp.TcpServer](/tcp/server.go) - Tcp Server Implementation
* [tcp.TcpClient](/tcp/client.go) - Tcp Client Implementation
* [tcp.builders.TcpActionBuilder](/tcp/builders/tcpactionbuilder.go) - TcpAction Builder Component
* [tcp.builders.TcpCallHandlerBuilder](/tcp/builders/tcpcallhandlerbuilder.go) - TcpCallHandler Builder Component
* [tcp.builders.ClientConfigBuilder](/tcp/builders/clientconfigbuilder.go) - TcpClientConfig Builder Component
* [tcp.builders.ServerConfigBuilder](/tcp/builders/serverconfigbuilder.go) - TcpServerConfig Builder Component



### Tcp Server components

Here main components, used to make working the Tcp Server component
* [tcp.TcpServer](/model/tcp.go) - Tcp Server type
* [tcp.NewTcpServer](/tcp/server.go) - Function that creates the Tcp Server (requires the logger application ame and log verbosity) 
* [tcp.builders/ServerConfigBuilder](/tcp/builders/serverconfigbuilder.go) - Fluent Builder for Tcp Server Configuration type instance
* [tcp.builders/TcpCallHandlerBuilder](/tcp/builders/tcpcallhandlerbuilder.go) - Fluent Builder for TcpCallHandler type instance
* [tcp.builders/TcpActionBuilder](/tcp/builders/tcpactionbuilder.go) - Fluent Builder for TcpAction type instance


#### Sample code for TcpServer creation

Following code for TcpServer instance, and definition of a single connection handler and a single function.

```
	logger := log.NewLogger("Sample Tcp Server", log.DEBUG)
	tcpServer := tcp.NewTcpServer("Sample Tcp Server", log.DEBUG)
	serverConfig, err := builders.
		NewTcpServerConfigBuilder().
		WithNetwork("tcp").
		WithEncoding(encoding.EncodingJSONFormat).
		WithHost("", 9998).
		Build()
	if err != nil {
		log2.Fatal(err)
	}
	logger.Infof("Server Config Encoding: %v", serverConfig.Encoding)
	tcpServer, err = tcpServer.Init(serverConfig)
	if err != nil {
		log2.Fatal(err)
	}
	handler, _ := builders.NewTcpCallHandlerBuilder().
		WithName("Main").
		WithTcpHandling(builders.NewTcpActionBuilder().
			WithName("read-sample-data").
			With(func(c context.TcpContext) error {
				logger.Info("Request handler Main.read-sample-data")
				defer func() {
					if r := recover(); r != nil {
						err = errors.New(fmt.Sprintf("%v", r))
					}
				}()
				logger.Debug("Reading request ...")
				req := emptyStruct()
				err = c.ParseRequest(&req)
				c.Logger.Debugf("Main Data Read Error: %v", err)
				c.Logger.Debugf("Main Data: %+v", req)
				res := sampleStruct()
				c.Logger.Debug("Writing the answer ...")
				err = c.WriteResponse(&res)
				if err != nil {
					time.Sleep(1 * time.Second)
					c.Logger.Errorf("Error during answer sending procedure: %v", err)
				} else {
					c.Logger.Info("Answer sent!!")
				}
				return err
			}).
			Build()).
		Build()
	err = tcpServer.AddPath(handler)
	if err != nil {
		log2.Fatal(err)
	}
	err = tcpServer.Start()
	if err != nil {
		log2.Fatal(err)
	}
	tcpServer.Wait()
```

As first we create an instance of the TcpServer using the `tcp.NewTcpServer` function.
Then we create a configuration (the example doesn't require the SSL/TLS encryption, available in the library)
using the fluent builder `tcp.builders.TcpServerConfigBuilder`  available calling the function 
`tcp.builders.NewTcpServerConfigBuilder`.
In order to configure the TcpServer we invoke the component function `tcp.TcpServer.Init` passing the 
server configuration. Now the server is ready to start. But we need to create the Main handler and the read-sample-data function.
We need an `TcpCallHandler` component instance, available executing the fluent builder `tcp.builders.TcpCallHandlerBuilder`, we can access this 
component calling the function `tcp.builders.NewTcpCallHandlerBuilder`.
We create a `TcpAction` using the `tcp.builders.TcpActionBuilder` component, using the function function `tcp.builders.NewTcpActionBuilder`, this builder
provides the function to store the action handler function content.

We can use in this code the [model.context.TcpContext](/model/context/tcpcontext.go) component. This component provides
multiple capabilities and accelerators to read the request, to save the response, to save data in
handler, server or global cache map.

We invoke the `tcp.TcpServer.AddPath` function to include the new handler in the
Tcp Server mux router, in the specified context path (`/`) and recorder method (`POST`).

Server is ready to start using the function `tcp.TcpServer.Start`. In order to make the main 
thread waiting for the completion of the server activities we invocate at the end of the code the
function `tcp.TcpServer.Wait`.



Here a sample Request/Response model structure :

```
struct{
		Id			string		`yaml:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
		Name		string		`yaml:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
		Surname		string		`yaml:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
		Age			int			`yaml:"age,omitempty" json:"age,omitempty" xml:"age,omitempty"`
	}{
	}
```



### Tcp Client components

Here main components, used to make working the Tcp Client component

* [tcp/TcpClient](/model/tcp.go) - Tcp Client type
* [tcp/NewTcpClient](/tcp/client.go) - Function that creates the Tcp Client (requires the logger application ame and log verbosity) 
* [tcp/builders/ClientConfigBuilder](/tcp/builders/clientconfigbuilder.go) - Fluent Builder for Client Configuration type instance


#### Sample code for TcpClient creation

Following code for TcpClient instance, and invocation of `POST` web mathod for the `root` Rest path

```
	logger := log.NewLogger("Sample Tcp Client", log.DEBUG)
	tcpClient := tcp.NewTcpClient("Sample Tcp Client", log.DEBUG)
	tcpClientConfig, err := builders.NewTcpClientConfigBuilder().
		WithHost("localhost", 9998).
		WithNetwork("tcp").
		WithEncoding(encoding.EncodingJSONFormat).
		Build()
	if err != nil {
		log2.Fatal(err)
	}
	logger.Info("Connecting server ...")
	err = tcpClient.Connect(tcpClientConfig)
	if err != nil {
		log2.Fatal(err)
	}
	defer func() {
		logger.Info("Closing connection ...")
		_ = tcpClient.Close()
	}()
	empty := emptyStruct()
	sample := sampleStruct()
	time.Sleep(1 * time.Second)
	logger.Infof("Request data: %v", sample)
	err = tcpClient.Encode(&sample, &empty, 0 * time.Second)
	logger.Infof("Response data: %+v", empty)
```

As first we create an instance of the tcp.TcpClient using the `NewTcpClient` function.
Then we create a configuration (the example doesn't require the SSL/TLS encryption, available in the library)
using the fluent builder `tcp.builders.ClientConfigBuilder`  available calling the function 
`tcp.builders.NewClientConfigBuilder`.
In order to connect the TcpClient we invoke the component function `tcp.TcpClient.Connect` passing the 
client configuration. Now the client is ready to send request to the server. 

As first we marshall the request object and we send the request to the server, using the TcpClient component 
function `tcp.TcpClient.Call`. Alternatively we can marshal the element to the
provided mime type in the Client Rest call using the component function `tcp.TcpClient.Encode`.

The response is available reading the response data, and marshalling at the request accept mime type.

Here a sample Request/Response model structure :

```
struct{
		Id			string		`yaml:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
		Name		string		`yaml:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
		Surname		string		`yaml:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
		Age			int			`yaml:"age,omitempty" json:"age,omitempty" xml:"age,omitempty"`
	}{
	}
```


#### Sample code

Sample code is available at [tcp.go](/sample/tcp.go).



## DevOps

Build procedures are reported in following sections.


### Create the sample executable

Install sample command :

```
go install github.com/hellgate75/go-network/sample/...
```


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
 

