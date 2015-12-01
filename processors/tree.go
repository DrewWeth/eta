package main

import (
	"time"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Node struct {
	name     string `json:"name"`
	children []*Node
}

type Obj struct{
	ID		string	`json:"ID"`
	Data 		string	`json:"data"`
	Children 	[]Obj	`json:"children"`
	Username	string
	Upvotes		int
	Downvotes	int
	Votes		int
	Posted		time.Time
	Parent		string	
}

var (
	nodeTable = map[string]*Node{}
	root      *Node
)

func add(id, name, parentId string) {
	fmt.Printf("add: id=%v name=%v parentId=%v\n", id, name, parentId)

	node := &Node{name: name, children: []*Node{}}

	if parentId == "0" {
		root = node
	} else {

		parent, ok := nodeTable[parentId]
		if !ok {
			fmt.Printf("add: parentId=%v: not found\n", parentId)
			return
		}

		parent.children = append(parent.children, node)
	}

	nodeTable[id] = node
}


func scan(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error reading lines: %v\n", err)
			return
		}
		tokens := strings.Fields(line)
		if t := len(tokens); t != 3 {
			fmt.Printf("bad input line %v: tokens=%d [%v]\n", lineCount, t, line)
			continue
		}
		add(tokens[0], tokens[1], tokens[2])
	}
}



func showNode(node *Node, prefix string, currentObjs []Obj) []Obj{
	if node == nil {
		var objs []Obj
		return objs
	}

	var obj Obj
	obj.Data = node.name
	
	var objs []Obj
	for _, n := range node.children {		
		for _, nn := range showNode(n, prefix + "--", objs){
			obj.Children = append(obj.Children, nn)		
		}
	}
	currentObjs = append(currentObjs, obj)
	return currentObjs
}

func show() []Obj{
	if root == nil {
		fmt.Printf("show: root node not found\n")
		return nil
	}
	fmt.Printf("RESULT:\n")
	var objs []Obj
	return showNode(root, "", objs)
}

func p(str interface{}) {
	fmt.Println(str)
}

func main() {
	p("Hello. Scanning...")
	scan(os.Args[1])
	fmt.Printf("main: reading input from stdin -- done\n")
	tree := show()
	p(tree)

	fmt.Printf("main: end\n")

	data, err := json.Marshal(tree)
	if err != nil {
		log.Fatal("Marshalling problem" + err.Error())
	}
	p(string(data))
}
