package draw

import "fmt"

func validateFlowData(f Flow) error {
	return validateShapes(f.Shapes.Shapes)
}

func validateShapes(shapes [][]any) error {
	if len(shapes) <= 0 {
		return fmt.Errorf("No shapes found")
	}
	for i, row := range shapes {
		for j, ishape := range row {
			switch shape := ishape.(type) {
			case *Arrow, *Op, *Merge, *Text:
				break
			case *Split:
				err := validateShapes(shape.Shapes)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf(
					"unsupported shape type %T at row index %d and column index %d",
					ishape, i, j)
			}
		}
	}
	return nil
}
