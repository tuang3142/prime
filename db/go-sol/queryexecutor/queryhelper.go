package queryexecutor

// TODO: filter func and projector func needs to know schema before hands. how can we design this?
type FilterFunc func(Row) bool

func And(fs ...FilterFunc) FilterFunc {
	return func(r Row) bool {
		ret := true
		for _, f := range fs {
			ret = ret && f(r)
		}
		return ret
	}
}

func Or(fs ...FilterFunc) FilterFunc {
	return func(r Row) bool {
		ret := false
		for _, f := range fs {
			ret = ret || f(r)
		}
		return ret
	}
}

type MapperFunc func(Row) Row

type SortFunc func(Row) int

type AggFunc func([]Row) Row
