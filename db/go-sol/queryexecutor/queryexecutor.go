package queryexecutor

type Row []any

type Node interface {
	SetChild(Node)
	GetChild() Node
	Next() Row
}

type childAccessor struct{ child Node }

func (c *childAccessor) GetChild() Node  { return c.child }
func (c *childAccessor) SetChild(n Node) { c.child = n }

type MemoryScan struct {
	childAccessor
	Rows []Row
	i    int
}

func (m *MemoryScan) Next() Row {
	if m.i >= len(m.Rows) {
		return nil
	}
	row := m.Rows[m.i]
	m.i++
	return row
}

type Selection struct {
	childAccessor
	Filter FilterFunc
}

func (s *Selection) Next() Row {
	for {
		row := s.GetChild().Next()
		if row == nil {
			return nil
		}
		if s.Filter(row) {
			return row
		}
	}
}

type Projection struct {
	childAccessor
	Mapper MapperFunc
}

func (p *Projection) Next() Row {
	row := p.GetChild().Next()
	if row != nil {
		return p.Mapper(row)
	}
	return nil
}

type Limit struct {
	childAccessor
	L int
}

func (l *Limit) Next() Row {
	if l.L <= 0 {
		return nil
	}
	l.L -= 1
	return l.GetChild().Next()
}
