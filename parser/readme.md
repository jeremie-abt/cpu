Implement a lexer / parser, the more generic possible for my project to build an ASM parser.

Let's see where the generics can be pushed without loosing any sense / usability in this package, will try to keep it simple and usable the more possible so that It can be used by other developpers to parse any language.

This has absolutely nothing to do with compilation, this package will only be responsible to generate a tree and return errors if it doesn't understand its inputs.

Insipired by https://www.youtube.com/watch?v=HxaD_trXwRE.



Roadmap :

- Routing Table(?) to parse known symbols (Variables, Functions name etc ...), to handle those.
- Test on memory, be sure that it has no memory pics so that it can be used in constraint envs.