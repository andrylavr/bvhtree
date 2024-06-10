package main

import "github.com/andrylavr/bvhtree"
import "fmt"

func main() {
	triangle0 := bvhtree.NewTriangle(
		0.0, 0.0, 0.0,
		1000.0, 0.0, 0.0,
		1000.0, 1000.0, 0.0,
	)

	triangle1 := bvhtree.NewTriangle(
		0.0, 0.0, 0.0,
		2000.0, 0.0, 0.0,
		2000.0, 1000.0, 0.0,
	)

	maxTrianglesPerNode := 7

	triangles := []bvhtree.Triangle{triangle0, triangle1}
	bvh := bvhtree.NewBVH(triangles, maxTrianglesPerNode)

	fmt.Println("BVH Tree created successfully!")

	// Define ray parameters
	rayOrigin := bvhtree.NewPoint(1500.0, 3.0, 1000.0)
	rayDirection := bvhtree.NewPoint(0.0, 0.0, -1.0) // Should be normalized
	backfaceCulling := true

	// Perform ray-triangle intersection
	intersectionResult := bvh.IntersectRay(rayOrigin, rayDirection, backfaceCulling)

	// Print intersection results
	fmt.Println("Intersection Results:")
	for _, result := range intersectionResult {
		fmt.Printf("Triangle Index: %d\n", result.TriangleIndex)
		fmt.Printf("Intersection Point: (%f, %f, %f)\n", result.IntersectionPoint.X, result.IntersectionPoint.Y, result.IntersectionPoint.Z)
	}
}
