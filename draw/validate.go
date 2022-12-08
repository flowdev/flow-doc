package draw

import "fmt"

const (
	unknown = iota
	comp
	arrow
)

type sequels struct {
	open map[int]bool
	all  map[int]bool
}

func newSequels() *sequels {
	return &sequels{
		open: make(map[int]bool, 32),
		all:  make(map[int]bool, 32),
	}
}

func (s *sequels) add(n int) {
	s.all[n] = true
	s.open[n] = true
}

func (s *sequels) exists(n int) bool {
	return s.all[n]
}

func (s *sequels) isOpen(n int) bool {
	return s.open[n]
}

func (s *sequels) close(n int) {
	delete(s.open, n)
}

func validateFlowData(f *Flow) error {
	merges := make(map[string]*Merge, 64)
	err := validateSplit(f.AllShapes, false, newSequels(), merges)
	cleanMerges(merges)
	return err
}

func validateSplit(split *Split, inner bool, seqs *sequels, merges map[string]*Merge) error {
	shapes := split.Shapes
	if len(shapes) <= 0 {
		if inner {
			return fmt.Errorf("no shapes found in inner split")
		}
		return fmt.Errorf("no shapes found at all")
	}

	for i, row := range shapes {
		lastShape := unknown
		if inner {
			lastShape = comp
		}
		lastIdx := len(row) - 1
		for j, ishape := range row {
			switch shape := ishape.(type) {
			case *Arrow:
				if lastShape == arrow {
					return fmt.Errorf(
						"two arrows in a row aren't allowed "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				lastShape = arrow
			case *Comp:
				if lastShape == comp && !validMergeComp(shape, merges) {
					return fmt.Errorf(
						"two components in a row aren't allowed; "+
							"second one is %q "+
							"(row index %d, column index %d, inner split: %t)",
						compID(shape), i, j, inner,
					)
				}
				lastShape = comp
			case *Sequel:
				if j != lastIdx && j != 0 {
					return fmt.Errorf(
						"a sequel has to either be the last or first shape of a row "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if j == 0 {
					if inner {
						return fmt.Errorf(
							"sequel Number %d can't end on an inner level "+
								"(row index %d, column index %d, inner split: %t)",
							shape.Number, i, j, inner,
						)
					}
					if !seqs.exists(shape.Number) {
						return fmt.Errorf(
							"sequel Number %d doesn't exist at all "+
								"(row index %d, column index %d, inner split: %t)",
							shape.Number, i, j, inner,
						)
					}
					if !seqs.isOpen(shape.Number) {
						return fmt.Errorf(
							"sequel Number %d has ended already "+
								"(row index %d, column index %d, inner split: %t)",
							shape.Number, i, j, inner,
						)
					}
					seqs.close(shape.Number)
				}
				if j == lastIdx {
					if lastShape != arrow {
						return fmt.Errorf(
							"a sequel has to follow an arrow "+
								"(row index %d, column index %d, inner split: %t)",
							i, j, inner,
						)
					}
					if seqs.exists(shape.Number) {
						return fmt.Errorf(
							"sequel Number %d has been started already "+
								"(row index %d, column index %d, inner split: %t)",
							shape.Number, i, j, inner,
						)
					}
					seqs.add(shape.Number)
				}
				lastShape = comp
			case *Merge:
				if j != lastIdx {
					return fmt.Errorf(
						"a merge has to be the last shape of a row "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if lastShape != arrow {
					return fmt.Errorf(
						"a merge has to follow an arrow "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				registerMerge(shape, merges)
				lastShape = comp
			case *Loop:
				if j != lastIdx {
					return fmt.Errorf(
						"a loop has to be the last shape of a row "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if lastShape != arrow {
					return fmt.Errorf(
						"a loop has to follow an arrow "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				lastShape = comp
			case *ExtPort:
				if j != lastIdx && j != 0 {
					return fmt.Errorf(
						"an external port has to either be the last or first shape of a row "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if j == lastIdx && lastShape != arrow {
					return fmt.Errorf(
						"an external port has to follow an arrow "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				lastShape = comp
			case *Split:
				if lastShape != comp {
					return fmt.Errorf(
						"a split can only follow a component "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if j != lastIdx {
					return fmt.Errorf(
						"a split has to be the last shape of a row "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if j == 0 {
					return fmt.Errorf(
						"a split can't be the first shape of a row "+
							"(row index %d, column index %d, inner split: %t)",
						i, j, inner,
					)
				}
				if err := validateSplit(shape, true, seqs, merges); err != nil {
					return err
				}
				lastShape = arrow
			default:
				return fmt.Errorf(
					"unsupported shape type %T found "+
						"(row index %d, column index %d, inner split: %t)",
					ishape, i, j, inner)
			}
		}
	}
	return nil
}

func cleanMerges(merges map[string]*Merge) {
	for _, merge := range merges {
		merge.drawData = nil
		merge.arrows = nil
	}
}

func validMergeComp(comp *Comp, merges map[string]*Merge) bool {
	m, ok := merges[compID(comp)]
	if !ok {
		return false
	}
	return m.Size == len(m.arrows)
}

func registerMerge(merge *Merge, merges map[string]*Merge) {
	merges[merge.ID] = merge // might be done already but we insist on the same object
	merge.arrows = append(merge.arrows, &Arrow{})
}
