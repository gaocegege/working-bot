# Dongyue Bot

[![Build Status](https://travis-ci.org/dyweb/working-bot.svg?branch=master)](https://travis-ci.org/dyweb/working-bot)

Inspired by [pouchrobot](https://github.com/pouchcontainer/pouchrobot).

## Usage

````bash
# help
working-bot -h
# start the server (no arguments, might change in the future)
working-bot
````

## Develop

````bash
make dep-install
make install
````

## Features

- Weekly generator automation

## TODO

- [ ] bot's http server listen to github webhook, once it detects an issue with label working is closed,
it will call the generator and send a PR to [dyweb/weekly][1]. After that, it will open an new issue with label `working`.
- [ ] telegram bot [dyweb/weekly#33][3]

## License

Apache 2.0

## About

We had the idea when discussing how to collaborate with other organizations and reduce the difficulty of contributing to weekly content [dyweb/weekly#29][2].  

[1]: https://github.com/dyweb/weekly/
[2]: https://github.com/dyweb/weekly/issues/29
[3]: https://github.com/dyweb/weekly/issues/33
