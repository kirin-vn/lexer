#  THE KIRIN_VN LEXER  #

KiriScript is the format for Kirin-VN game scripts.
The `lexer` program is used to turn a KiriScript file into a series of tokens which can then be processed by other programs.

`lexer` can be built as a plugin and used inside other Go programs, or it can be run from the command-line to view a plain-text representation of the lexed output.

##  The KiriScript Syntax  ##

>  This document is a work in progress! Pls don't get mad 🙏

The KiriScript syntax is inspired by a number of related markup and scripting languages, including [Ren'Py](https://www.renpy.org/), [Fountain](https://fountain.io/) and [ink](http://www.inklestudios.com/ink/).
It is designed to integrate with the Kirin-VN engine.
The syntax is as follows.

###  Basic syntax elements:

####  Pages

Each script consists of one or more __*pages*__, each of which consists of a single file.
Files must have the extension `.kiri`.
Files are referred to by their filename (sans-extension).

####  Whitespace

Line breaks are considered significant by the KiriScript syntax.
The following character sequences are recognized as single line breaks:

 -  A U+000D CARRIAGE RETURN character, followed by a U+000A LINE FEED character
 -  A single U+000D CARRIAGE RETURN character
 -  A single U+000A LINE FEED character

Whitespace (as defined by Unicode) at the beginning and ending of a token is considered insignificant by the KiriScript lexer, and is removed prior to output.
All other whitespace is collapsed; ie, replaced with a single U+0020 SPACE character.

####  Special characters

The following characters have special meaning in KiriScript and may not appear in identifiers:

    !#$%&()*+,-./:;<=>?@[\]^_`{|}~

There are no rules against including whitespace in identifiers; however, recall that all whitespace will be collapsed.

####  Normalization

All text is normalized to Unicode NFC form by the KiriScript lexer.

####  Comments

KiriScript supports two kinds of comments.
The first is a __*line comment*__, which uses the `%%` like so:

    %%  This is a comment
    here is some text %% with a comment

Line comments are completely removed during lexing and replaced with the empty string.
Line comments extend to the end of the source line, and not to the end of the verse.

The second type of comment is a __*boneyard*__ (multi-line) comment, which uses `%(` and `)%` as delimiters:

    here is some %( contains
        a comment

    and any amount of other stuff
    )% text

Boneyard comments are removed during lexing and replaced with:

 -  If the comment contains at least one newline: a single LINE FEED character.
 -  Otherwise: a single space.

####  Blocks

Each KiriScript page is broken up into a number of __*blocks*__, each of which is comprised of one or more verses (see below).
The rules for dividing the page into blocks is as follows:

 -  A block is open at the start of the file.
 -  The old block is closed and a new block is opened for every blank line that is not a part of a note or comment.
 -  At the end of the file, the last block is closed.
 -  Empty blocks are discarded.

####  Verses

A __*verse*__ is a span of text which is treated as though it were a single line.
Usually these are *written* on a single line, but they may span multiple if the succeeding lines begin with `;`.
For example:

    This is a verse
    ;   which spans multiple lines in the source.

Upon processing, the intervening whitespace and `;` characters are replaced with a single space.

####  Spans

A __*span*__ is the main body text for certain kinds of verses.
Spans can consist of any of the following three components, in order:

 -  A list of boolean checks
 -  The span content
 -  A command

Each __*boolean check*__ must consist of a single boolean value (see below).
Any number of boolean checks can be included in the span.
The span content is only rendered if all of the checks evaluate to `[true]`.

    [check 1][check 2] A span with checks.

__*Span content*__ is a string value consisting of the remaining content of the span, discounting the terminal command (if present; see below).
It may optionally be preceded with an `_`, which will be ignored.

    Here is some span content.
    _[We need an underscore because this is not a check.]

Span content may contain formatting; see the section on formatting below.

Commands are executed after all span content has been processed and displayed.
They are only executed if all boolean checks passed.

    [if this check passes] =>> GOTO THIS MOMENT

For a list of commands, see the section on commands below.

Spans can span multiple lines in the source, but they are restricted to a single verse.

###  Data types:

####  Identifiers

An __*identifier*__ is the name for a variable.
Identifiers may not contain special characters but may contain spaces.
Identifiers are case-sensitive.

    This is a valid identifier

Identifiers which are case-insensitive matches for `true`, `yes`, `on`, and `y` are non-overwritable values which correspond to the boolean `[true]` value, and identifiers which are case-insensitive matches for `false`, `no`, `off`, and `n` are non-overwritable values which correspond to the boolean `[false]` value.

Identifiers do not need to be declared, and have the initial value of the empty string.

####  Attributes

An __*attribute*__ is an identifier which describes the current setting or character.
Attributes may be either __*present*__ or __*absent*__.
Present attributes evaluate to `[true]`, while absent attributes evaluate to `[false]`.

Character attributes can only be evaluated inside of character blocks.

####  Numbers

An number is a sequence of digits (0–9) preceded by a mandatory `+` or `-` sign.
Whitespace may appear anywhere inside of an number; it will be ignored.

    + 3279
    -10 020

All KiriScript numbers are signed integers.
Floats are not supported.

####  Booleans

A boolean is a true/false value.
It is accessed via an identifier enclosed in square brackets.
Defined identifiers which are nonzero and not the empty string evaluate to `[true]`; all other identifiers evaluate to `[false]`.

    [True]
    [ NO ]

A `!` character may be placed before the identifier to negate the boolean's value.

    [!False]
    [! YES ]

####  Strings

Any sequence of characters which does not fall into another category is a string.
Strings may optionally be preceded by an `_` character, which is ignored.

    This is a string.
    _[Yes]
    __Only one underscore will be printed.

####  Lists

A list is a sequence of numbers, lists, or identifiers, placed inside of curly braces and separated by `|` characters.

    {some variable | another variable | +1000 | true }
    { this list | { contains another list | inside } }

Strings cannot be included in lists verbatim; they must first be assigned to a variable.

###  Commands

A __*command*__ is a span of text which terminates a span to add an additional effect.
It is okay for a span to consist of only a command.

####  Tags

__*Tags*__ are like line comments except that they are processed and output by the lexer.
Tags are delimited by `#` characters.

    This is some text  # TODO: Actually write something here.

####  Directions

__*Directions*__ are used to control the script flow.
They must begin with of the following character sequences:

| Sequence | Name | Description |
| :------: | :--: | ----------- |
|   `=>`   | CALL | Plays the specified moment, then returns to this point in the script. |
|   `=>>`  | GOTO | Plays the specified moment and continues from there. |
|   `=<`   | DONE | Ends the specified moment. |
|   `=<<`  | EXIT | Ends the script. |
|   `=<>`  | WAIT | Waits for an engine response. |

`=>` and `=>>` must be followed by an identifier specifying the moment to travel to.
For example, the following code travels to the moment titled `MY MOMENT`, then returns to the given point in the script:

    => MY MOMENT

Finally, a direction may end with an argument list, which must contain values or identifiers to pass in/out of the script.
For example, the following code passes `variable` to the engine, and waits for a response.

    =<> { variable }

###  Verse types:

####  Moments

A __*moment*__ is a type of verse which must begin with a period.
Moments identify a particular location in the script, and can be used to break a page up into sections.
The contents of the verse (sans–initial period) provide the moment's identifier, which must be unique within a page.

    .A SIMPLE MOMENT

Moment identifiers may be used with directions for navigation.
Elsewhere in the script, the moment identifier returns the number of times the moment has been viewed.
Conveniently, this means `[A SIMPLE MOMENT]` evaluates to `[true]` for moments which have been viewed at least once, and `[false]` for moments which have not.

Moments can take arguments and be used like functions.
These can be specified with a parenthetical list of identifiers, separated by commas.

    .A MOMENT WITH ARGUMENTS (arg1, arg2)

Moments describe all of the blocks which follow them, until another moment is declared.

####  Setting

A __*setting*__ is a verse identifying a setting.
It begins with a `>`, followed by the setting's identifier, and optionally followed by a parenthetical list of attributes, separated by commas.

    > Setting (morning, raining)

The identifier of a setting can be used to access the setting's name.
This defaults to the string representation of the identifier, but can be changed in the engine.

    > Basement

        We are currently in the `Basement`. %% defaults to Basement, but can be changed.

If a parenthetical is given, it first removes all attributes from the specified setting before adding those provided.
Otherwise, the attributes are inherited from the previous time a setting with that identifier was used.

####  Character

A __*character*__ is a verse identifying a character.
It begins with an `@`, followed by the character's identifier, and optionally followed by a parenthetical list of attributes, separated by commas.

    @CHARACTER (happy, blushing)

The identifier of a character can be used to access the character's name.
This defaults to the string representation of the identifier, but can be changed in the engine.

                @GIRLFRIEND
        What is it?

                @PLAYER
        I really like you, `GIRLFRIEND`! %% defaults to GIRLFRIEND, but can be changed.

If a parenthetical is given, it first removes all attributes from the specified character before adding those provided.
Otherwise, the attributes are inherited from the previous time a character with that identifier was used.

####  Parenthetical

A __*parenthetical*__ verse can be used to add or remove attributes from the current character or setting.
It must begin with a `(` and end with a `)`.
Inside these parentheses must be a list of values, optionally separated by commas, each of which must be either:

 -  `+` followed by the name of the attribute to add
 -  `-` followed by the name of the attribute to subtract
 -  `?` followed by the name of an attribute to subtract if present, or add if not (this toggles the attribute)
 -  `:0` to remove all attributes currently specified
 -  `:^` to reset the attributes to those declared at the beginning of the block.

These values are evaluated from left-to-right, meaning that later values can override previous ones.
Since `:0` and `:^` will remove/reset any attributes previously specified in the verse, these should always come first.
For example, in the following verse:

    (+happy :0 +sad)

…the attribute `happy` is removed by `:0` and only the attribute `sad` is applied.

If `+attr` is specified but `attr` is already present, it is ignored.
Similarly, if `-attr` is specified but `attr` is not present, it is ignored.

####  Choice

A __*choice*__ is a verse which begins with a `*`, `+`, or `-`, and is used to signify a user choice.
Choices which begin with `*` are __*once-only choices*__, and can only be selected once.
Choices which begin with `+` are __*sticky choices*__, and can be selected any number of times.
Choices which begin with `-` are __*fallback choices*__, and can only be selected when no other choices are available.
These characters should be followed by a span labelling the choice.

        * This is a choice.
        + This is a sticky choice.
        - This is a fallback choice.

The initial `*`, `+`, or `-` character may be repeated; this signifies a sub-choice.

          ** This is a choice inside of another choice.

Like all spans, choice spans may begin with a series of boolean values, inside of square brackets.
All of these values must evaluate to `[true]` for the choice to be selectable.

        * [test 1][test 2] Both tests must pass to pick this option.

The spans of choices are evaluated when the choice is displayed.
This includes any formatting or commands inside the span.

        * This WAIT command will be executed immediately =<>

The verses in-between choice verses of the same level are only executed if the choice is selected.
You can use this to display choice-specific text, cause redirection, or perform other advanced functions.

        * This choice displays text.
          Here is some text.
        * This choice executes a GOTO command.
          =>> GOTO THIS MOMENT
        * This choice sets a variable.
          ~ variable = oh yeah
        * This choice changes attributes.
               (-old attribute +new attribute)

####  Operation

An __*operation*__ is a verse which is used to manipulate the value of an identifier.
It begins with a `~` character and is followed by an expression.

    ~ var = 5
    ~ var++
    ~ var =<> {some | data}

>   TODO: More on expressions

####  Continuation verse

A __*continuation verse*__ is a verse that is intended to continue uninterrupted from the preceding verse.
It begins with a `<` character and is followed by a span.

    < This is a continuation verse.

If the first verse which is output by a block is a continuation verse, then it (and any other remaining verses in the block) are treated as though they were a part of the preceding block.

                @CHARACTER
        Here is some dialogue.

        <  %%  It's okay for the continuation verse to be empty.
        This is still dialogue even though it appears in the next block.
        (+and this parenthetical +sets character attributes +not setting ones)

This is particularly handy when combined with moments and cycle blocks; for example, the following code can be used to cycle through character dialogue:

                @CHARACTER
        I need to tell you something.
        => I REALLY LIKE YOU

    .I REALLY LIKE YOU

        :{
            < I really
            < really
            < really
            < really
            < really
            < really like you.
        }

Because only one of them will be processed at a time, each verse in the cycle block needs to be a continuation verse.

Continuation verses are distinguished from continuation *lines* (which begin with a `;`) in that continuation verses are verses in their own right.

####  Plain verse

A __*plain verse*__ is an unadorned verse which does not fit into any of the categories above.
It consists solely of a span.

        This is a line of plain verse.

As with all spans, the span content of plain verses may optionally be preceded by an `_`, which will be ignored.
This is useful in instances where the verse would otherwise be interpreted as a different type.
You can use two `_` characters if you need one to be rendered.

        _...I didn't even know what to think.
        __emily, that was her username, with a single initial underscore.
        _(What kind of a username was that?)

###  Block types:

####  Operation blocks

A __*operation block*__ groups together a number of operations into a single block.
The first and last verse of this block must consist of three `~` characters, optionally separated by whitespace.

    ~~~
        var = 1
        second var = 2
    ~~~

####  Moment blocks

A __*moment block*__ is a block which contains a single moment verse, optionally followed by a setting verse.
If a setting verse is not provided, the setting is not changed.

    .SCENE ONE: MY FAVOURITE COLOUR
    > A Field Of Roses

####  Setting blocks

A __*setting block*__ sets the current setting.
It consists of a single setting verse.

    > FOREST PATH (sunny, autumn)

####  Description blocks

A __*description block*__ is a block containing plain, parenthetical, continuation, choice, or operation verses.
Parenthetical verses in this block affect the current setting.
This block should be used for background narration or setting description.

        The morning was cool and refreshing.
        I was a little tired.

####  Dialogue blocks

A __*dialogue block*__ consists of a character verse, optionally followed by any number of plain, parenthetical, continuation, choice, or operation verses.
Parenthetical verses in this block affect the current character.
It is used to represent dialogue.

                @ALICE (questioning)
        So you really think that they're coming?
                                   (=0 +worried)
        What if they don't like my dress?
                                            (=^)

####  Cycle blocks

Cycle blocks are used to display one from a list of verses.
The first verse in a cycle block must consist of one of the following character sequences:

|  Characters  |   Name    |  Meaning  |
| :----------: | :-------: | --------- |
|     `:{`     |   LIST    | Each time the block is reached, the next verse is displayed. When all of the verses have been cycled through, the last one is displayed perpetually. |
|     `&{`     |   LOOP    | Each time the block is reached, the next verse is displayed. When all of the verses have been cycled through, the cycle starts again from the beginning. |
|     `^{`     | ONLY-ONCE | Each time the block is reached, the next verse is displayed. When all of the verses have been cycled through, nothing is displayed. |
|     `${`     |  SHUFFLE  | Each time the block is reached, a random verse is displayed. |

The last verse in a cycle block must consist of a solitary `}`.

###  Formatting:

The content of spans can optionally contain special formatting.
The following options are available:

####  Accessing variables

You can access the value of a variable by using `` ` `` characters.

        I had seen this before `MOMENT` times.

####  Text formatting

Custom text formatting can be applied using the following syntax:

        \fmt|text content/

####  Character escaping

Newlines can be represented in span content using the character sequence `\n`.
Spaces can be represented using the character sequence `\ ` (a backslash followed by a space); such spaces will not be collapsed.

In addition, the following characters can be escaped by preceding them with a `\` character.
Characters escaped in this manner cannot be used for formatting or to start commands.

    \|/=#%()`:

All other `\` characters are rendered literally.

When escaping comments, the *second* character of the comment delimeter should be escaped; for example, these comments are correctly escaped:

        This is an example %\% of a correctly escaped comment
        as is %\( this %)

These comments, however, are *not* escaped:

        This text has a comment \%% which will be removed by the lexer
        As does this text \%( unfortunately %)

####  Emoticons

Words which begin with a `:` are processed as __*emoticons*__, which can be used to concisely change character attributes or add additional effects.
Here, the `:)` emoticon might be used to add the `smiling` character attribute:

        I really like that idea! :) However, maybe we should consider…

The meaning of emoticons are left to the engine to process.
For example, the following code uses the `:add_apple:` emoticon to add an apple to the user's inventory.

        I picked up the apple and put it in my bag. :add_apple:

Emoticons are processed when the text is rendered, so any emoticons placed in image options will be activated whenever the option text is displayed.
The following code is broken because it will add the apple to the user's inventory regardless of whether the option is selected or not:

        * Pick up the apple. :add_apple:
        * Leave the apple behind.

Instead, the emoticon should be placed on a separate verse:

        * Pick up the apple.
          :add_apple:
        * Leave the apple behind.

Emoticons cannot be passed arguments.
If you need to pass arguments, you should use a `=>` command instead.
The following short script shows this in action:

    =>> YOU ARE OFFERED AN APPLE

    ~~~  %%  Basic variable definitions
        add to inventory = INV_ADD
        apple = apple
    ~~~

    .ADD ITEM (type)

        =<> {add to inventory | apple}

    .YOU ARE OFFERED AN APPLE

        A man offers you an apple.

        * Pick up the apple.
          => ADD ITEM {apple}
        * Leave the apple behind

Unlike redirections, emoticons are asynchronous, and the game engine will not pause while the emoticon is being processed.

###  Loading external files:

>  TODO: This section ;P
