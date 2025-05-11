package main

type (
	nodeType int

	node interface {
		Id() id
		Entity() entity
		NodeType() nodeType
	}
)

const (
	stationNode nodeType = iota
	intersectionNode
)

// Implementors of node

func (s station) Id() id             { return s.id }
func (s station) Entity() entity     { return s.entity }
func (s station) NodeType() nodeType { return stationNode }

func (s intersection) Id() id              { return s.id }
func (in intersection) Entity() entity     { return in.entity }
func (in intersection) NodeType() nodeType { return intersectionNode }
