// Command plan 打印按帧预算排出的特效播放计划。
package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"example.com/vfx-budget/plan"
)

func scaleRequests(count int) []plan.Request {
	out := make([]plan.Request, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, plan.Request{
			ID: fmt.Sprintf("v%04d", i), Frame: 1,
			Key: fmt.Sprintf("k%d", i%20), Priority: (i % 5) + 1,
		})
	}
	return out
}

func scenario(name string) (*plan.Planner, []plan.Request) {
	switch name {
	case "budget":
		return &plan.Planner{Budget: 2, CooldownFrames: 1}, []plan.Request{
			{ID: "a", Frame: 1, Key: "fire", Priority: 3},
			{ID: "b", Frame: 1, Key: "ice", Priority: 2},
			{ID: "c", Frame: 1, Key: "bolt", Priority: 1},
		}
	case "preempt":
		return &plan.Planner{Budget: 1, CooldownFrames: 1}, []plan.Request{
			{ID: "a", Frame: 1, Key: "fire", Priority: 1},
			{ID: "b", Frame: 1, Key: "ice", Priority: 5},
		}
	case "cooldown":
		return &plan.Planner{Budget: 4, CooldownFrames: 5}, []plan.Request{
			{ID: "a", Frame: 1, Key: "fire", Priority: 1},
			{ID: "b", Frame: 2, Key: "fire", Priority: 1},
		}
	case "order":
		return &plan.Planner{Budget: 5, CooldownFrames: 1}, []plan.Request{
			{ID: "c", Frame: 1, Key: "fire", Priority: 1},
			{ID: "a", Frame: 1, Key: "ice", Priority: 1},
		}
	case "work":
		return &plan.Planner{Budget: 100000, CooldownFrames: 1}, scaleRequests(3000)
	}
	fmt.Fprintf(os.Stderr, "未知场景 %q\n", name)
	os.Exit(2)
	return nil, nil
}

func main() {
	sample := flag.String("sample", "budget", "场景名")
	flag.Parse()

	p, requests := scenario(*sample)
	rows := p.Plan(requests)

	if *sample == "work" {
		fmt.Printf("{\"played\":%d,\"work\":%d}\n", len(rows), p.Work)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "frame\tid\tkey")
	for _, r := range rows {
		fmt.Fprintf(w, "%d\t%s\t%s\n", r.Frame, r.ID, r.Key)
	}
	fmt.Fprintf(w, "-\tpreempted\t%d\n", p.Preempt)
	w.Flush()
}
