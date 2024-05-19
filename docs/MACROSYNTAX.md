# Identification

The version will always be at the top of the file like this:
```
1.0|Full| goco macro file starts below this line!
```

- Full
    - Uncompressed, full name of the command
- Compressed
    - Compressed, short name of the commands

# Syntax
Basic key press commands will look like this:
```
1.0|Full| goco macro file starts below this line!

# Keyboard
keyboard|press|a
delay|1000
keyboard|release|a

# Typing
keyboard|type|Hello, World!|100

# Mouse
mouse|set|100|100
mouse|click|left

mouse|down|left
delay|1000
mouse|up|left

mouse|scroll|100

# Numeric Varience
delay|rnd|1000|500
```
