package main

import (
	"github.com/andrylavr/bvhtree"
	"github.com/andrylavr/wasmo"
	"syscall/js"
)

var window, Object, Array js.Value

func pointToJS(point bvhtree.Point) js.Value {
	pointJS := Object.New()
	pointJS.Set("x", point.X)
	pointJS.Set("y", point.Y)
	pointJS.Set("z", point.Z)
	return pointJS
}

func jsToPoint(pointJS js.Value) bvhtree.Point {
	point := bvhtree.NewPoint(
		pointJS.Get("x").Float(),
		pointJS.Get("y").Float(),
		pointJS.Get("z").Float(),
	)
	return point
}

func triangleToJS(triangle bvhtree.Triangle) js.Value {
	triangleJS := Array.New()
	for _, point := range triangle {
		pointJS := pointToJS(point)
		triangleJS.Call("push", pointJS)
	}
	return triangleJS
}

func printPoint(name string, p bvhtree.Point) {
	//fmt.Printf("%s X: %d; Y: %d; Z: %d\n", name, p.X, p.Y, p.Z)
}

func intersectRayJS(this js.Value, args []js.Value) interface{} {
	intersections := Array.New()
	bvh := wasmo.GetLinkedVar(this, "bvh").(*bvhtree.BVH)

	////fmt.Println("len bvh.VertexArray()", len(bvh.VertexArray()))

	rayOrigin := jsToPoint(args[0])
	rayDirection := jsToPoint(args[1])
	backfaceCulling := args[2].Bool()

	//fmt.Println("intersectRayJS:")
	printPoint("rayOrigin", rayOrigin)
	printPoint("rayDirection", rayDirection)
	//fmt.Println("backfaceCulling", backfaceCulling)

	intersectionResults := bvh.IntersectRay(rayOrigin, rayDirection, backfaceCulling)
	//fmt.Println("intersectionResults len", len(intersectionResults))
	for _, intersectionResult := range intersectionResults {
		intersectionJS := Object.New()
		intersectionJS.Set("triangle", triangleToJS(intersectionResult.Triangle))
		intersectionJS.Set("triangleIndex", intersectionResult.TriangleIndex)
		intersectionJS.Set("intersectionPoint", pointToJS(intersectionResult.IntersectionPoint))
		intersections.Call("push", intersectionJS)
	}
	return intersections
}

func bvhFromVertexArrayJS(this js.Value, args []js.Value) interface{} {
	bvhJS := Object.New()

	vertexArrayJS := args[0]
	vertexArrayLen := vertexArrayJS.Get("length").Int()
	vertexArray := make([]float64, vertexArrayLen)
	for i := 0; i < vertexArrayLen; i++ {
		vertexArray[i] = vertexArrayJS.Index(i).Float()
	}

	maxTrianglesPerNode := args[1].Int()
	bvh := bvhtree.NewBVHFromVertexArray(vertexArray, maxTrianglesPerNode)
	wasmo.LinkVar(bvhJS, "bvh", bvh)
	bvhJS.Set("intersectRay", js.FuncOf(intersectRayJS))
	return bvhJS
}

type Tezt struct {
	A int
}

func main() {
	window = js.Global()
	Object = window.Get("Object")
	Array = window.Get("Array")

	//teztHolder := Object.New()
	//window.Set("teztHolder", teztHolder)
	//tezt := &Tezt{}
	//tezt.A = 456
	//linkVar(teztHolder, tezt)
	//tezt2 := getVar[Tezt](teztHolder)
	////fmt.Println("tezt2.A", tezt2.A)

	bvhTreeJS := Object.New()
	bvhTreeJS.Set("bvhFromVertexArray", js.FuncOf(bvhFromVertexArrayJS))
	window.Set("bvhtree", bvhTreeJS)

	select {}
}
