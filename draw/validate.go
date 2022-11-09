package draw

import "fmt"

func validateFlowData(f *Flow) error {
	return validateSplit(f.AllShapes)
}

func validateSplit(split *Split) error {
	shapes := split.Shapes

	if len(shapes) <= 0 {
		return fmt.Errorf("No shapes found")
	}
	for i, row := range shapes {
		for j, ishape := range row {
			switch shape := ishape.(type) {
			case *Arrow, *Comp, *Merge, *Text, *Sequel, *Loop:
				break
			case *Split:
				err := validateSplit(shape)
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
