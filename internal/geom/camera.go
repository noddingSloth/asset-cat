package geom

import "math"

// Camera defines the properties for viewing a 3D scene.
type Camera struct {
	Position Vector3
	Target   Vector3
	Up       Vector3
	FOV      float64 // Field of view in degrees
	Near     float64
	Far      float64
}

// ViewMatrix returns the view (look-at) matrix for this camera.
func (c *Camera) ViewMatrix() Matrix4 {
	return LookAt(c.Position, c.Target, c.Up)
}

// ProjectionMatrix returns the perspective projection matrix for this camera.
// aspect is width / height of the viewport.
func (c *Camera) ProjectionMatrix(aspect float64) Matrix4 {
	fovRad := c.FOV * math.Pi / 180.0
	near := c.Near
	far := c.Far
	if near <= 0 {
		near = 0.1
	}
	if far <= 0 {
		far = 100.0
	}
	return Perspective(fovRad, aspect, near, far)
}

// ViewProjectionMatrix returns the combined view-projection matrix.
func (c *Camera) ViewProjectionMatrix(aspect float64) Matrix4 {
	view := c.ViewMatrix()
	proj := c.ProjectionMatrix(aspect)
	return proj.Multiply(view)
}

// DefaultCamera returns a camera with sensible defaults for viewing a model.
func DefaultCamera() *Camera {
	return &Camera{
		Position: Vector3{X: 0, Y: 0, Z: 3},
		Target:   Vector3{X: 0, Y: 0, Z: 0},
		Up:       Vector3{X: 0, Y: 1, Z: 0},
		FOV:      60.0,
		Near:     0.1,
		Far:      100.0,
	}
}
