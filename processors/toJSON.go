package main

import (
	"encoding/json"
	"fmt"
)

type Node struct {
	Id       string  `json:"-"`
	ParentId string  `json:"-"`
	Name     string  `json:"name"`
	Leaf     string  `json:"leaf,omitempty"`
	Children []*Node `json:"children,omitempty"`
}
func (this *Node) Size() int {
	var size int = len(this.Children)
	for _, c := range this.Children {
		size += c.Size()
	}
	return size
}
func (this *Node) Add(nodes... *Node) bool {
	var size = this.Size();
	for _, n := range nodes {
		if n.ParentId == this.Id {
			this.Children = append(this.Children, n)
		} else { 
			for _, c := range this.Children {
				if c.Add(n) {
					break
				}
			}
		}
	}
	return this.Size() == size + len(nodes)
}

func main() {
	var root *Node = &Node{"001", "", "root", "", nil}
	data := []*Node{
		&Node{"002", "001", "Shooping","0", nil},
		&Node{"003", "002", "Housewares","0", nil},
		&Node{"004", "003", "Kitchen","1", nil},
		&Node{"005", "003", "Officer","1", nil},
		&Node{"006", "002", "Remodeling","0", nil},
		&Node{"007", "006", "Retile kitchen","1", nil},
		&Node{"008", "006", "Paint bedroom","1", nil},
		&Node{"009", "008", "Ceiling","1", nil},
		&Node{"010", "006", "Other","1", nil},
		&Node{"011", "001", "Misc","1", nil},
	}
	fmt.Println(root.Add(data...), root.Size())
	bytes, _ := json.MarshalIndent(root, "", "\t") //formated output
	//bytes, _ := json.Marshal(root)
	fmt.Println(string(bytes))
}
