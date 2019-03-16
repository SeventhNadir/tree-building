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

	root, err := getRoot(records)
	if err != nil {
		return nil, err
	}

	children := getChildren(root.ID, records)
	node := Node{root.ID, nil}

	if children != nil {
		node = Node{root.ID, children}
	}

	return &node, nil
}

func buildNode(record Record) *Node {
	node := Node{ID: record.ID}
	return &node
}

func getRoot(records []Record) (*Node, error) {
	for _, record := range records {
		if record.ID == 0 {
			if record.Parent != 0 {
				return nil, errors.New("root node has parent")
			}
			node := buildNode(record)
			return node, nil
		}
	}
	return &Node{}, errors.New("no root node")
}

func getChildren(parent int, allRecords []Record) []*Node {
	var nodeArray []*Node
	for _, child := range allRecords {
		if child.Parent == parent && child.ID != parent {
			node := buildNode(child)
			node.Children = getChildren(node.ID, allRecords)
			nodeArray = append(nodeArray, node)
		}
	}
	if len(nodeArray) == 0 {
		return nil
	}

	return nodeArray
}

func errorCheck(records []Record) error {
	duplicate := map[int]bool{}
	for _, record := range records {
		if duplicate[record.ID] == true {
			if record.ID == 0 {
				return errors.New("duplicate node")
			}
			return errors.New("duplicate root")
		}
		duplicate[record.ID] = true
	}

	for i, record := range records {
		if i != record.ID {
			return errors.New("non-contiguous")
		}
		if record.ID < record.Parent {
			return errors.New("higher id parent of lower id")
		}
	}

	for _, record := range records {
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
