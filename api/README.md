<p align="right">
 <img src="https://github.com/hellgate75/go-network/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://api.travis-ci.com/hellgate75/go-network.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-network">Check last build on Travis-CI</a>
 </p>

<p align="center">
<image width="150" height="146" src="../images/network.png"></image>&nbsp;
<image width="260" height="410" src="../images/golang-logo.png">
&nbsp;<image width="150" height="150" src="../images/library.png"></image>
</p><br/>
<br/>

# go-network

Go Network Library


## Api library

This module manages the API Rest Server and client.

* [Model](/model/api.go) - Api Rest Server and Client model
* [api.ApiServer](/api/server.go) - Api Rest Server Implementation
* [api.ApiClient](/api/client.go) - Api Rest Client Implementation
* [api.builders.ApiActionBuilder](/api/builders/apiactionbuilder.go) - ApiAction Builder Component
* [api.builders.ApiCallHandlerBuilder](/api/builders/apicallhandlerbuilder.go) - ApiCallHandler Builder Component
* [api.builders.ClientConfigBuilder](/api/builders/clientconfigbuilder.go) - ClientConfig Builder Component
* [api.builders.ServerConfigBuilder](/api/builders/serverconfigbuilder.go) - ServerConfig Builder Component


### Api Rest Server components

Here main components, used to make working the Api Rest Server component
* [api.ApiServer](/model/api.go) - Api Rest Server type
* [api.NewApiServer](/api/server.go) - Function that creates the Api Rest Server (requires the logger application ame and log verbosity) 
* [api.builders.ServerConfigBuilder](/api/builders/serverconfigbuilder.go) - Fluent Builder for Server Configuration type instance
* [api.builders.ApiCallHandlerBuilder](/api/builders/apicallhandlerbuilder.go) - Fluent Builder for ApiCallHandler type instance
* [api.builders.ApiActionBuilder](/api/builders/apiactionbuilder.go) - Fluent Builder for ApiAction type instance


#### Sample code for ApiServer creation

Following code for ApiServer instance, and definition of an only `POST` for the `root` Rest path

```
	apiServer := api.NewApiServer("Sample Api Rest Server", log.DEBUG)
	serverConfig, err := builders.
		NewServerConfigBuilder().
		WithHost("", 9999).
		Build()
	if err != nil {
		log2.Fatal(err)
	}
	apiServer, err = apiServer.Init(serverConfig)
	if err != nil {
		log2.Fatal(err)
	}
	handler, _ := builders.NewApiCallHandlerBuilder().
		WithPath("/").
		WithWebMethodHandling("POST", builders.NewApiActionBuilder().
			With(func(c context.ApiCallContext) error {
				defer func() {
					if r := recover(); r != nil {
						err = errors.New(fmt.Sprintf("%v", r))
					}
				}()
				req := <sample request structure>
				err = c.ParseBody(&req)
				c.Logger.Infof("POST Data Read Error: %v", err)
				c.Logger.Infof("POST Data: %+v", req)
				res := <empty response structure>
				err = c.WriteResponse(&res, http.StatusOK)
				return err
			}).
			Build()).
		Build()
	err = apiServer.AddPath(handler)
	if err != nil {
		log2.Fatal(err)
	}
	err = apiServer.Start()
	if err != nil {
		log2.Fatal(err)
	}
	apiServer.Wait()
```

As first we create an instance of the ApiServer using the `api.NewApiServer` function.
Then we create a configuration (the example doesn't require the SSL/TLS encryption, available in the library)
using the fluent builder `api.builders.ServerConfigBuilder`  available calling the function 
`api.builders.NewServerConfigBuilder`.
In order to configure the ApiServer we invoke the component function `api.ApiServer.Init()` passing the 
server configuration. Now the server is ready to start. But we need to create the root path rest available in POST web method.
We need an `ApiCallHandler` component instance, available executing the fluent builder `api.builders.ApiCallHandlerBuilder`, we can access this 
component calling the function `api.builders.NewApiCallHandlerBuilder`.
We create an `ApiAction` using the `api.builders.ApiActionBuilder` component, using the function function `api.builders.NewApiActionBuilder`, this builder
provides the function to store the action handler function content.

We can use in this code the [model.context.ApiCallContext](/model/context/apicontext.go) component. This component provides
multiple capabilities and accelerators to read the request, to save the response, to save data in
handler, server or global cache map.

We invoke the `api.ApiServer.AddPath` function to include the new handler in the
Api Rest Server mux router, in the specified context path (`/`) and recorder method (`POST`).

Server is ready to start using the function `api.ApiServer.Start`. In order to make the main 
thread waiting for the completion of the server activities we invocate at the end of the code the
function `api.ApiServer.Wait`.



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



### Api Rest Client components

Here main components, used to make working the Api Rest Client component
* [api.ApiClient](/model/api.go) - Api Rest Client type
* [api.NewApiClient](/api/client.go) - Function that creates the Api Rest Client (requires the logger application ame and log verbosity) 
* [api.builders.ClientConfigBuilder](/api/builders/clientconfigbuilder.go) - Fluent Builder for Client Configuration type instance


#### Sample code for ApiClient creation

Following code for ApiClient instance, and invocation of `POST` web mathod for the `root` Rest path

```
	apiClient := api.NewApiClient("Sample Api Rest Client", log.DEBUG)
	apiClientConfig, err := builders.NewClientConfigBuilder().
		WithHost("http", "localhost", 9999).
		Build()
	if err != nil {
		log2.Fatal(err)
	}
	err = apiClient.Connect(apiClientConfig)
	if err != nil {
		log2.Fatal(err)
	}
	empty := <empty response structure>
	sample := <sample request structure>
	
	b, err := io.Marshal(encoding.EncodingJSONFormat, sample)
	if err != nil {
		log2.Fatal(err)
	}
	reader := bytes.NewBuffer(b)
	mime := encoding.JsonMimeType
	resp, err := apiClient.Call("/", http.MethodPost, &mime, &mime, reader)
	fmt.Printf("Status: %s\n", resp.Status)
	bts, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", resp)
	_ = io.Unmarshal(bts, encoding.EncodingJSONFormat, &empty)
	fmt.Printf("Response data: %+v\n", empty)
```

As first we create an instance of the api.ApiClient using the `NewApiClient` function.
Then we create a configuration (the example doesn't require the SSL/TLS encryption, available in the library)
using the fluent builder `api.builders.ClientConfigBuilder`  available calling the function 
`api.builders.NewClientConfigBuilder`.
In order to connect the ApiClient we invoke the component function `api.ApiClient.Connect` passing the 
client configuration. Now the client is ready to send request to the server. 

As first we marshall the request object and we send the request to the server, using the ApiClient component 
function `api.ApiClient.Call`. Alternatively we can marshal the element to the
provided mime type in the Client Rest call using the component function `api.ApiClient.Encode`.

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

Sample code is available at [api.go](/sample/api.go).



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
 

