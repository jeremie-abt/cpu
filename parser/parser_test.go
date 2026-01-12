package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode_FirstNode(t *testing.T) {
	tree := &rootMetadataTree{
		ctx: new(nodeContextMap),
	}
	tok := &token{value: []rune("root")}

	tree.Node(tok)

	assert.NotNil(t, tree.c, "c should not initialized and not nil")
	assert.Same(t, tree.c.token, tok, "c should point to root token")
	assert.Nil(t, tree.c.parent, "root should not have a parent")
}

func TestNode_AddChild(t *testing.T) {
	tree := &rootMetadataTree{
		ctx: new(nodeContextMap),
	}
	parent := &token{value: []rune("parent")}
	child := &token{value: []rune("child")}

	tree.Node(parent)
	tree.Node(child)

	assert.Same(t, tree.c.token, child, "c should point to child, because this is the last node pushed")
	if assert.NotNil(t, tree.c.parent, "c should have a parent set") {
		assert.Same(t, tree.c.parent, parent, "c parents is badly set")
	}
}

// TODO: remove this tests ?
func TestNode_ParentKeepsChildren(t *testing.T) {
	tree := &rootMetadataTree{
		ctx: new(nodeContextMap),
	}
	parent := &token{value: []rune("parent")}
	child := &token{value: []rune("child")}

	tree.Node(parent)
	parentNode := tree.c
	tree.Node(child)

	if len(parentNode.children) != 1 {
		t.Fatalf("parent should have 1 child, got %d", len(parentNode.children))
	}
	if parentNode.children[0].token != child {
		t.Error("parent's child should be child token")
	}
}

func TestNode_ThreeLevelChain(t *testing.T) {
	tree := &rootMetadataTree{
		ctx: new(nodeContextMap),
	}

	aNode := &token{value: []rune("A")}
	bNode := &token{value: []rune("B")}
	cNode := &token{value: []rune("C")}
	tree.Node(aNode)
	tree.Node(bNode)
	tree.Node(cNode)

	assert.Same(t, tree.c, cNode, "c should be the last node pushed")
	assert.Same(t, tree.c.parent, bNode, "bNode should be the parent of cNode")
	assert.Same(t, tree.c.parent.parent, aNode, "aNode should be the parent of bNode")
	assert.Nil(t, tree.c.parent.parent.parent, "aNode should not have a parent because this is root node")
}
