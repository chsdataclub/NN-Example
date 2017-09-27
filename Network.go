package main

import (
	"fmt"
)

//NOTE most of the calculating work is networked by nodes inside the struct

type Network struct {
	nodeList []Node //master list of nodes
	numConnections int
	innovation []int //list of inovation numbers this network has (SORTED)
	id int //network id
	learningRate float64 //learning rate for backprop
	output []*Node //output nodes
	input []*Node //input nodes
	fitness float64
	adjustedFitness float64
	numInnovation int
	networkId int
}

//processes the network
func (n *Network) Process(input []float64) {
	for i := 0; i < len(n.input); i++ {
		n.input[i].setValue(input[i])
	}
}

//todo test when time
//backpropogates the network to desired one time
func (n *Network) BackProp(input []float64, desired []float64) {
	n.Process(input) //need to do so that you are performing the algorithm on that set of values

	//this will calc all the influence
	for i := 0; i < len(n.output); i++ {
		n.output[i].setInfluence(n.output[i].value-desired[i])
	}

	//actually adjusts the weights
	for i := 0; i < len(n.nodeList); i++ {
		derivative := sigmoidDerivative(n.nodeList[i].value)
		for a := 0; a < len(n.nodeList[i].receive); a++ {
			n.nodeList[i].receive[a].nextWeight +=  derivative * (n.nodeList[i].receive[a].nodeFrom.value) * n.nodeList[i].influence * n.learningRate
		}
	}
	//backprop output and hidden
	/*for z := 2; z >= 1; z++ {
		for i := 0; i < len(n.nodes[z]); i++ {
			node := n.nodes[z][i]

			node.influence = 0
			derivative := sigmoidDerivative(node.value)

			if z < 2 {
				for a := 0; a < len(node.receive); a++ {
					node.influence += (*node.receive[a].connectInfluence) * (node.receive[a].weight)
				}
			} else {
				node.influence = node.value-desired[i]
			}

			for a := 0; a < len(node.receive); a++ {
				node.receive[a].nextWeight += derivative * (*node.receive[a].sendValue) * node.influence * n.learningRate
			}
		}
	}*/
}

//todo test
func (n *Network) mutateConnection(from int, to int) {
	if len(n.nodeList[len(n.nodeList)-from].send) <= (1+n.nodeList[len(n.nodeList)-from].numConOut) {
		n.nodeList[len(n.nodeList)-from].send = append(n.nodeList[len(n.nodeList)-from].send, Connection{weight: 1, nextWeight: 0, disable: false, nodeFrom: &n.nodeList[len(n.nodeList)-from], nodeTo: &n.nodeList[len(n.nodeList)-to]})
	} else {
		n.nodeList[from].send[len(n.nodeList[from].send)-1] = Connection{weight: 1, nextWeight: 0, disable: false, nodeFrom: &n.nodeList[from], nodeTo: &n.nodeList[to]}
	}

	if len(n.nodeList[len(n.nodeList)-to].receive) <= (1+n.nodeList[len(n.nodeList)-to].numConIn) {
		n.nodeList[len(n.nodeList)-to].receive = append(n.nodeList[len(n.nodeList)-to].receive,  &n.nodeList[len(n.nodeList)-from].send[len( n.nodeList[len(n.nodeList)-from].send)-1])
	} else {
		n.nodeList[to].receive[len(n.nodeList[to].receive)-1] =  &n.nodeList[from].send[len( n.nodeList[from].send)-1]
	}

	n.nodeList[len(n.nodeList)-to].numConIn++
	n.nodeList[len(n.nodeList)-from].numConOut++
}

func (n *Network) addInnovation(num int) {
	if len(n.innovation) <= n.numConnections+1 {
		n.innovation = append(n.innovation,  num)
	} else {
		n.innovation[n.numInnovation+1] =  num
	}
	n.numInnovation++
}
//todo test
/*
change from nodes connection to one with new node
change to nodes pointer to one sent by by new node
 */
func (n *Network) mutateNode(from int, to int) int {
	fromNode := &n.nodeList[len(n.nodeList)-from]
	toNode := &n.nodeList[len(n.nodeList)-to]
	newNode := &n.nodeList[len(n.nodeList)-(1+n.createNode().id)]

	//creates and modfies the connection to the toNode
	for i := 0; i < len(toNode.receive); i++ {
		if fromNode == toNode.receive[len(toNode.receive)-i].nodeFrom { //compares the memory location
			newNode.send[0] = Connection{weight: 1, nextWeight: 0, disable:false, nodeFrom: newNode, nodeTo:toNode}
			toNode.receive[i] = &newNode.send[0]
		}
	}

	//todo find a better way?
	for i := 0; i < len(fromNode.send); i++ {
		if fromNode.send[i].nodeTo == toNode {
			fromNode.send[i].nodeTo = newNode
			newNode.receive[0] = &fromNode.send[i]
		}
	}

	return newNode.id
}

func (n *Network) createNode() *Node {
	node := Node {value:0, influenceRecieved: 0, inputRecieved: 0, id:n.id, receive:make([]*Connection, len(n.input)), send:make([]Connection, len(n.output))}
	n.id++

	if (node.id+1) >= len(n.nodeList) {
		n.nodeList = append(n.nodeList,  node)
	} else {
		n.nodeList[len(n.nodeList)-(1+node.id)] =  node
	}

	return &n.nodeList[len(n.nodeList)-(1+node.id)]
}

func GetNetworkInstance(input int, output int, id int) Network {
	n := Network{networkId: id, id: 0, learningRate: .1, numConnections:0, nodeList:make([]Node, (input+output)*2), output: make([]*Node, output), input: make([]*Node, input)}

	count := 1

	fmt.Print("initialized")

	//create output nodes
	for i := 0; i < output; i++ {
		n.nodeList[len(n.nodeList)-count] = Node {value:0, influenceRecieved: 0, inputRecieved: 0, id:n.id, receive:make([]*Connection, input)}
		n.output[i] = &n.nodeList[count]
		count++
		n.id++
	}
	fmt.Print("output")

	//creates the input nodes and adds them to the network
	for i := 0; i < input; i++ {
		n.nodeList[len(n.nodeList)-count] = Node {value:0, id:n.id, influenceRecieved: 0, inputRecieved: 0, send:make([]Connection, output)}
		n.input[i] = &n.nodeList[count]

		//creates the connections
		for a := 0; a < output; a++ {
			n.nodeList[len(n.nodeList)-count].send[len(n.nodeList[len(n.nodeList)-count].send)-(1+a)] = Connection{weight: 1, nextWeight: 0, disable:false, nodeFrom: n.input[i], nodeTo: n.output[a]}
			n.nodeList[a].receive[i] = &n.nodeList[len(n.nodeList)-count].send[len(n.nodeList[len(n.nodeList)-count].send)-(1+a)]
			n.numConnections++
		}

		n.id++
		count++
	}
	fmt.Print("input")

	return n
}

func (n *Network) getNode(i int) *Node {
	return &n.nodeList[len(n.nodeList)-i-1]
}
