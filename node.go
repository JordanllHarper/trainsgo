package main

type (
	nodeType int

	Node interface {
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

func (s Station) Id() id             { return s.E.Id }
func (s Station) Entity() entity     { return s.E }
func (s Station) NodeType() nodeType { return stationNode }

func (s intersection) Id() id              { return s.entity.Id }
func (in intersection) Entity() entity     { return in.entity }
func (in intersection) NodeType() nodeType { return intersectionNode }
