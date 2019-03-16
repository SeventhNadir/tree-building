package tree

import (
	"errors"
	"sort"
)

//Record of post ID's and their parents
type Record struct {
	ID     int
	Parent int
}

//Node representation, recreated from Records
type Node struct {
	ID       int
	Children []*Node
}

//Build a tree of nodes from a set of related records.
func Build(records []Record) (*Node, error) {
	if len(records) == 0 {
		return nil, nil
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })

	err := errorCheck(records)
	if err != nil {
		return nil, err
	}

	parentChild := getParentChildMap(records)
	root := &Node{ID: records[0].ID}
	root.Children = getChildren(root, parentChild)
	return root, nil
}

func getParentChildMap(records []Record) map[int][]*Node {
	childMap := make(map[int][]*Node)
	for _, record := range records {
		parentID := record.Parent
		childNode := &Node{record.ID, nil}
		childMap[parentID] = append(childMap[parentID], childNode)
	}
	return childMap
}

func getChildren(parent *Node, parentChild map[int][]*Node) []*Node {
	var nodeArray []*Node
	for _, child := range parentChild[parent.ID] {
		if child.ID != parent.ID {
			node := &Node{ID: child.ID}
			node.Children = getChildren(node, parentChild)
			nodeArray = append(nodeArray, node)
		}
	}
	return nodeArray
}

func errorCheck(records []Record) error {
	for i, record := range records {
		if i != record.ID {
			return errors.New("non-contiguous")
		}
		if record.ID < record.Parent {
			return errors.New("higher id parent of lower id")
		}
		originalRecord := record
		for record.ID != 0 {
			record = records[record.Parent]
			if record == originalRecord {
				return errors.New("cycle indirectly")
			}
		}
	}
	return nil
}
