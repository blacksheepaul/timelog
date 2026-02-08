package model

import (
	"testing"
)

// TestBuildCategoryTreePointers tests that the tree building uses pointers correctly
// so that children are preserved in the returned structure
func TestBuildCategoryTreePointers(t *testing.T) {
	// Create test categories
	rootID := uint(1)
	childID := uint(2)
	grandchildID := uint(3)

	categories := []Category{
		{
			ID:          rootID,
			Name:        "Root",
			Color:       "#FF0000",
			Description: "Root category",
			ParentID:    nil,
			Level:       0,
			Path:        "/",
		},
		{
			ID:          childID,
			Name:        "Child",
			Color:       "#00FF00",
			Description: "Child category",
			ParentID:    &rootID,
			Level:       1,
			Path:        "/Root",
		},
		{
			ID:          grandchildID,
			Name:        "Grandchild",
			Color:       "#0000FF",
			Description: "Grandchild category",
			ParentID:    &childID,
			Level:       2,
			Path:        "/Root/Child",
		},
	}

	// Build the tree
	tree := buildCategoryTree(categories)

	// Verify tree structure
	if len(tree) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(tree))
	}

	rootNode := tree[0]
	if rootNode.Category.ID != rootID {
		t.Errorf("Expected root ID %d, got %d", rootID, rootNode.Category.ID)
	}

	if len(rootNode.Children) != 1 {
		t.Fatalf("Expected root to have 1 child, got %d. This suggests children array was copied by value instead of using pointers.", len(rootNode.Children))
	}

	childNode := rootNode.Children[0]
	if childNode.Category.ID != childID {
		t.Errorf("Expected child ID %d, got %d", childID, childNode.Category.ID)
	}

	if len(childNode.Children) != 1 {
		t.Fatalf("Expected child to have 1 grandchild, got %d. This suggests children array was copied by value instead of using pointers.", len(childNode.Children))
	}

	grandchildNode := childNode.Children[0]
	if grandchildNode.Category.ID != grandchildID {
		t.Errorf("Expected grandchild ID %d, got %d", grandchildID, grandchildNode.Category.ID)
	}

	if len(grandchildNode.Children) != 0 {
		t.Errorf("Expected grandchild to have no children, got %d", len(grandchildNode.Children))
	}
}

// TestBuildCategoryTreeMultipleRoots tests tree building with multiple root nodes
func TestBuildCategoryTreeMultipleRoots(t *testing.T) {
	root1ID := uint(1)
	root2ID := uint(2)
	child1ID := uint(3)
	child2ID := uint(4)

	categories := []Category{
		{
			ID:       root1ID,
			Name:     "Root1",
			ParentID: nil,
			Level:    0,
		},
		{
			ID:       root2ID,
			Name:     "Root2",
			ParentID: nil,
			Level:    0,
		},
		{
			ID:       child1ID,
			Name:     "Child1",
			ParentID: &root1ID,
			Level:    1,
		},
		{
			ID:       child2ID,
			Name:     "Child2",
			ParentID: &root2ID,
			Level:    1,
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 2 {
		t.Fatalf("Expected 2 root nodes, got %d", len(tree))
	}

	// Verify each root has exactly one child
	for _, root := range tree {
		if len(root.Children) != 1 {
			t.Errorf("Expected root '%s' to have 1 child, got %d", root.Category.Name, len(root.Children))
		}
	}
}

// TestBuildCategoryTreeEmptyInput tests with empty category list
func TestBuildCategoryTreeEmptyInput(t *testing.T) {
	categories := []Category{}
	tree := buildCategoryTree(categories)

	if len(tree) != 0 {
		t.Errorf("Expected empty tree, got %d nodes", len(tree))
	}
}

// TestBuildCategoryTreeSingleRoot tests with only one root category
func TestBuildCategoryTreeSingleRoot(t *testing.T) {
	categories := []Category{
		{
			ID:       uint(1),
			Name:     "OnlyRoot",
			ParentID: nil,
			Level:    0,
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(tree))
	}

	if len(tree[0].Children) != 0 {
		t.Errorf("Expected root to have no children, got %d", len(tree[0].Children))
	}
}
