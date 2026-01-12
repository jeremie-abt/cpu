// This file code will receive tokens, verify that they are semantically correct and return a syntax nodeTree.
package parser

type Parser interface {
	Interpret(t *token)
}

// TODO: Implement those two interface.
type Tree interface {
	// Node push the token into the current node as a child, moove the c on the child and make the link
	// between the parent and the child.
	Node(t *token)

	// swap place the token at the current c and place the current c token as child.
	// does not move the c.
	swap(t *token)

	// Next returns the next token from the lexer.
	Next() *token

	// Root place the c on the root of the nodeTree at the global scope.
	Root()
}

type rootMetadataTree struct {
	ctx NodeContext

	// c means the current leaf
	c *leafTree
}

type leafTree struct {
	token *token

	parent   *leafTree
	children []*leafTree
}

func newLeafTree(t *token) *leafTree {
	return &leafTree{
		token: t,
	}
}

func newTree() *rootMetadataTree {
	return &rootMetadataTree{
		ctx: new(nodeContextMap),
	}
}

func (t *rootMetadataTree) Node(token *token) {
	if t.c == nil {
		t.c = newLeafTree(token)
		return
	}

	leaf := newLeafTree(token)
	leaf.parent = t.c
	t.c.children = append(t.c.children, leaf)
	t.c = leaf
}

func (t *rootMetadataTree) swap(token *token) {
	//TODO implement me
	panic("implement me")
}

func (t *rootMetadataTree) Next() *token {
	//TODO implement me
	panic("implement me")
}

func (t *rootMetadataTree) Root() {
	//TODO implement me
	panic("implement me")
}

// NodeContext define a context for a node where you can store or retrieve data, that will work as closure, going up
// until your key is found.
type NodeContext interface {
	// SetContextKey add a key value pair on the current node that can be accessed later.
	SetContextKey(key string, val any)

	// GetContextKey retrieve a token attached to a key
	GetContextKey(key string) *token
}

type nodeContextMap map[string]*token

func (n nodeContextMap) SetContextKey(key string, val any) {
	//TODO implement me
	panic("implement me")
}

func (n nodeContextMap) GetContextKey(key string) *token {
	//TODO implement me
	panic("implement me")
}

// scope maintain a quick table link for variable names, this assume that any symbol can only be
// referenced by going up the Scope nodeTree, from function to global variable for exemple.
type scope struct {
	parent *scope
	s      map[string]*token
}

type parserFunc func(t Tree) parserFunc
