# Github Bot for Monorepo Automation

This service uses:
 - Ruby on Rails 5.1
 - Octokit Gem for Github API v3
 - Github event webhooks [https://developer.github.com/webhooks](https://developer.github.com/webhooks)

## Events

List down events that need to be covered

### On Pull Request Opened `event: pull_request`
 - Read `labels` list on the commit or pull request body, then assign labels to the pull request
 - ignore in case the label haven't created 

### On Pull Request Comment `event:issue_comment`
 - Check whether this is comment on a pull request
 - if content is `/test` then retest the pull request.


 ## Github Authentication

 Place github credentials in `~/.netrc` file.