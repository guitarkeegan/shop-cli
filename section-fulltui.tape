Output demos/section-fulltui.gif
Require echo
Require nvim
Set Shell "zsh"
Set FontSize 24
Set Width 2000
Set Height 1200
Set Margin 60
Set TypingSpeed 100ms

Type "echo 'This is how we build a section!🚀🚀🚀'"
Sleep 500ms
Enter
Sleep 1s

Type "shop"
Sleep 500ms
Enter
Sleep 2s

# select section
Enter
Sleep 500ms

# enter filename
Sleep 500ms
Type "Banner"
Sleep 1s
Enter
Sleep 1s

# choose action, include settings
Enter
Sleep 1s
Type "y"
Sleep 1s

# choose setting
Type "/"
Sleep 500ms
Type "color"
Sleep 1s
Enter

# fill out color scheme
Sleep 500ms
Type "color scheme"
Enter
Sleep 2s
Backspace
Backspace
Backspace
Backspace
Backspace
Backspace
Backspace
Type "_scheme"
Enter
Type "Color Scheme"
Enter

# no more inputs
Sleep 1s
Type "n"

# include section blocks
Sleep 2s
Type "y"
Sleep 1s
# block name
Type "Content"
Enter
Sleep 500ms
Type "content"
Sleep 500ms
Enter

# content block settings
Sleep 500ms
Enter
Type "heading"
Enter
Type"Heading"
Enter

# more inputs?
Sleep 500ms
Type "y"
Sleep 500ms
Type "j"
Type "j"
Type "j"
Sleep 500ms
Enter

# rich text
Sleep 500ms
Type "copy"
Enter
Type "Copy"
Enter

# more inputs
Sleep 500ms
Type "n"

#add more blocks
Sleep 1s
Type "y"
Sleep 500ms
Type "Logo"
Enter
Type "logo"
Enter

# logo settings
Sleep 500ms
Type "/"
Sleep 500ms
Type "im"
Sleep 500ms
Enter

# logo logo id and label
Sleep 500ms
Type "image"
Enter
Type "Image"
Sleep 500ms
Enter

# more inputs?
Sleep 1s
Type "n"
Sleep 1s
Type "n"
Type "jjjjjjjjjjjjjjj"
Sleep 1s
Type "h"
Enter

# open in nvim
Sleep 500ms
Type "cat sections/banner.liquid"
Sleep 500ms
Enter
Sleep 2s
