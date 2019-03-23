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
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})

	for i, record := range records {
		if i != record.ID {
			return nil, errors.New("non-contiguous")
		}
		if record.ID < record.Parent {
			return nil, errors.New("higher id parent of lower id")
		}
		originalRecord := record
		for record.ID != 0 {
			record = records[record.Parent]
			if record == originalRecord {
				return nil, errors.New("cycle indirectly")
			}
		}
	}

	nodes := make([]*Node, len(records))
	records = records[1:]

	nodes[0] = &Node{ID: 0}
	for _, record := range records {
		parent := nodes[record.Parent]
		node := &Node{ID: record.ID}
		nodes[node.ID] = node
		parent.Children = append(parent.Children, node)
	}
	return nodes[0], nil

}
