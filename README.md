# Minitalk

Minitalk is a lightweight, object-oriented programming language inspired by Smalltalk, implemented in Go for performance and simplicity. This README provides an overview of how to run the Minitalk interpreter and explains the key features and syntax of the language.

## Running the Minitalk Interpreter

### Prerequisites

- Install Go (version 1.16 or later recommended).
- Clone the Minitalk repository:

  ```bash
  git clone https://github.com/CostyRo/minitalk.git
  cd minitalk
  ```

### Building and Running

1. **Build the interpreter**:

   ```bash
   go build
   ```

   This generates an executable file (e.g., `minitalk` on Unix-like systems or `minitalk.exe` on Windows).

2. **Run the REPL**:

   ```bash
   ./minitalk
   ```

   This starts the interactive REPL, where you can type Minitalk code and see immediate results.

3. **Run a Minitalk file**:

   ```bash
   ./minitalk path/to/your/script.sm
   ```

   Replace `path/to/your/script.sm` with the path to a file containing Minitalk code.

The source code for the Minitalk interpreter is available in the GitHub repository.

## Minitalk Language Specification

Minitalk is an object-oriented language where everything is an object, and computation is driven by message passing. Below is an overview of its core features and syntax.

### Literals

Minitalk supports various literals for representing data:

- **Numbers**: Integers and real numbers, including different bases (e.g., binary, hexadecimal).

  ```minitalk
  42
  -42
  123.45
  1.2345e2
  2r10010010
  16rA000
  ```
- **Strings**: Enclosed in single quotes. Escaping a single quote requires doubling it.

  ```minitalk
  'quote'
  'escaped''quote'
  ```
- **Symbols**: Preceded by a hashtag (`#`). Quotes are optional if the symbol contains no spaces.

  ```minitalk
  #space
  #'spa ce'
  ```
- **Byte Arrays**: Denoted by `#[]` with space-separated elements.

  ```minitalk
  #[1 2 3 4]
  ```
- **Regular Arrays**: Enclosed in parentheses `()`.

  ```minitalk
  #(#(1 2 3 4) [1 2 3 4] 'four' 4.0 #four)
  ```
- **Code Blocks**: Defined in square brackets `[]`, with an optional argument list (prefixed by `:`) and a body, separated by `|`.

  ```minitalk
  ['empty arguments']
  [:x | 'one argument']
  [:i :x | 'more arguments']
  ```

### Assignment, Printing, and Statement Separator

- **Assignment**: Uses `:=`. Multiple variables can be assigned the same value in a chain.

  ```minitalk
  a := 1
  b := c := 1
  ```
- **Printing**: Uses the `Transcript` object with the `show:` message.

  ```minitalk
  Transcript show: 'this is the printed message'
  ```
- **Statement Separator**: Unlike Smalltalk, Minitalk does not require a dot (`.`) to separate most statements.

  ```minitalk
  Transcript show: 'first message'
  Transcript show: 'second message'
  ```

### Messages

Messages are the core mechanism for computation in Minitalk, categorized into three types:

- **Unary Messages**: Applied directly to an object without arguments.

  ```minitalk
  true not
  #[1 2 3 4] reversed
  1 toString
  ```
- **Binary Messages**: Involve a receiver and one argument, including operators like `==`, `&`, or custom messages like `show:`.

  ```minitalk
  Transcript show: 'receiver object'
  a := #(1 2 3 4) removeAt: 0
  'multiple words in string' splitBy: ' '
  1 == 1  "or equivalently: 1 eq: 1"
  ```
- **Keyword Messages**: Accept multiple arguments, each prefixed by a keyword, for readable and expressive code.

  ```minitalk
  a := #(1 2 3) at: 0 put: 0
  1 to: 5 step: 5
  ```

### Expressions

Expressions are formed by chaining messages, where each message’s result is an object that can receive further messages. Evaluation proceeds from left to right.

```minitalk
1 > 0 not not ifTrue: ['1 2 3 4'] splitBy: ' ' at: 0 toInteger - 1 toBool not
```

This evaluates step-by-step:

 1. `1 > 0` → `true`
 2. `true not` → `false`
 3. `false not` → `true`
 4. `true ifTrue: ['1 2 3 4']` → `'1 2 3 4'`
 5. `'1 2 3 4' splitBy: ' '` → `#('1' '2' '3' '4')`
 6. `#('1' '2' '3' '4') at: 0` → `'1'`
 7. `'1' toInteger` → `1`
 8. `1 - 1` → `0`
 9. `0 toBool` → `false`
10. `false not` → `true`

**Operator Precedence**: Minitalk evaluates left-to-right, which may lead to unexpected results for arithmetic (e.g., `1+2*3` yields `9` instead of `7`). Use parentheses to enforce precedence:

```minitalk
1+(2*3)  "or equivalently: 1 plus: (2 mul: 3)"
```

**Semicolon (**`;`**)**: Sends multiple messages to the same object, requiring a single dot (`.`) to terminate the chain.

```minitalk
Transcript show: 'message1'; show: 'message2'; show: 'message3'.
```

### Code Blocks

Code blocks are powerful, supporting lazy evaluation and closure-like behavior.

- **Lazy Evaluation**: Blocks are not executed until explicitly evaluated with `value`.

  ```minitalk
  a  "raises NameError"
  [a]  "no error, block is not evaluated"
  [a] value  "raises NameError when evaluated"
  ```
- **Block Structure**: A dot (`.`) separates lines in a block’s body.

  ```minitalk
  [1 + 1] value  "returns 2"
  [1. +1] value  "returns +1, as 1 is discarded"
  ```
- **Argument Handling**: The `value` message adapts to the block’s arity. For multiple arguments, Minitalk uses currying.

  ```minitalk
  [:x | x + 1] value: 1  "returns 2"
  [:x :y | x + y] value: 1 value: 2  "returns 3"
  ```

### Built-in Objects

Minitalk provides built-in objects for common tasks:

- `nl`: A string containing a newline (`"\n"`).
- `Transcript`: For console output via `show:`.
- `stdin`: For console input via `nextLine`.
- `FileSystem`: For file system operations.
  - `FileSystem disk` returns a `disk` object.
  - `disk ls: '.'` lists files in the current directory.
  - `disk referenceTo: 'file.txt'` returns a `file` object.
- `file` **Object**: Represents a file with properties (`basename`, `extension`, `size`, etc.) and methods (`contents`, `write`, `append`, etc.).

### Control Structures

Minitalk uses code blocks for control flow:

- **Conditional Branching**:

  ```minitalk
  true ifTrue: [true]  "returns true"
  true ifFalse: [false]  "returns nil"
  true ifTrue: [true] ifFalse: [false]  "returns true"
  ```
- **Iteration**: Uses `do:` with an array, often generated by `to:step:`.

  ```minitalk
  index := 1 to: 10
  index do: [:i | Transcript show: (i toString + nl)]
  ```

  This prints numbers 1 to 10, each on a new line.
- **Error Handling**: Uses `onError` and specific error handlers (`onNameError`, `onZeroDivisionError`, etc.).

  ```minitalk
  'no error' onError: ['onError']  "returns 'no error'"
  1/0 onZeroDivisionError: [0]  "returns 0"
  ```
