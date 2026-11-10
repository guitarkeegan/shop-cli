Output demos/section-cli.gif
Require echo
Set Shell "zsh"
Set FontSize 32
Set Width 2000
Set Height 1200
Set TypingSpeed 100ms

Type "echo 'This is how we build a Section with blocks!'"
Sleep 500ms
Enter
Sleep 1s

Type "./tmp/shop section create -n 'My Newest Section' -b 'content,text:title:Title,richtext:copy:Copy' -b 'logo,image_picker:image:Image' 'image_picker:image:Background Image'"
Sleep 500ms
Enter
Sleep 2s

Type "cat sections/my-newest-section.liquid"
Sleep 500ms
Enter
Sleep 2s
