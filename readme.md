git-request-list
================

Command-line tool for listing pull- and merge-requests from github and gitlab.

Configuration
-------------

`git-request-list` accepts two flags: `-v` enables verbose logging and `-c <path>` overrides the path to the configuration
file to use. All other behaviour is controlled by the configuration file.

The configuration file is written in YAML and must contain a list of sources from which to fetch pull- or merge-requests.
For example, the configuration file below is for fetching all github pull requests for the identified by a given access
token and merge requests from two repositories in a privately hosted Gitlab.

    ---
    sources:
      - api: github
        token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

      - api: gitlab
        host: https://gitlab.example.com
        token: xxxxxxxxxxxxxxxxxxxx
        repositories:
          - foo/bar/baz
          - foo/barbar

Each source can have the following properties:

 - `api: <string>`: Either `github` or `gitlab`. Required for each source.
 - `token: <string>`: A personal API access token. Required for each source.
 - `host: <string>`: The protocol and hostname of the source. For `github` sources the default is `https://api.github.com`, for `gitlab` sources this is a required parameter.
 - `skip_wip: <bool>`: Ignore WIP merge requests from `gitlab` source. This is ignored for `github` sources.
 - `repositories: <list of strings>`: A list of repository names from which to include requests. If this is not given, all visible projects are searched.
