Casperlang
==========

This repository contains the golang code for the main interpreter of the Casper programming language.

More information about *casper* can be found [here](https://www.openengineer.dev/casperlang.html).

# Dev guide

Build dependencies:
* make
* go

## Interpreter components

### 1. tokenization

A file is loaded as an absolute path and source content.

The tokenization is a single scan forward of a source file, resulting in a list of tokens. Indentation of every line is also a token.

Each token has a context attached to it, with a start and stop position (start is inclusive, stop is exclusive). Each position indicates line number, col number and global char number of the token's first character. 

Tokenization also groups tokens based on parentheses, braces or brackets. At this stage braces are already assumed to form dicts (with colons and commas), and brackets are already assumed to form lists (with commas).

### 2. parsing (syntax tree construction)

1. import statements are parsed
2. function statements are parsed
3. function headers are built (type patterns)
   1. `::` operator
   2. words and groups
   3. recurse nested patterns
4. function expressions are built
   1. semicolons, assignments
   2. groups
   3. pipes
   4. function calls (essentially a list of expressions)
   5. remaining operators
   6. remaining tokens 
   7. recurse expression parsing for nested expressions

An important realization here is that expressions are essentially deferred values, so expressions respect the `Value` interface.

### 3. import resolution

1. load all the modules by following the import tree. at this point circular imports are detected
2. merge functions between files of the same modules (looping over the modules one by one)
3. load imported functions (package.json prelude acts as an import!), 
  ignore the functions beginning with underscores
4. push methods attached to referenceable constructors to other files in same module, and then up import tree
  ignore the functions beginning with underscores
  (attached means method was defined in same module as constructor was defined)
5. load the core global functions into each file
  
File acts as a collection of functions. At this point non unique constructors throw errors. Each function keeps a reference to its original Scope/File.

### 4. eval

The program entry point is evaluated returning a final IO value.

The `IO.Run()` is called. This starts a cascade of nested IO actions, interspersed with pure function calls.

When functions are called its arguments are substituted immediately into a copy of the function's rhs (unevaluated) value.

Note that it is actually the pattern-matching that initiates non-lazy evaluation (the built-in `fold` function is also non-lazy).

### 5. serialization

In case the program should be run in another environment (eg. javascript in the browser), it could be serialized in a CBOR format. The serialization result is a binary representation of a single File, without any imports.

Optionally the expressions can be simplified as much as possible before serializing.

Separate VM's that run these binary versions of *casper* programs would need to be implemented independently from the compiler.

## Profiling

Start the profiling tool:
```
go tool pprof "profile.dat"
```

Inside the profiling REPL:
```
top 100
```

See top 100 most used functions.
