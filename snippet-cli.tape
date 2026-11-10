Output demos/snippet-cli.gif
Require echo
Set Shell "zsh"
Set FontSize 32
Set Width 2000
Set Height 1200
Set TypingSpeed 100ms

Type "echo 'This is how we build a snippet!'"
Sleep 500ms
Enter
Sleep 1s

Type "./tmp/shop snippet create -n 'My New Snippet' -d 'This is my fab description' 'title:string,image:object,is_good?:boolean'"
Sleep 500ms
Enter
Sleep 2s

Type "cat snippets/my-new-snippet.liquid"
Sleep 500ms
Enter
Sleep 2s
