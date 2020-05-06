package protometry

//// Implements the Unmarshaler interface of the yaml pkg.
//func (v *Vector3) UnmarshalYAML(value *yaml.Node) error {
//	t := Vector3{}
//	err := yaml.Unmarshal([]byte(value.), &t)
//	if err != nil {
//		return err
//	}
//
//	// make sure to dereference before assignment,
//	// otherwise only the local variable will be overwritten
//	// and not the value the pointer actually points to
//	v.X = yamlVector3[0]
//	v.Y = yamlVector3[1]
//	v.Z = yamlVector3[2]
//
//	return nil
//}
