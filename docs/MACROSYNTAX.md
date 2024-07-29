# Sulfurite Script Syntax Documentation

This document describes the syntax of a Sulfurite script. The script is used to automate keyboard and mouse actions. This is made for macro interpreter v1

## File Structure
The first line will always be the interpreter version in the number form EX: (Notice how it doesn't use semantic versioning. This is to prevent confusion with the version of sulfurite itself)
- 1
- 2
- 13

## Commands

### `mouseset`

The `mouseset` command moves the mouse cursor to a specified position on the screen. It takes three arguments: the x and y coordinates, and a smoothness factor. If the smoothness factor is 0, the cursor moves instantly. Otherwise, the cursor moves smoothly over a specified time and speed.

Syntax: `mouseset <x> <y> <smoothness> <time> <speed>`

### `mousemove`

The `mousemove` command moves the mouse cursor relative to its current position. It takes the same arguments as `mouseset`.

Syntax: `mousemove <x> <y> <smoothness> <time> <speed>`

### `scroll`

The `scroll` command scrolls the mouse wheel. It takes two arguments: the x and y scroll amounts, and a smoothness factor. If the smoothness factor is 0, the scroll is instant. Otherwise, the scroll is smooth over a specified time and speed.

Syntax: `scroll <x> <y> <smoothness> <time> <speed>`

### `drag`

The `drag` command drags the mouse from its current position to a specified position. It takes two arguments: the x and y coordinates.

Syntax: `drag <x> <y>`

### `mousedown`

The `mousedown` command simulates a mouse button press. It takes one argument: the button to press.

Syntax: `mousedown "<button>"`

### `mouseup`

The `mouseup` command simulates a mouse button release. It takes one argument: the button to release.

Syntax: `mouseup "<button>"`

### `keydown`

The `keydown` command simulates a key press event. The key to be pressed is specified in quotes.

Syntax: `keydown "<key>"`

### `keyup`

The `keyup` command simulates a key release event. The key to be released is specified in quotes.

Syntax: `keyup "<key>"`

### `keytap`

The `keytap` command simulates a key press and release event. The key to be tapped is specified in quotes.

Syntax: `keytap "<key>"`

### `type`

The `type` command simulates typing a string of text. The text to be typed is specified in quotes.

Syntax: `type "<text>"`

### `sleep`

The `sleep` command pauses the execution of the script for a specified amount of time. The time is specified in milliseconds.

Syntax: `sleep <milliseconds>`

### `loop`

The `loop` command repeats a block of commands a specified number of times.

Syntax: 

```
loop <count>
<commands>
end
```

### `forever`

The `forever` command is used to create an infinite loop. The commands inside the `forever` block will be executed indefinitely until the script is manually stopped.

Syntax: 

```
forever
<commands>
```

### `end`

The `end` command is used to mark the end of a `loop` or `forever` block.

Syntax: `end`
