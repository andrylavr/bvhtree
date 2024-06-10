package bvhtree

import "math"

// Node represents a node in the BVH structure
type Node struct {
	ExtentsMin, ExtentsMax Point
	StartIndex, EndIndex   int
	Level                  int
	Node0, Node1           *Node
}

// NewBVHNode creates a new BVH node
func NewBVHNode(extentsMin, extentsMax Point, startIndex, endIndex, level int) *Node {
	return &Node{
		ExtentsMin: extentsMin,
		ExtentsMax: extentsMax,
		StartIndex: startIndex,
		EndIndex:   endIndex,
		Level:      level,
		Node0:      nil,
		Node1:      nil,
	}
}

// ElementCount returns the number of elements in the node
func (node *Node) ElementCount() int {
	return node.EndIndex - node.StartIndex
}

// CenterX returns the center X coordinate of the node
func (node *Node) CenterX() float64 {
	return (node.ExtentsMin.X + node.ExtentsMax.X) * 0.5
}

// CenterY returns the center Y coordinate of the node
func (node *Node) CenterY() float64 {
	return (node.ExtentsMin.Y + node.ExtentsMax.Y) * 0.5
}

// CenterZ returns the center Z coordinate of the node
func (node *Node) CenterZ() float64 {
	return (node.ExtentsMin.Z + node.ExtentsMax.Z) * 0.5
}

// ClearShapes clears the shapes in the node
func (node *Node) ClearShapes() {
	node.StartIndex = -1
	node.EndIndex = -1
}

// CalcBoundingSphereRadius calculates the radius of the bounding sphere for the node
func CalcBoundingSphereRadius(extentsMin, extentsMax Point) float64 {
	centerX := (extentsMin.X + extentsMax.X) * 0.5
	centerY := (extentsMin.Y + extentsMax.Y) * 0.5
	centerZ := (extentsMin.Z + extentsMax.Z) * 0.5

	extentsMinDistSqr :=
		(centerX-extentsMin.X)*(centerX-extentsMin.X) +
			(centerY-extentsMin.Y)*(centerY-extentsMin.Y) +
			(centerZ-extentsMin.Z)*(centerZ-extentsMin.Z)

	extentsMaxDistSqr :=
		(centerX-extentsMax.X)*(centerX-extentsMax.X) +
			(centerY-extentsMax.Y)*(centerY-extentsMax.Y) +
			(centerZ-extentsMax.Z)*(centerZ-extentsMax.Z)

	return math.Sqrt(math.Max(extentsMinDistSqr, extentsMaxDistSqr))
}
