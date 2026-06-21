package geom

import "math"

// Matrix4 represents a 4x4 matrix for 3D transformations and projections.
// Stored in column-major order: m[col][row]
type Matrix4 [4][4]float64

// Identity returns the identity matrix.
func Identity() Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Multiply returns the product of two 4x4 matrices (this * n).
func (m Matrix4) Multiply(n Matrix4) Matrix4 {
	var result Matrix4
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += m[k][row] * n[col][k]
			}
			result[col][row] = sum
		}
	}
	return result
}

// Translate returns a translation matrix by (x, y, z).
// Column-major: translation goes in column 3, rows 0-2.
func Translate(x, y, z float64) Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{x, y, z, 1},
	}
}

// ScaleMatrix returns a scaling matrix.
func ScaleMatrix(sx, sy, sz float64) Matrix4 {
	return Matrix4{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	}
}

// RotateX returns a rotation matrix around the X axis (radians).
// In column-major: col 1 and col 2 get the sin/cos terms.
func RotateX(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		{1, 0, 0, 0},
		{0, c, s, 0},
		{0, -s, c, 0},
		{0, 0, 0, 1},
	}
}

// RotateY returns a rotation matrix around the Y axis (radians).
func RotateY(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		{c, 0, -s, 0},
		{0, 1, 0, 0},
		{s, 0, c, 0},
		{0, 0, 0, 1},
	}
}

// RotateZ returns a rotation matrix around the Z axis (radians).
func RotateZ(angle float64) Matrix4 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix4{
		{c, s, 0, 0},
		{-s, c, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

// Perspective returns a perspective projection matrix.
// fovY: field of view in radians
// aspect: width / height
// near, far: clipping planes
// Perspective returns a perspective projection matrix.
func Perspective(fovY, aspect, near, far float64) Matrix4 {
	f := 1.0 / math.Tan(fovY/2.0)

	return Matrix4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, -(far + near) / (far - near), -1},
		{0, 0, -2 * far * near / (far - near), 0},
	}
}

// LookAt returns a view matrix that positions the camera at eye, looking at center, with up as the upward direction.
// In OpenGL convention, the camera looks down the -Z axis.
func LookAt(eye, center, up Vector3) Matrix4 {
	f := center.Sub(eye).Normalize()         // forward (towards center from eye)
	s := f.Cross(up.Normalize()).Normalize() // side (right)
	u := s.Cross(f)                          // up (corrected)

	return Matrix4{
		{s.X, u.X, -f.X, 0},
		{s.Y, u.Y, -f.Y, 0},
		{s.Z, u.Z, -f.Z, 0},
		{-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1},
	}
}
