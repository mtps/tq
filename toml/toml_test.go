package toml

import "testing"

func TestTree_GetNode(t *testing.T) {
	tree := newTree()

	subName := "node1"
	subtree := newTree()
	subtree.comment = subName

	tree.values[subName] = subtree

	{
		node := tree.getNode([]string{subName})
		if node == nil {
			t.Errorf("%s should not be nil", subName)
		}

		tv, ok := node.(*Tree)
		if !ok {
			t.Errorf("%s should be a *Tree", subName)
		}

		if tv.comment != subName {
			t.Errorf("%s.comment should be %s", subName, subName)
		}
	}

	{
		subName := "none"
		node := tree.getNode([]string{subName})
		if node != nil {
			t.Errorf("%s should be nil", subName)
		}
	}
}
