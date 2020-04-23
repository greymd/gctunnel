# Gmail to Google Calendar tunneling CLI 


## Why I made it ?
I know there are several CLI tools providing Gmail or Google Calendar functionalities.

i.e:
* [insanum/gcalcli](https://github.com/insanum/gcalcli)
* [ThomasHabets/cmdg](https://github.com/ThomasHabets/cmdg)
* [yoshinari-nomura/glima](https://github.com/yoshinari-nomura/glima)

Technically, what this tool aiming can be done by combining above tools.

However, combining multiple tools may cause many troubles in the future.
For example, we have to authorize multiple applications and requires multiple Client ID / Secret.
And also, I want to avoid combining multiple tools developed by individuals or multiple language like Ruby, Python, Node.js like above.
Because even one of them loses compatibility or stops maintenance, entire script will stop to work.
Unfortunately, as far as I searched, there is no CLI tool providing entire functionalities of Gmail/Google Calendar both.
