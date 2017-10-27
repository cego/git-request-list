git-request-list
================

Command-line tool for listing pull- and merge-requests from github and gitlab.

Configuration
-------------

`git-request-list` accepts two flags: `-v` enables verbose logging and `-c <path>` overrides the path to the configuration
file to use. All other behaviour is controlled by the configuration file.

The configuration file is written in YAML and should contain a list of sources from which to fetch pull- or merge-requests.
For example, the configuration file below is for fetching github pull requests for the user identified by a given access
token and merge requests from all repositories in the `foo` namespace and the `bar/baz` project in a privately hosted Gitlab.

    ---
    sort_by: updated
    sources:
      - api: github
        token: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

      - api: gitlab
        host: https://gitlab.example.com
        token: xxxxxxxxxxxxxxxxxxxx
        repositories:
          - ^foo/
          - ^bar/baz$

Each source can have the following properties:

 - `api: <string>`: Either `github` or `gitlab`. Required for each source.
 - `token: <string>`: A personal API access token. Required for each source.
 - `host: <string>`: The protocol and hostname of the source. For `github` sources the default is `https://api.github.com`, for `gitlab` sources this is a required parameter.
 - `repositories: <list of strings>`: A list of regular expressions. Only requests from repositories with names matching one of these will be included in the output. If this is not given, all visible projects are searched.

At top-level, you can specify the following properties:
  - `sort_by: <string>`: The property to sort output by. Either `repository`, `name`, `state`, `url`, `created` or `updated`.
