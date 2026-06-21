package geom

import "math"

// Vector3 represents a point or direction in 3D space.
type Vector3 struct {
	X, Y, Z float64
}

// Add returns the sum of two vectors.
func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Sub returns the difference between two vectors.
func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Scale multiplies the vector by a scalar.
func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{v.X * s, v.Y * s, v.Z * s}
}

// Dot returns the dot product of two vectors.
func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of two vectors.
func (v Vector3) Cross(other Vector3) Vector3 {
	return Vector3{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}

// Length returns the magnitude of the vector.
func (v Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a unit vector in the same direction.
func (v Vector3) Normalize() Vector3 {
	l := v.Length()
	if l == 0 {
		return Vector3{}
	}
	return v.Scale(1.0 / l)
}

// Transform applies a 4x4 matrix transformation (as a position, w=1).
// Matrix is column-major: m[col][row].
// The w component uses row 3 of each column.
func (v Vector3) Transform(m Matrix4) Vector3 {
	x := m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0]
	y := m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1]
	z := m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2]
	w := m[0][3]*v.X + m[1][3]*v.Y + m[2][3]*v.Z + m[3][3]

	if w == 0 {
		w = 1
	}
	return Vector3{
		X: x / w,
		Y: y / w,
		Z: z / w,
	}
}

// TransformDirection applies a 4x4 matrix transformation (as a direction, w=0).
func (v Vector3) TransformDirection(m Matrix4) Vector3 {
	return Vector3{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z,
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z,
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z,
	}
}
