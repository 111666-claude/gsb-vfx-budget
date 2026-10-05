// Package plan 按帧预算挑选要播放的特效。
package plan

import "sort"

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
	prio      []int
	byFrame   map[int][]int
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
	if p.lastFrame == nil {
		p.lastFrame = map[string]int{}
	}
	if p.byFrame == nil {
		p.byFrame = map[int][]int{}
	}
	for _, r := range requests {
		if p.seen[r.ID] {
			continue
		}
		p.seen[r.ID] = true

		if last, ok := p.lastFrame[r.Key]; ok && r.Frame-last < p.CooldownFrames {
			continue
		}

		idxs := p.byFrame[r.Frame]
		if len(idxs) >= p.Budget {
			low := -1
			for _, i := range idxs {
				p.Work++
				if low == -1 || p.prio[i] < p.prio[low] {
					low = i
				}
			}
			if low == -1 || r.Priority <= p.prio[low] {
				continue
			}
			removed := p.played[low]
			if p.lastFrame[removed.Key] == removed.Frame {
				delete(p.lastFrame, removed.Key)
			}
			p.played[low] = Row{Frame: r.Frame, ID: r.ID, Key: r.Key}
			p.prio[low] = r.Priority
			p.lastFrame[r.Key] = r.Frame
			p.Preempt++
			continue
		}

		p.byFrame[r.Frame] = append(idxs, len(p.played))
		row := Row{Frame: r.Frame, ID: r.ID, Key: r.Key}
		p.played = append(p.played, row)
		p.prio = append(p.prio, r.Priority)
		p.lastFrame[r.Key] = r.Frame
	}

	out := make([]Row, len(p.played))
	copy(out, p.played)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Frame != out[j].Frame {
			return out[i].Frame < out[j].Frame
		}
		return out[i].ID < out[j].ID
	})
	return out
}
