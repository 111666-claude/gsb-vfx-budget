// Package plan 按帧预算挑选要播放的特效。
package plan

// Request 是一次特效播放请求。
type Request struct {
	ID       string
	Frame    int
	Key      string
	Priority int
}

// Row 是计划里的一条播放记录。
type Row struct {
	Frame int
	ID    string
	Key   string
}

// Planner 生成播放计划。
type Planner struct {
	Budget         int
	CooldownFrames int

	played    []Row
	seen      map[string]bool
	lastFrame map[string]int
	Work      int
	Preempt   int
}

// Plan 依次处理请求，返回每一帧实际播放的记录。
func (p *Planner) Plan(requests []Request) []Row {
	if p.seen == nil {
		p.seen = map[string]bool{}
	}
	var out []Row
	for _, r := range requests {
		p.seen[r.ID] = true

		used := 0
		for _, prev := range p.played {
			p.Work++
			if prev.Frame == r.Frame {
				used++
			}
		}
		if used >= p.Budget {
			continue
		}

		row := Row{Frame: r.Frame, ID: r.ID, Key: r.Key}
		p.played = append(p.played, row)
		out = append(out, row)
	}
	return out
}
