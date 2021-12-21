# Flags
## Pre deploy
If you run your binary with argument `-pre-deploy` the program will check for alters and if there is no alters it will exit with code 0 but if there is an alters it will exit with code 1.

## Force alters
If you run your binary with argument `-force-alters` the program will check for DB and RediSearch alters and it will execute them(only in local mode).
You can use this command on localhost if you use make file:
`mç`

You can use this feature during the deployment process check if you need to execute the alters before you deploy it
