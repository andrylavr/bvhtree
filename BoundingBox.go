package bvhtree

// SetBox sets the bounding box information in the bboxArray at the specified position
func SetBox(bboxArray []float64, pos int, triangleID int, minX, minY, minZ, maxX, maxY, maxZ float64) {
	bboxArray[pos*7] = float64(triangleID)
	bboxArray[pos*7+1] = minX
	bboxArray[pos*7+2] = minY
	bboxArray[pos*7+3] = minZ
	bboxArray[pos*7+4] = maxX
	bboxArray[pos*7+5] = maxY
	bboxArray[pos*7+6] = maxZ
}

// CopyBox copies the bounding box information from sourceArray at sourcePos to destArray at destPos
func CopyBox(sourceArray []float64, sourcePos int, destArray []float64, destPos int) {
	destArray[destPos*7] = sourceArray[sourcePos*7]
	destArray[destPos*7+1] = sourceArray[sourcePos*7+1]
	destArray[destPos*7+2] = sourceArray[sourcePos*7+2]
	destArray[destPos*7+3] = sourceArray[sourcePos*7+3]
	destArray[destPos*7+4] = sourceArray[sourcePos*7+4]
	destArray[destPos*7+5] = sourceArray[sourcePos*7+5]
	destArray[destPos*7+6] = sourceArray[sourcePos*7+6]
}

// GetBox retrieves the bounding box information from bboxArray at the specified position and stores it in outputBox
func GetBox(bboxArray []float64, pos int, outputBox *BoundingBox) {
	outputBox.TriangleID = int(bboxArray[pos*7])
	outputBox.MinX = bboxArray[pos*7+1]
	outputBox.MinY = bboxArray[pos*7+2]
	outputBox.MinZ = bboxArray[pos*7+3]
	outputBox.MaxX = bboxArray[pos*7+4]
	outputBox.MaxY = bboxArray[pos*7+5]
	outputBox.MaxZ = bboxArray[pos*7+6]
}

// BoundingBox represents a bounding box in 3D space
type BoundingBox struct {
	TriangleID int     // ID of the triangle associated with the bounding box
	MinX       float64 // Minimum X coordinate
	MinY       float64 // Minimum Y coordinate
	MinZ       float64 // Minimum Z coordinate
	MaxX       float64 // Maximum X coordinate
	MaxY       float64 // Maximum Y coordinate
	MaxZ       float64 // Maximum Z coordinate
}
