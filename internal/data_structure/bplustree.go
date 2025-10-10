package data_structure

type Item struct {
	Score  float64
	Member string
}

func (i *Item) CompareTo(other *Item) int {
	if i.Score < other.Score {
		return -1
	}
	if i.Score > other.Score {
		return 1
	}
	// Tie-break by member name
	if i.Member < other.Member {
		return -1
	}
	if i.Member > other.Member {
		return 1
	}
	return 0
}

type Node struct {
	Items    []*Item
	Children []*Node
	IsLeaf   bool
	Parent   *Node
	Next     *Node // Linked list of leaves for range scans
	Size     int   // Total items in subtree (for O(log N) rank)
}

type BPlusTree struct {
	Root      *Node
	Degree    int
	MemberMap map[string]*Item // O(1) lookup
}

func NewBPlusTree(degree int) *BPlusTree {
	if degree < 3 {
		degree = 64 // Optimal for cache lines
	}
	return &BPlusTree{
		Root: &Node{
			IsLeaf: true,
			Size:   0,
		},
		Degree:    degree,
		MemberMap: make(map[string]*Item),
	}
}

func (t *BPlusTree) Score(member string) (float64, bool) {
	item, exist := t.MemberMap[member]
	if !exist {
		return 0, false
	}
	return item.Score, true
}

func (t *BPlusTree) Add(score float64, member string) int {
	if len(member) == 0 {
		return 0
	}

	if existingItem, exist := t.MemberMap[member]; exist {
		oldScore := existingItem.Score
		if oldScore == score {
			return 1
		}
		// Reposition: remove, update score, re-insert
		t.removeFromTree(member, oldScore)
		existingItem.Score = score
		t.insertIntoTree(existingItem)
		return 1
	}

	// New member
	item := &Item{Score: score, Member: member}
	t.MemberMap[member] = item
	t.insertIntoTree(item)
	return 1
}

func (t *BPlusTree) insertIntoTree(item *Item) {
	node := t.Root

	// Navigate to leaf, increment sizes
	for !node.IsLeaf {
		node.Size++
		i := 0
		for i < len(node.Items) && item.Score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// Insert in sorted order
	i := 0
	for i < len(node.Items) && item.CompareTo(node.Items[i]) >= 0 {
		i++
	}
	node.Items = append(node.Items[:i], append([]*Item{item}, node.Items[i:]...)...)
	node.Size++

	// Split if over capacity
	if len(node.Items) > t.Degree-1 {
		t.splitNode(node)
	}
}

func (t *BPlusTree) removeFromTree(member string, score float64) {
	node := t.Root

	// Navigate to leaf, decrement sizes
	for !node.IsLeaf {
		node.Size--
		i := 0
		for i < len(node.Items) && score >= node.Items[i].Score {
			i++
		}
		node = node.Children[i]
	}

	// Remove from leaf (no rebalancing)
	for i, item := range node.Items {
		if item.Member == member {
			node.Items = append(node.Items[:i], node.Items[i+1:]...)
			node.Size--
			return
		}
	}
}

func (t *BPlusTree) splitNode(node *Node) {
	if node.Parent == nil {
		t.splitRoot()
		return
	}

	if node.IsLeaf {
		t.splitLeaf(node)
	} else {
		t.splitInternal(node)
	}
}

func (t *BPlusTree) splitLeaf(node *Node) {
	medianIndex := len(node.Items) / 2

	newLeaf := &Node{
		IsLeaf: true,
		Parent: node.Parent,
		Next:   node.Next,
		Size:   len(node.Items) - medianIndex,
	}

	// Move second half to new leaf
	newLeaf.Items = append(newLeaf.Items, node.Items[medianIndex:]...)
	node.Items = node.Items[:medianIndex]
	node.Size = len(node.Items)
	node.Next = newLeaf

	// Promote first key of new leaf to parent
	parent := node.Parent
	promotedItem := newLeaf.Items[0]

	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*Node{newLeaf}, parent.Children[childIndex+1:]...)...)

	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) splitInternal(node *Node) {
	medianIndex := len(node.Items) / 2

	// Save promoted item before trimming
	promotedItem := node.Items[medianIndex]

	newInternal := &Node{
		IsLeaf: false,
		Parent: node.Parent,
	}

	// Move second half to new node
	newInternal.Items = append(newInternal.Items, node.Items[medianIndex+1:]...)
	newInternal.Children = append(newInternal.Children, node.Children[medianIndex+1:]...)

	// Recalculate sizes
	newInternal.Size = 0
	for _, child := range newInternal.Children {
		newInternal.Size += child.Size
	}

	node.Items = node.Items[:medianIndex]
	node.Children = node.Children[:medianIndex+1]

	node.Size = 0
	for _, child := range node.Children {
		node.Size += child.Size
	}

	// Update parent pointers
	for _, child := range newInternal.Children {
		child.Parent = newInternal
	}

	parent := node.Parent
	childIndex := 0
	for childIndex < len(parent.Children) {
		if parent.Children[childIndex] == node {
			break
		}
		childIndex++
	}

	parent.Items = append(parent.Items[:childIndex], append([]*Item{promotedItem}, parent.Items[childIndex:]...)...)
	parent.Children = append(parent.Children[:childIndex+1], append([]*Node{newInternal}, parent.Children[childIndex+1:]...)...)

	if len(parent.Items) > t.Degree-1 {
		t.splitNode(parent)
	}
}

func (t *BPlusTree) splitRoot() {
	oldRoot := t.Root
	newRoot := &Node{
		IsLeaf: false,
		Size:   oldRoot.Size,
	}
	t.Root = newRoot
	oldRoot.Parent = newRoot
	newRoot.Children = append(newRoot.Children, oldRoot)

	if oldRoot.IsLeaf {
		t.splitLeaf(oldRoot)
	} else {
		t.splitInternal(oldRoot)
	}
}

func (t *BPlusTree) GetRank(member string) int {
	item, exist := t.MemberMap[member]
	if !exist {
		return -1
	}

	// O(log N) rank query using augmented Size fields
	return t.getRankByItem(item)
}

func (t *BPlusTree) getRankByItem(targetItem *Item) int {
	accumulatedItems := 0
	node := t.Root

	// Accumulate sizes of left subtrees
	for !node.IsLeaf {
		i := 0
		for i < len(node.Items) && targetItem.Score >= node.Items[i].Score {
			accumulatedItems += node.Children[i].Size
			i++
		}
		node = node.Children[i]
	}

	// Find position in leaf
	for i, item := range node.Items {
		if item.Member == targetItem.Member {
			return accumulatedItems + i
		}
	}

	return -1
}
