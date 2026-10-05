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

type entry struct {
	row      Row
	priority int
}

// Planner 生成播放计划。
type Planner struct {
	Budget         int
	CooldownFrames int

	frames    map[int][]entry
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
	if p.frames == nil {
		p.frames = map[int][]entry{}
	}
	if p.lastFrame == nil {
		p.lastFrame = map[string]int{}
	}

	for _, r := range requests {
		if p.seen[r.ID] {
			continue
		}
		p.seen[r.ID] = true

		if p.Budget <= 0 {
			continue
		}

		if last, ok := p.lastFrame[r.Key]; ok && r.Frame-last < p.CooldownFrames {
			continue
		}

		entries := p.frames[r.Frame]
		e := entry{row: Row{Frame: r.Frame, ID: r.ID, Key: r.Key}, priority: r.Priority}

		if len(entries) < p.Budget {
			p.frames[r.Frame] = append(entries, e)
			p.lastFrame[r.Key] = r.Frame
			continue
		}

		low := 0
		for i := range entries {
			p.Work++
			if entries[i].priority < entries[low].priority {
				low = i
			}
		}
		if r.Priority <= entries[low].priority {
			continue
		}
		old := entries[low]
		if p.lastFrame[old.row.Key] == old.row.Frame {
			delete(p.lastFrame, old.row.Key)
		}
		entries[low] = e
		p.frames[r.Frame] = entries
		p.lastFrame[r.Key] = r.Frame
		p.Preempt++
	}

	frameNums := make([]int, 0, len(p.frames))
	for f := range p.frames {
		frameNums = append(frameNums, f)
	}
	sort.Ints(frameNums)

	var out []Row
	for _, f := range frameNums {
		entries := p.frames[f]
		sort.Slice(entries, func(i, j int) bool { return entries[i].row.ID < entries[j].row.ID })
		for _, e := range entries {
			out = append(out, e.row)
		}
	}
	return out
}
