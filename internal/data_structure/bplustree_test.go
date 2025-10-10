package data_structure

import (
	"testing"
)

func TestCompareTo(t *testing.T) {
	a := &Item{Score: 1, Member: "a"}
	b := &Item{Score: 2, Member: "b"}
	if a.CompareTo(b) >= 0 {
		t.Fatal("expected a < b")
	}
	if b.CompareTo(a) <= 0 {
		t.Fatal("expected b > a")
	}
	c := &Item{Score: 1, Member: "c"}
	if a.CompareTo(c) >= 0 {
		t.Fatal("expected a < c when same score")
	}
}

func TestScoreAndAdd(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(10, "a")
	tree.Add(20, "b")
	tree.Add(15, "c")
	if s, ok := tree.Score("b"); !ok || s != 20 {
		t.Fatalf("expected 20 got %v", s)
	}
	tree.Add(25, "b")
	if s, _ := tree.Score("b"); s != 25 {
		t.Fatalf("score not updated")
	}
}

func TestGetRankSimple(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(10, "a")
	tree.Add(20, "b")
	tree.Add(15, "c")
	rank := tree.GetRank("a")
	if rank < 0 {
		t.Fatal("rank should exist")
	}
}

func TestSplitLeaf(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(1, "a")
	tree.Add(2, "b")
	tree.Add(3, "c")
	tree.Add(4, "d")
	if tree.Root == nil {
		t.Fatal("root should exist")
	}
	if !tree.Root.IsLeaf && len(tree.Root.Children) == 0 {
		t.Fatal("expected root to have children after split")
	}
}

func TestUpdateExistingMember(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(5, "x")
	tree.Add(5, "x") // same score -> no change
	s, _ := tree.Score("x")
	if s != 5 {
		t.Fatalf("expected score 5, got %v", s)
	}
	tree.Add(10, "x")
	s, _ = tree.Score("x")
	if s != 10 {
		t.Fatalf("expected updated score 10")
	}
}

func TestRankOrder(t *testing.T) {
	tree := NewBPlusTree(3)
	members := []string{"a", "b", "c", "d"}
	scores := []float64{1, 2, 3, 4}
	for i := range members {
		tree.Add(scores[i], members[i])
	}
	for i, m := range members {
		if r := tree.GetRank(m); r != i {
			t.Fatalf("expected rank %d for %s, got %d", i, m, r)
		}
	}
}

func TestNewBPlusTree(t *testing.T) {
	tests := []struct {
		name   string
		degree int
	}{
		{"small degree", 3},
		{"medium degree", 5},
		{"large degree", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewBPlusTree(tt.degree)
			if tree == nil {
				t.Fatal("NewBPlusTree returned nil")
			}
			if tree.Degree != tt.degree {
				t.Errorf("expected degree %d, got %d", tt.degree, tree.Degree)
			}
			if tree.Root == nil {
				t.Fatal("Root node is nil")
			}
			if !tree.Root.IsLeaf {
				t.Error("Root should be a leaf node initially")
			}
		})
	}
}

func TestBPlusTree_Add_SingleItem(t *testing.T) {
	tree := NewBPlusTree(3)

	result := tree.Add(10.5, "member1")

	if result != 1 {
		t.Errorf("expected Add to return 1, got %d", result)
	}

	score, found := tree.Score("member1")
	if !found {
		t.Error("member1 should be found after adding")
	}
	if score != 10.5 {
		t.Errorf("expected score 10.5, got %f", score)
	}
}

func TestBPlusTree_Add_EmptyMember(t *testing.T) {
	tree := NewBPlusTree(3)

	result := tree.Add(10.5, "")

	if result != 0 {
		t.Errorf("expected Add to return 0 for empty member, got %d", result)
	}
}

func TestBPlusTree_Add_UpdateExisting(t *testing.T) {
	tree := NewBPlusTree(3)

	// Add initial item
	tree.Add(10.0, "member1")

	// Update with new score
	result := tree.Add(20.0, "member1")

	if result != 1 {
		t.Errorf("expected Add to return 1 when updating, got %d", result)
	}

	score, found := tree.Score("member1")
	if !found {
		t.Error("member1 should still be found after update")
	}
	if score != 20.0 {
		t.Errorf("expected updated score 20.0, got %f", score)
	}
}

func TestBPlusTree_Add_MultipleItems(t *testing.T) {
	tree := NewBPlusTree(3)

	items := []struct {
		score  float64
		member string
	}{
		{10.0, "alice"},
		{20.0, "bob"},
		{15.0, "charlie"},
		{5.0, "david"},
		{25.0, "eve"},
	}

	for _, item := range items {
		result := tree.Add(item.score, item.member)
		if result != 1 {
			t.Errorf("failed to add %s with score %f", item.member, item.score)
		}
	}

	// Verify all items exist
	for _, item := range items {
		score, found := tree.Score(item.member)
		if !found {
			t.Errorf("member %s not found", item.member)
		}
		if score != item.score {
			t.Errorf("member %s: expected score %f, got %f", item.member, item.score, score)
		}
	}
}

func TestBPlusTree_Add_TriggerSplit(t *testing.T) {
	tree := NewBPlusTree(3) // Max 2 items per node before split

	// Add items to trigger splits
	for i := 0; i < 10; i++ {
		result := tree.Add(float64(i), string(rune('a'+i)))
		if result != 1 {
			t.Errorf("failed to add item %d", i)
		}
	}

	// Verify tree structure changed (root should no longer be a leaf)
	if tree.Root.IsLeaf {
		t.Error("root should not be a leaf after multiple splits")
	}

	// Verify all items are still accessible
	for i := 0; i < 10; i++ {
		member := string(rune('a' + i))
		score, found := tree.Score(member)
		if !found {
			t.Errorf("member %s not found after splits", member)
		}
		if score != float64(i) {
			t.Errorf("member %s: expected score %d, got %f", member, i, score)
		}
	}
}

func TestBPlusTree_Score_NotFound(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(10.0, "alice")

	score, found := tree.Score("bob")

	if found {
		t.Error("should not find non-existent member")
	}
	if score != 0 {
		t.Errorf("expected score 0 for non-existent member, got %f", score)
	}
}

func TestBPlusTree_Score_EmptyTree(t *testing.T) {
	tree := NewBPlusTree(3)

	score, found := tree.Score("alice")

	if found {
		t.Error("should not find member in empty tree")
	}
	if score != 0 {
		t.Errorf("expected score 0 for empty tree, got %f", score)
	}
}

func TestBPlusTree_GetRank_SingleItem(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(10.0, "alice")

	rank := tree.GetRank("alice")

	if rank != 0 {
		t.Errorf("expected rank 0, got %d", rank)
	}
}

func TestBPlusTree_GetRank_NotFound(t *testing.T) {
	tree := NewBPlusTree(3)
	tree.Add(10.0, "alice")

	rank := tree.GetRank("bob")

	if rank != -1 {
		t.Errorf("expected rank -1 for non-existent member, got %d", rank)
	}
}

func TestBPlusTree_GetRank_MultipleItems(t *testing.T) {
	tree := NewBPlusTree(5)

	// Add items in specific order
	tree.Add(10.0, "alice")   // rank 0
	tree.Add(20.0, "bob")     // rank 1
	tree.Add(15.0, "charlie") // rank 1, alice moves to 0, bob to 2
	tree.Add(5.0, "david")    // rank 0
	tree.Add(25.0, "eve")     // rank 4

	tests := []struct {
		member       string
		expectedRank int
	}{
		{"david", 0},   // score 5.0
		{"alice", 1},   // score 10.0
		{"charlie", 2}, // score 15.0
		{"bob", 3},     // score 20.0
		{"eve", 4},     // score 25.0
	}

	for _, tt := range tests {
		t.Run(tt.member, func(t *testing.T) {
			rank := tree.GetRank(tt.member)
			if rank != tt.expectedRank {
				t.Errorf("member %s: expected rank %d, got %d", tt.member, tt.expectedRank, rank)
			}
		})
	}
}

func TestBPlusTree_GetRank_AfterSplits(t *testing.T) {
	tree := NewBPlusTree(3)

	// Add many items to cause splits
	for i := 0; i < 20; i++ {
		tree.Add(float64(i*10), string(rune('a'+i)))
	}

	// Verify ranks are correct
	for i := 0; i < 20; i++ {
		member := string(rune('a' + i))
		rank := tree.GetRank(member)
		if rank != i {
			t.Errorf("member %s: expected rank %d, got %d", member, i, rank)
		}
	}
}

func TestBPlusTree_GetRank_EmptyTree(t *testing.T) {
	tree := NewBPlusTree(3)

	rank := tree.GetRank("alice")

	if rank != -1 {
		t.Errorf("expected rank -1 for empty tree, got %d", rank)
	}
}

func TestItem_CompareTo(t *testing.T) {
	tests := []struct {
		name     string
		item1    *Item
		item2    *Item
		expected int
	}{
		{
			name:     "item1 score less than item2",
			item1:    &Item{Score: 10.0, Member: "alice"},
			item2:    &Item{Score: 20.0, Member: "bob"},
			expected: -1,
		},
		{
			name:     "item1 score greater than item2",
			item1:    &Item{Score: 30.0, Member: "charlie"},
			item2:    &Item{Score: 20.0, Member: "bob"},
			expected: 1,
		},
		{
			name:     "equal scores, item1 member less than item2",
			item1:    &Item{Score: 20.0, Member: "alice"},
			item2:    &Item{Score: 20.0, Member: "bob"},
			expected: -1,
		},
		{
			name:     "equal scores, item1 member greater than item2",
			item1:    &Item{Score: 20.0, Member: "charlie"},
			item2:    &Item{Score: 20.0, Member: "bob"},
			expected: 1,
		},
		{
			name:     "completely equal items",
			item1:    &Item{Score: 20.0, Member: "alice"},
			item2:    &Item{Score: 20.0, Member: "alice"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item1.CompareTo(tt.item2)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestBPlusTree_Add_SameScore_DifferentMembers(t *testing.T) {
	tree := NewBPlusTree(5)

	// Add multiple members with same score
	tree.Add(100.0, "alice")
	tree.Add(100.0, "bob")
	tree.Add(100.0, "charlie")

	// Verify all exist
	for _, member := range []string{"alice", "bob", "charlie"} {
		score, found := tree.Score(member)
		if !found {
			t.Errorf("member %s not found", member)
		}
		if score != 100.0 {
			t.Errorf("member %s: expected score 100.0, got %f", member, score)
		}
	}

	// Verify ranks are in lexicographic order
	tests := []struct {
		member       string
		expectedRank int
	}{
		{"alice", 0},
		{"bob", 1},
		{"charlie", 2},
	}

	for _, tt := range tests {
		rank := tree.GetRank(tt.member)
		if rank != tt.expectedRank {
			t.Errorf("member %s: expected rank %d, got %d", tt.member, tt.expectedRank, rank)
		}
	}
}

func TestBPlusTree_StressTest(t *testing.T) {
	tree := NewBPlusTree(4)

	// Add 100 items
	for i := 0; i < 100; i++ {
		result := tree.Add(float64(i), string(rune('a'+(i%26)))+string(rune('a'+(i/26))))
		if result != 1 && i >= 26 { // First 26 are unique, after that we have duplicates that should update
			if result != 1 {
				t.Errorf("failed to add/update item %d", i)
			}
		}
	}

	// Verify tree integrity by checking that we can retrieve items
	node := tree.Root
	for !node.IsLeaf {
		node = node.Children[0]
	}

	itemCount := 0
	for node != nil {
		itemCount += len(node.Items)
		node = node.Next
	}

	if itemCount == 0 {
		t.Error("tree should contain items after stress test")
	}
}

func TestBPlusTree_LeafLinking(t *testing.T) {
	tree := NewBPlusTree(3)

	// Add items to create multiple leaf nodes
	for i := 0; i < 10; i++ {
		tree.Add(float64(i), string(rune('a'+i)))
	}

	// Find first leaf
	node := tree.Root
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse leaf nodes and count items
	totalItems := 0
	for node != nil {
		if !node.IsLeaf {
			t.Error("encountered non-leaf node while traversing leaf chain")
		}
		totalItems += len(node.Items)
		node = node.Next
	}

	if totalItems != 10 {
		t.Errorf("expected 10 items across all leaves, got %d", totalItems)
	}
}

func TestBPlusTree_NegativeScores(t *testing.T) {
	tree := NewBPlusTree(3)

	tree.Add(-10.0, "alice")
	tree.Add(-5.0, "bob")
	tree.Add(-15.0, "charlie")

	tests := []struct {
		member        string
		expectedScore float64
		expectedRank  int
	}{
		{"charlie", -15.0, 0},
		{"alice", -10.0, 1},
		{"bob", -5.0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.member, func(t *testing.T) {
			score, found := tree.Score(tt.member)
			if !found {
				t.Errorf("member %s not found", tt.member)
			}
			if score != tt.expectedScore {
				t.Errorf("expected score %f, got %f", tt.expectedScore, score)
			}

			rank := tree.GetRank(tt.member)
			if rank != tt.expectedRank {
				t.Errorf("expected rank %d, got %d", tt.expectedRank, rank)
			}
		})
	}
}

func TestBPlusTree_FloatingPointScores(t *testing.T) {
	tree := NewBPlusTree(5)

	tree.Add(1.1, "alice")
	tree.Add(1.11, "bob")
	tree.Add(1.111, "charlie")

	// Verify order is maintained
	rank1 := tree.GetRank("alice")
	rank2 := tree.GetRank("bob")
	rank3 := tree.GetRank("charlie")

	if !(rank1 < rank2 && rank2 < rank3) {
		t.Errorf("ranks not in expected order: %d, %d, %d", rank1, rank2, rank3)
	}
}
