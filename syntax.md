
# Common types

rdo     // readonly
const   // known at compiletime

\$int    // ownptr<br>
*int    // mutable ptr, not owned<br>
&int    // shorthand for *rdo int<br>
ref int // reference counted variable<br>

[]int   // dynamic array, will resize if necessery <br/>
[N] int // fixed size array where N >= 0<br/>
[int]   // array slice, dynamic size but not resizable<br/>
[string: int] // hashmap(map), key=string, val=in<br/>

N..M    // Range from N to M, iterable

# All keywords

- **var**: Declares a mutable variable
- **rdo**: Declares a readonly variable or marks an expression or type as readonly
- **const**: Same as rdo except forces the value to be known at compiletime
- **type**: Declares a type or used as generic type argument to allow any type
- **struct**: A type that describes the memory layout of specified fields
- **interface**: Defines a blueprint to generate vtables of all types that fit the requirements
- **class**: More advanced struct which supports polymorphism and inheritance
- **enum**: Defines a type that can match as any of the enum entries
- **attr**: Defines attributes for a function or variable
- **alias**: Defined a more complex as a simpler identifier, works similar to C macros but require that their type is known at compiletime
- **async**: Either calls a function on a new thread, similar to goroutines or marks types as atomic[could be changed to usage like a promise in JS]
- **use**: Imports a module or brings something into scope
- **mod**: First line of every file, defines the module that the file and its contents belong to
- **void**: Can be both used as a type or a value and indicates a unit type
- **do**: Used to defined a block or a single line expression, useful for oneline if-statements
- **if**: Defines a conditional branch
- **else**: Defines the false case of a conditional branch
- **for**: Defines a loop, can be used as while true loop, while loop or for loop
- **this**: Denotes a reference to the object being executed on, like a methods target
- **in**: Used for iterating an iterable or checking if a value is in an iterable
- **return**: Stops execution of a function and optionally returns a value 
- more coming
