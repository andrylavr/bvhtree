package bvhtree

import (
	"math"
	"sort"
)

const EPSILON = 1e-6

// Triangle represents a triangle in 3D space
type Triangle [3]Point

func NewTriangle(f0, f1, f2, f3, f4, f5, f6, f7, f8 float64) Triangle {
	t := Triangle{
		NewPoint(f0, f1, f2),
		NewPoint(f3, f4, f5),
		NewPoint(f6, f7, f8),
	}
	return t
}

// IntersectionResult represents the result of a ray-triangle intersection
type IntersectionResult struct {
	Triangle          Triangle
	TriangleIndex     int
	IntersectionPoint Point
}

// BVH represents a bounding volume hierarchy
type BVH struct {
	vertexArray         []float64
	maxTrianglesPerNode int
	bboxArray           []float64
	bboxHelper          []float64
	rootNode            *Node
	nodesToSplit        []*Node
}

// NewBVH creates a new BVH from a list of triangles
func NewBVHFromVertexArray(vertexArray []float64, maxTrianglesPerNode int) *BVH {
	bvh := &BVH{
		vertexArray:         vertexArray,
		maxTrianglesPerNode: maxTrianglesPerNode,
	}

	bvh.bboxArray = bvh.CalcBoundingBoxes(vertexArray)
	bvh.bboxHelper = make([]float64, len(bvh.bboxArray))
	copy(bvh.bboxHelper, bvh.bboxArray)

	triangleCount := len(vertexArray) / 9
	extents := bvh.CalcExtents(0, triangleCount, EPSILON)
	bvh.rootNode = NewBVHNode(extents[0], extents[1], 0, triangleCount, 0)
	bvh.nodesToSplit = []*Node{bvh.rootNode}

	for len(bvh.nodesToSplit) > 0 {
		node := bvh.nodesToSplit[len(bvh.nodesToSplit)-1]
		bvh.nodesToSplit = bvh.nodesToSplit[:len(bvh.nodesToSplit)-1]
		bvh.SplitNode(node)
	}

	return bvh
}

func NewBVH(triangles []Triangle, maxTrianglesPerNode int) *BVH {
	vertexArray := make([]float64, len(triangles)*9)
	for i, tri := range triangles {
		vertexArray[i*9] = tri[0].X
		vertexArray[i*9+1] = tri[0].Y
		vertexArray[i*9+2] = tri[0].Z

		vertexArray[i*9+3] = tri[1].X
		vertexArray[i*9+4] = tri[1].Y
		vertexArray[i*9+5] = tri[1].Z

		vertexArray[i*9+6] = tri[2].X
		vertexArray[i*9+7] = tri[2].Y
		vertexArray[i*9+8] = tri[2].Z
	}

	return NewBVHFromVertexArray(vertexArray, maxTrianglesPerNode)
}

func (bvh *BVH) VertexArray() []float64 {
	return bvh.vertexArray
}

// IntersectRay returns a list of all the triangles in the BVH which intersected a specific ray
func (bvh *BVH) IntersectRay(rayOrigin, rayDirection Point, backfaceCulling bool) []IntersectionResult {
	nodesToIntersect := []*Node{bvh.rootNode}
	var trianglesInIntersectingNodes []int
	var intersectingTriangles []IntersectionResult

	invRayDirection := &Vector3{
		X: 1.0 / rayDirection.X,
		Y: 1.0 / rayDirection.Y,
		Z: 1.0 / rayDirection.Z,
	}

	for len(nodesToIntersect) > 0 {
		node := nodesToIntersect[len(nodesToIntersect)-1]
		nodesToIntersect = nodesToIntersect[:len(nodesToIntersect)-1]

		if IntersectNodeBox(rayOrigin, invRayDirection, node) {
			if node.Node0 != nil {
				nodesToIntersect = append(nodesToIntersect, node.Node0)
			}
			if node.Node1 != nil {
				nodesToIntersect = append(nodesToIntersect, node.Node1)
			}
			for i := node.StartIndex; i < node.EndIndex; i++ {
				trianglesInIntersectingNodes = append(trianglesInIntersectingNodes, int(bvh.bboxArray[i*7]))
			}
		}
	}

	a := &Vector3{}
	b := &Vector3{}
	c := &Vector3{}
	rayOriginVec3 := &Vector3{X: rayOrigin.X, Y: rayOrigin.Y, Z: rayOrigin.Z}
	rayDirectionVec3 := &Vector3{X: rayDirection.X, Y: rayDirection.Y, Z: rayDirection.Z}

	for _, triIndex := range trianglesInIntersectingNodes {
		a.SetFromArray(bvh.vertexArray, triIndex*9)
		b.SetFromArray(bvh.vertexArray, triIndex*9+3)
		c.SetFromArray(bvh.vertexArray, triIndex*9+6)

		if intersectionPoint := IntersectRayTriangle(a, b, c, rayOriginVec3, rayDirectionVec3, backfaceCulling); intersectionPoint != nil {
			intersectingTriangles = append(intersectingTriangles, IntersectionResult{
				Triangle:          Triangle{a, b, c},
				TriangleIndex:     triIndex,
				IntersectionPoint: intersectionPoint,
			})
		}
	}

	return intersectingTriangles
}

// CalcBoundingBoxes calculates the bounding box for each triangle and adds an index to the triangle's position in the array.
// Each bbox is saved as 7 values in a float64 slice: (position, minX, minY, minZ, maxX, maxY, maxZ).
func (bvh *BVH) CalcBoundingBoxes(vertexArray []float64) []float64 {
	triangleCount := len(vertexArray) / 9
	bboxArray := make([]float64, triangleCount*7)

	for i := 0; i < triangleCount; i++ {
		p0x := vertexArray[i*9]
		p0y := vertexArray[i*9+1]
		p0z := vertexArray[i*9+2]
		p1x := vertexArray[i*9+3]
		p1y := vertexArray[i*9+4]
		p1z := vertexArray[i*9+5]
		p2x := vertexArray[i*9+6]
		p2y := vertexArray[i*9+7]
		p2z := vertexArray[i*9+8]

		minX := math.Min(math.Min(p0x, p1x), p2x)
		minY := math.Min(math.Min(p0y, p1y), p2y)
		minZ := math.Min(math.Min(p0z, p1z), p2z)
		maxX := math.Max(math.Max(p0x, p1x), p2x)
		maxY := math.Max(math.Max(p0y, p1y), p2y)
		maxZ := math.Max(math.Max(p0z, p1z), p2z)

		SetBox(bboxArray, i, i, minX, minY, minZ, maxX, maxY, maxZ)
	}

	return bboxArray
}

// CalcExtents calculates the extents (i.e., the min and max coordinates) of a list of bounding boxes in the bboxArray.
// It takes startIndex, endIndex, and expandBy as parameters to define the range and the safety margin.
func (bvh *BVH) CalcExtents(startIndex, endIndex int, expandBy float64) [2]Point {
	if startIndex >= endIndex {
		return [2]Point{{0, 0, 0}, {0, 0, 0}}
	}

	minX := math.MaxFloat64
	minY := math.MaxFloat64
	minZ := math.MaxFloat64
	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64
	maxZ := -math.MaxFloat64

	for i := startIndex; i < endIndex; i++ {
		minX = math.Min(bvh.bboxArray[i*7+1], minX)
		minY = math.Min(bvh.bboxArray[i*7+2], minY)
		minZ = math.Min(bvh.bboxArray[i*7+3], minZ)
		maxX = math.Max(bvh.bboxArray[i*7+4], maxX)
		maxY = math.Max(bvh.bboxArray[i*7+5], maxY)
		maxZ = math.Max(bvh.bboxArray[i*7+6], maxZ)
	}

	return [2]Point{
		{minX - expandBy, minY - expandBy, minZ - expandBy},
		{maxX + expandBy, maxY + expandBy, maxZ + expandBy},
	}
}

// SplitNode splits a node if it contains more elements than maxTrianglesPerNode
func (bvh *BVH) SplitNode(node *Node) {
	if node.ElementCount() <= bvh.maxTrianglesPerNode || node.ElementCount() == 0 {
		return
	}

	startIndex := node.StartIndex
	endIndex := node.EndIndex

	leftNode := [3][]int{}
	rightNode := [3][]int{}
	extentCenters := [3]float64{node.CenterX(), node.CenterY(), node.CenterZ()}

	extentsLength := [3]float64{
		node.ExtentsMax.X - node.ExtentsMin.X,
		node.ExtentsMax.Y - node.ExtentsMin.Y,
		node.ExtentsMax.Z - node.ExtentsMin.Z,
	}

	objectCenter := [3]float64{}

	for i := startIndex; i < endIndex; i++ {
		objectCenter[0] = (bvh.bboxArray[i*7+1] + bvh.bboxArray[i*7+4]) * 0.5
		objectCenter[1] = (bvh.bboxArray[i*7+2] + bvh.bboxArray[i*7+5]) * 0.5
		objectCenter[2] = (bvh.bboxArray[i*7+3] + bvh.bboxArray[i*7+6]) * 0.5

		for j := 0; j < 3; j++ {
			if objectCenter[j] < extentCenters[j] {
				leftNode[j] = append(leftNode[j], i)
			} else {
				rightNode[j] = append(rightNode[j], i)
			}
		}
	}

	splitFailed := [3]bool{
		len(leftNode[0]) == 0 || len(rightNode[0]) == 0,
		len(leftNode[1]) == 0 || len(rightNode[1]) == 0,
		len(leftNode[2]) == 0 || len(rightNode[2]) == 0,
	}

	if splitFailed[0] && splitFailed[1] && splitFailed[2] {
		return
	}

	splitOrder := []int{0, 1, 2}
	sort.Slice(splitOrder, func(i, j int) bool {
		return extentsLength[splitOrder[j]] > extentsLength[splitOrder[i]]
	})

	var leftElements, rightElements []int

	for _, candidateIndex := range splitOrder {
		if !splitFailed[candidateIndex] {
			leftElements = leftNode[candidateIndex]
			rightElements = rightNode[candidateIndex]
			break
		}
	}

	node0Start := startIndex
	node0End := node0Start + len(leftElements)
	node1Start := node0End
	node1End := endIndex

	helperPos := node.StartIndex
	concatenatedElements := append(leftElements, rightElements...)

	for _, currElement := range concatenatedElements {
		CopyBox(bvh.bboxArray, currElement, bvh.bboxHelper, helperPos)
		helperPos++
	}

	subArr := bvh.bboxHelper[node.StartIndex*7 : node.EndIndex*7]
	copy(bvh.bboxArray[node.StartIndex*7:], subArr)

	node0Extents := bvh.CalcExtents(node0Start, node0End, EPSILON)
	node1Extents := bvh.CalcExtents(node1Start, node1End, EPSILON)

	node0 := NewBVHNode(node0Extents[0], node0Extents[1], node0Start, node0End, node.Level+1)
	node1 := NewBVHNode(node1Extents[0], node1Extents[1], node1Start, node1End, node.Level+1)

	node.Node0 = node0
	node.Node1 = node1
	node.ClearShapes()

	bvh.nodesToSplit = append(bvh.nodesToSplit, node0, node1)
}

// TValues represents the tmin and tmax values for ray-box intersection
type TValues struct {
	Min, Max float64
}

// CalcTValues calculates the tmin and tmax values for ray-box intersection
func CalcTValues(minVal, maxVal, rayOriginCoord, invdir float64) TValues {
	res := TValues{}
	if invdir >= 0 {
		res.Min = (minVal - rayOriginCoord) * invdir
		res.Max = (maxVal - rayOriginCoord) * invdir
	} else {
		res.Min = (maxVal - rayOriginCoord) * invdir
		res.Max = (minVal - rayOriginCoord) * invdir
	}
	return res
}

func isNaN(f float64) bool {
	return math.IsInf(f, 1) || math.IsInf(f, -1)
}

// IntersectNodeBox checks if a ray intersects with a node's bounding box
func IntersectNodeBox(rayOrigin, invRayDirection Point, node *Node) bool {
	t := CalcTValues(node.ExtentsMin.X, node.ExtentsMax.X, rayOrigin.X, invRayDirection.X)
	ty := CalcTValues(node.ExtentsMin.Y, node.ExtentsMax.Y, rayOrigin.Y, invRayDirection.Y)

	if t.Min > ty.Max || ty.Min > t.Max {
		return false
	}

	if ty.Min > t.Min || isNaN(t.Min) {
		t.Min = ty.Min
	}
	if ty.Max < t.Max || isNaN(t.Max) {
		t.Max = ty.Max
	}

	tz := CalcTValues(node.ExtentsMin.Z, node.ExtentsMax.Z, rayOrigin.Z, invRayDirection.Z)

	if t.Min > tz.Max || tz.Min > t.Max {
		return false
	}
	if tz.Min > t.Min || isNaN(t.Min) {
		t.Min = tz.Min
	}
	if tz.Max < t.Max || isNaN(t.Max) {
		t.Max = tz.Max
	}

	if t.Max < 0 {
		return false
	}

	return true
}

// IntersectRayTriangle determines if a ray intersects with a triangle in 3D space
func IntersectRayTriangle(a, b, c, rayOrigin, rayDirection Point, backfaceCulling bool) Point {
	// Compute the offset origin, edges, and normal.
	var diff = &Vector3{}
	var edge1 = &Vector3{}
	var edge2 = &Vector3{}
	var normal = &Vector3{}

	// Compute the offset origin, edges, and normal.
	edge1.SubVectors(b, a)
	edge2.SubVectors(c, a)
	normal.CrossVectors(edge1, edge2)

	// Solve Q + t*D = b1*E1 + bL*E2 (Q = kDiff, D = ray direction,
	// E1 = kEdge1, E2 = kEdge2, N = Cross(E1,E2)) by
	//   |Dot(D,N)|*b1 = sign(Dot(D,N))*Dot(D,Cross(Q,E2))
	//   |Dot(D,N)|*b2 = sign(Dot(D,N))*Dot(D,Cross(E1,Q))
	//   |Dot(D,N)|*t = -sign(Dot(D,N))*Dot(Q,N)
	DdN := rayDirection.Dot(normal)
	var sign float64

	if DdN > 0 {
		if backfaceCulling {
			return nil
		}
		sign = 1
	} else if DdN < 0 {
		sign = -1
		DdN = -DdN
	} else {
		return nil
	}

	diff.SubVectors(rayOrigin, a)
	DdQxE2 := sign * rayDirection.Dot(edge2.CrossVectors(diff, edge2))

	// b1 < 0, no intersection
	if DdQxE2 < 0 {
		return nil
	}

	DdE1xQ := sign * rayDirection.Dot(edge1.Cross(diff))

	// b2 < 0, no intersection
	if DdE1xQ < 0 {
		return nil
	}

	// b1+b2 > 1, no intersection
	if DdQxE2+DdE1xQ > DdN {
		return nil
	}

	// Line intersects triangle, check if ray does.
	QdN := -sign * diff.Dot(normal)

	// t < 0, no intersection
	if QdN < 0 {
		return nil
	}

	// Ray intersects triangle.
	t := QdN / DdN
	result := &Vector3{0, 0, 0}
	result.X = rayDirection.X*t + rayOrigin.X
	result.Y = rayDirection.Y*t + rayOrigin.Y
	result.Z = rayDirection.Z*t + rayOrigin.Z
	return result
}
