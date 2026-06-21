package geom

// Camera defines the properties for viewing a 3D scene.
type Camera struct {
	Position Vector3
	Target   Vector3
	Up       Vector3
}
