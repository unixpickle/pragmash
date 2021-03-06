# Syntax

I am trying to keep the syntax minimal. Here are all the built-in constructs the language will support. Note that this guide is designed for people who already have some programming experience.

## Commands

Scripts in pragmash run a series of commands. Each command has a string output. Similar to Bash and many other shells, a command is a name followed by space-delimited arguments. If an argument needs to have a space in it, the spaces can be escaped or the entire argument can be surrounded by quotes. Other escapes (like `\n`, `\"`, `\\`, etc.) are allowed as well. Here's some examples:

    cp /path1 /path2
    cp "/my path" /newpath
    ls /this/path\ has\ spaces\ without\ quotes
    ls /this/path\\has\\backslashes\nand\nnewlines

In addition, a command's output can be used as an argument to another command using parentheses:

    read (replace http://aqnichol.com aqnichol google)

It is also possible to nest backticks:

    read (replace http://aqnichol.com (read old_domain.txt) google)

Something to note is that *all* tokens and keywords are evaluated the same way. As a result, code like this is valid (as you will see later on):

    "if" a b "{"
        "echo" hey
    "}"

## Comments

Comments are ignored by the parser and runtime. A line is marked as a comment if its first non-whitespace character is "#". Here are some examples:

    # This is a top level comment
    if ... {
        # This is an indented comment
    }

Pragmash does not support block comments. This makes it easy to be sure whether a line is in a comment or not.

## Line continuations

Sometimes, you might want to run a command with a long set of arguments. You can continue a line onto the next line using a backslash. For example:

    puts hey there this is a long string \
        and now i'm making it even longer.

It is recommended that you indent continuations, but it is not required. Note, however, that any whitespace on the next line is counted when the continuation is inside quotes or backticks. For example, the following code would output "hey&nbsp;&nbsp;&nbsp;&nbsp;there":

    puts "hey\
        there"

Line continuations are evaluated before comments. This means that the following code would print "hey#there":

    puts hey\
    #there

However, comments may not themselves have line continuations. For example, this would print "hey there":

    # this is a comment \
    puts hey there

## Variables

Variables exist in a global scope, just like environment variables in Bash scripts.

The `set` pseudo-command sets a variable. The following example would set the variable "x" to the contents of a URL:

    set x (read http://aqnichol.com)

The `get` pseudo-command gets a variable. For example, this would write the contents of a variable "x" to a file:

    write ./home.html (get x)

As a shorthand for `get`, you can use a `$` followed by the variable name:

    write ./home.html $x

This means that, in order to pass the "$" character to a command, you must escape it or include it in quotes:

    puts \$100.00 is your account balance.
    puts "$100.00 is your account balance."

## If conditionals

An empty string is considered "false" in pragmash. Thus, you can use a basic `if` statement to check if a command outputs a non-empty string like this:

    if (read /might/be/empty) {
        puts The file wasn't empty.
    }

Certain commands might be crafted to output some sort of boolean result in this manner. For instance, an `exists` command might return "" if a file doesn't exist and "true" if it does:

    if (exists /some/path) {
        puts The file exists.
    }

But checking if a command outputs an empty string only goes so far. You can also check if any number of arguments are equal:

    if "Not found." (read http://google.com) {
        puts The page couldn't be found.
    }
    ...
    if "a" (echo a) a {
        puts It works!
    }

You can also use the `else if` and `else` keywords as expected:

    if "a" $x {
        echo It's A.
    } else if "b" $x {
        echo It's B.
    } else {
        echo It's not A or B.
    }

To negate a condition in an if statement, you can insert the `not` token:

    if not "a" $x {
        echo It's not A.
    }

Note: because of the above feature, you must take care when checking if an expression is the string "not". I recommend putting "not" as the last argument in such cases:

    if $a not {
        echo It's 'not'.
    }

The following **is not** a solution, since the tokenizer runs before the semantic processor:

    # This probably won't do what you think...
    if "not" $a {
        echo A is 'not'.
    }

## While loops

A `while` loop repeats a block as long as a condition remains true. Conditions for `while` loops work exactly the same way as conditions for `if` statements.

    set x 0
    while (< $x 10) {
        set x (+ $x 1)
        puts $x
    }
    puts Lift-off!

An iteration of a while loop can be skipped with the `continue` built-in command. The while loop can be completely terminated using the `break` built-in command.

## For loops

There are no array types in pragmash; instead, arrays are represented as strings with newline delimiters. You can loop over the lines in a string like this:

    for x (ls /foo/bar) {
        puts Found file called $x
    }

This could be used for other purposes as well, such as iterating through a small range of numbers:

    for x (range 1 11) {
        puts $x
    }
    puts Lift-off!

The variable parameter of the for loop can be omitted if each element is unneeded:

    for (range 10) {
        ...
    }

You can also add a second variable argument to get the index of each iteration:

    for i x (arr a b c) {
        puts index is $i element is $x
	}

An iteration of a for loop can be skipped with the `continue` built-in command. The for loop can be completely terminated using the `break` built-in command.

## Try blocks

Some commands might trigger errors. You can catch errors using a try-catch block:

    try {
        throw "Error, yo."
    } catch {
    }

You can optionally capture the error message in a variable:

    try {
        throw "Error, yo."
    } catch e {
        puts Got error $e
    }
