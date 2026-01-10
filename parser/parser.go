// This file code will receive tokens, verify that they are semantically correct and return a syntax tree.
package parser

type Parser interface {
	Interpret(t *token)
}

type Tree interface {
	// Node push the token into the current node as a child, moove the cursor on the child and make the link
	// between the parent and the child.
	Node(t *token)

	// swap place the token at the current cursor and place the current cursor token as child.
	// does not move the cursor.
	swap(t *token)

	// Next returns the next token from the lexer.
	Next() *token

	// Root place the cursor on the root of the tree at the global scope.
	Root()

	// SetKey add a key value pair on the current node that can be accessed later.
	SetKey(key string, val any)

	// GetKey retrieve a token attached to a key
	GetKey(key string) *token
}

// Scope maintain a quick table link for variable names, this assume that any symbol can only be
// referenced by going up the Scope tree, from function to global variable for exemple.
type Scope struct {
	parent *Scope
	s      map[string]*token
}

var test = 54

func doSomeThing() {}

func premierF() {
	doSomeThing()
	secondF()

	// TODO: le scope doit être mit dans le tree direct ? je pense car en sois tu dois dire que tu changes de scope
	// d'une manière ou d'une autre
	// Enfaite la question a se poser est comment le changement de scope est matérialisé au sein de l'arbre ?
	// à Savoir quand j'ai un '
	// Edit: les bloc n'existent pas en golang, faire des recherche en ASM intel et pratiquer un peu:
	// -> https://shikaan.github.io/assembly/x86/guide/2024/09/08/x86-64-introduction-hello.html
	// -> https://genxcyber.com/x86-and-x64-assembly-from-scratch/?utm_source
	// -> https://gist.github.com/lancejpollard/1db84c233bcd849b237df76b3a6c4d9e?utm_source
	test1 := 8
	{
		test1 = 5
		{
			test2 := 5
		}
	}
	test1 = 5
	test := 78
	chocolat := test + 5
}

func secondF() {
	test := 5
}

// Un truc bete ou tu ajoute quand besoins, quite a sauté des scopes
func AddScope(parent *Scope) *Scope {
	return &Scope{parent: parent, s: make(map[string]*token)}
}

type parserFunc func(t Tree) parserFunc
