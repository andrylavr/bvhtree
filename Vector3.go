package bvhtree

// Vector3 is a 3D Vector class
type Vector3 struct {
	X, Y, Z float64
}

// Point represents a point in 3D space
type Point = *Vector3

func NewPoint(arr ...float64) Point {
	p := &Vector3{}
	switch len(arr) {
	case 1:
		p.X = arr[0]
	case 2:
		p.X, p.Y = arr[0], arr[1]
	case 3:
		p.X, p.Y, p.Z = arr[0], arr[1], arr[2]
	}
	return p
}

// Copy copies the values from another vector
func (v *Vector3) Copy(src *Vector3) *Vector3 {
	v.X = src.X
	v.Y = src.Y
	v.Z = src.Z
	return v
}

// Set sets the vector's components
func (v *Vector3) Set(x, y, z float64) *Vector3 {
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// SetFromArray sets the vector's components from an array
func (v *Vector3) SetFromArray(array []float64, firstElementPos int) {
	v.X = array[firstElementPos]
	v.Y = array[firstElementPos+1]
	v.Z = array[firstElementPos+2]
}

// Add adds another vector to this vector
func (v *Vector3) Add(src *Vector3) *Vector3 {
	v.X += src.X
	v.Y += src.Y
	v.Z += src.Z
	return v
}

// MultiplyScalar multiplies the vector's components by a scalar
func (v *Vector3) MultiplyScalar(scalar float64) *Vector3 {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	return v
}

// SubVectors sets this vector to be the difference of two vectors
func (v *Vector3) SubVectors(a, b *Vector3) *Vector3 {
	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z
	return v
}

// Dot returns the dot product of this vector and another vector
func (v *Vector3) Dot(src *Vector3) float64 {
	return v.X*src.X + v.Y*src.Y + v.Z*src.Z
}

// Cross sets this vector to be the cross product of this vector and another vector
func (v *Vector3) Cross(src *Vector3) *Vector3 {
	x := v.X
	y := v.Y
	z := v.Z
	v.X = y*src.Z - z*src.Y
	v.Y = z*src.X - x*src.Z
	v.Z = x*src.Y - y*src.X
	return v
}

// CrossVectors sets this vector to be the cross product of two other vectors
func (v *Vector3) CrossVectors(a, b *Vector3) *Vector3 {
	var ax, ay, az = a.X, a.Y, a.Z
	var bx, by, bz = b.X, b.Y, b.Z

	v.X = ay*bz - az*by
	v.Y = az*bx - ax*bz
	v.Z = ax*by - ay*bx

	return v
}

// Clone creates a new vector with the same components as this vector
func (v *Vector3) Clone() *Vector3 {
	return &Vector3{v.X, v.Y, v.Z}
}
