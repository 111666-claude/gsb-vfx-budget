package plan

import "testing"

func TestFirstRequestIsPlayed(t *testing.T) {
	p := &Planner{Budget: 4, CooldownFrames: 2}
	got := p.Plan([]Request{{ID: "a", Frame: 1, Key: "fire", Priority: 1}})
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("第一条应播放：%+v", got)
	}
}

func TestPlayedRowsCarryFrameAndKey(t *testing.T) {
	p := &Planner{Budget: 4, CooldownFrames: 2}
	got := p.Plan([]Request{{ID: "a", Frame: 7, Key: "spark", Priority: 1}})
	if got[0].Frame != 7 || got[0].Key != "spark" {
		t.Fatalf("帧与特效键应原样带上：%+v", got[0])
	}
}

func TestDifferentFramesBothPlay(t *testing.T) {
	p := &Planner{Budget: 4, CooldownFrames: 1}
	got := p.Plan([]Request{
		{ID: "a", Frame: 1, Key: "fire", Priority: 1},
		{ID: "b", Frame: 9, Key: "fire", Priority: 1},
	})
	if len(got) != 2 {
		t.Fatalf("相隔很远的两帧都该播：%+v", got)
	}
}

func TestEmptyInput(t *testing.T) {
	p := &Planner{Budget: 4}
	if got := p.Plan(nil); len(got) != 0 {
		t.Fatalf("空输入应为空：%+v", got)
	}
}
