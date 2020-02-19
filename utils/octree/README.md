# Octree

```go
// Create
oct := NewOctree(erutan.NetVector3{X: 0, Y: 0, Z: 0}, erutan.NetVector3{X: 1, Y: 1, Z: 1})

// Add element at point
oct.Add(1, erutan.NetVector3{X: 0.1, Y:  0.2, Z: 0.3})
oct.Add(2, erutan.NetVector3{X: 0.2, Y: 0.3, Z: 0.4})
oct.Add(3, erutan.NetVector3{X: 0.3, Y: 0.4, Z: 0.5})
node4 := oct.Add(4, erutan.NetVector3{X: 0.3, Y: 0.4, Z; 0.5}) // save for removal later

// Retrieval at point
oct.ElementsAt(erutan.NetVector3{X: 0.1, Y: 0.2, Z: 0.3}) // [1]
oct.ElementsAt(erutan.NetVector3{X: 0.2, Y: 0.3, Z: 0.4}) // [2]
oct.ElementsAt(erutan.NetVector3{X: 0.3, Y: 0.4, Z: 0.5}) // [3 4]

// Retrieval in box
oct.ElementsIn(Box{erutan.NetVector3{X: 0.1, Y: 0.2, Z: 0.3}, erutan.NetVector3{X: 0.2, Y: 0.3, Z: 0.4}}) // [1 2]

// Remove first of element in tree (slower)
oct.Remove(1) // true

// Remove first of element within node (faster)
oct.RemoveUsing(4, node4) // true

// Clear contents
oct.Clear() // true
```

## Tests

```bash
go test github.com/user/erutan/utils/octree -v
```

## Roadmap

- [ ] Fix naming