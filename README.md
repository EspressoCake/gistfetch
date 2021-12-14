# gistfetch
Command-line utility to grab Github gists from your own account.

### How do I use this?
- Add an [API token](https://github.com/settings/tokens) with permissions to read `Gists`
- Fetch this module: `go get -u github.com/EspressoCake/gistfetch` (yes, I know `get` is deprecated)
- Supply the aforementioned token as an argument to the application (`-token YOUR_API_TOKEN`)
- Resulting files are written to the current directory
    - An avoidance measure is taken to prevent possible name collisions; the file's length will prepend the suffix if an existing file of the same name is found to exist

### What does the output look like?
![](https://i.ibb.co/fqh86VW/image.jpg)
