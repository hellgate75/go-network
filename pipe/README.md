<p align="right">
 <img src="https://github.com/hellgate75/go-network/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://pipe.travis-ci.com/hellgate75/go-network.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-network">Check last build on Travis-CI</a>
 </p>

<p align="center">
<image width="150" height="146" src="../images/network.png"></image>&nbsp;
<image width="260" height="410" src="../images/golang-logo.png">
&nbsp;<image width="150" height="150" src="../images/library.png"></image>
</p><br/>
<br/>

# go-network

Go Network Library


## Pipe library

This module manages the Tcp Network Pipe Nodes (in available modes: Input, Output, Input/Output).

* [Model](/model/pipe.go) - Tcp Network Pipe Nodes model
* [pipe.PipeNode](/pipe/pipenode.go) - Pipe Node Implementation
* [pipe.builders.PipeNodeConfigBuilder](/pipe/builders/pipenodeconfigbuilder.go) - PipeNodeConfig Builder Component



### Pipe Node modes

Tcp Pipe Node component may be of three type:

* *InputPipe* - Reads data from a server and send data to a bytes array message channel
* *OutputPipe* - Sends data to a client reading from a bytes array message channel
* *InputOutputPipe* - Contains both the capabilities, running both the features

Using the wrong message channel will occur and error, because only used channels will be created by the Node.


#### Sample code for Pipe Node in Input Mode

Following code for PipeNode instance, describing steps used for opening the reading tcp channel.

```
	logger := log.NewLogger("Sample Input Pipe Node", log.DEBUG)
	pipeInputNode := pipe.NewPipeNode("Sample Input Pipe Node", log.DEBUG)
	pipeNodeConfig, err := builders.
		NewPipeNodeConfigBuilder().
		WithInHost("", 9997).
		Build()
	if err != nil {
		panic(err)
	}
	pipeInputNode, err = pipeInputNode.Init(pipeNodeConfig)
	if err != nil {
		panic(err)
	}
	err = pipeInputNode.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = pipeInputNode.Stop()
		if err != nil {
			panic(err)
		}
	}()
	pipeInputNode.UntilStarted()
	outputMessageChannel := pipeInputNode.GetInputPipeChannel()
	logger.Debug("Start message network reader ...")
	for {
		select {
			case msg := <-outputMessageChannel:
				logger.Warnf("Received message: %s", string([]byte(msg)))
		}
	}
```

As first we create an instance of the PipeNode using the `pipe.NewPipeNode` function.
Then we create a configuration (the example doesn't require the SSL/TLS encryption, available in the library)
using the fluent builder `pipe.builders.PipeNodeConfigBuilder`  available calling the function 
`pipe.builders.NewPipeNodeConfigBuilder` (we can set In, Out or both Host options).
In order to configure the TcpServer we invoke the component function `pipe.PipeNode.Init` passing the 
node configuration. Now the node is ready to start. 

Pipe Node is ready to start using the function `pipe.PipeNode.Start`. In order to make the main 
thread waiting for the completion of the server activities we invocate at the end of the code the
function `pipe.TcpServer.Wait`, instead to wait for the channels are ready for use you must invoke the
function `pipe.TcpServer.UntilStarted`.

Reading data from the network is easy, just call function `pipe.TcpServer.GetInputPipeChannel` and use the
returned channel to read messages (model.PipeMessage === []byte). 
Place the reading event, loop or anything matching with your desing, and enjoy the powerful library.

The message is format-aware and you can use any content your need in the model.PipeMessage type, accordingly 
to the kind of information you want to transfer in the network.


#### Sample code for Pipe Node in Output Mode

Following code for PipeNode instance, describing steps used for opening the writing tcp channel.

```
	logger := log.NewLogger("Sample Output Pipe Node", log.DEBUG)
	pipeOutputNode := pipe.NewPipeNode("Sample Output Pipe Node", log.DEBUG)
	pipeNodeConfig, err := builders.
		NewPipeNodeConfigBuilder().
		WithOutHost("", 9997).
		Build()
	if err != nil {
		panic(err)
	}
	pipeOutputNode, err = pipeOutputNode.Init(pipeNodeConfig)
	if err != nil {
		panic(err)
	}
	err = pipeOutputNode.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = pipeOutputNode.Stop()
		if err != nil {
			panic(err)
		}
	}()
	pipeOutputNode.UntilStarted()
	inputMessageChannel := pipeOutputNode.GetOutputPipeChannel()
	for {
		msg := createPipeMessage()
		logger.Debugf("Sending message: %s", string([]byte(msg)))
		inputMessageChannel <- msg
		logger.Warnf("Sent message: %s", string([]byte(msg)))
		time.Sleep(2 * time.Second)
	}

```

As first we create an instance of the PipeNode using the `pipe.NewPipeNode` function.
Then we create a configuration (the example doesn't require the SSL/TLS encryption, available in the library)
using the fluent builder `pipe.builders.PipeNodeConfigBuilder`  available calling the function 
`pipe.builders.NewPipeNodeConfigBuilder` (we can set In, Out or both Host options).
In order to configure the TcpServer we invoke the component function `pipe.PipeNode.Init` passing the 
node configuration. Now the node is ready to start. 

Pipe Node is ready to start using the function `pipe.PipeNode.Start`. In order to make the main 
thread waiting for the completion of the server activities we invocate at the end of the code the
function `pipe.TcpServer.Wait`, instead to wait for the channels are ready for use you must invoke the
function `pipe.TcpServer.UntilStarted`.

Reading data from the network is easy, just call function `pipe.TcpServer.GetOutputPipeChannel` and use the
returned channel to send messages (model.PipeMessage === []byte) to the remote pipe node or tcp server. 
Place the writer event, loop or anything matching with your desing, and enjoy the powerful library.

The message is format-aware and you can use any content your need in the model.PipeMessage type, accordingly 
to the kind of information you want to transfer in the network.


Here a sample model.PipeMessage generation function:

```
var pipeMessageCount int64

func createPipeMessage() model.PipeMessage {
	pipeMessageCount++
	return model.PipeMessage([]byte(fmt.Sprintf("This is message # %v", pipeMessageCount)))
}
```


####  Observations

For the Tcp Input/Output Pipe Node, you can mix-in the privious samples and calling in the configuration
both the `pipe.builders.PipeNodeConfigBuilder.WithInHost` and  `pipe.builders.PipeNodeConfigBuilder.WithOutHost`
automatically the builder will merge the pipe types in the `model.InputOutputPipe`.
This action will force the Pipe Node instance, configured with this mode to execute both the
Network reader on the input port and the writer on the Output Address and port, using the tcp protocol.



#### Sample code

Sample code is available at [pipe.go](/sample/pipe.go).



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
 

