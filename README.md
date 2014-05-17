# Winston Resource Images

A resource type identifies abstractly how to pull down, push up, and check for
new versions of some external resource.

For example, the `git` resource type can handle cloning repos, and checking
for new commits.

## Anatomy of a Resource

### `/tmp/resource/check`: Check for new versions.

A resource type's `/check` script is invoked to detect new versions of the
resource. It is given the current "latest" version and the source's location.

These parameters are passed on stdin as a JSON payload like so:

```json
{
  "cursor": {...},
  "location": {...}
}
```

The script should then output a list of all versions between `cursor` and the
version identified by `location`, in order:

```json
[{...}, {...}, {...}]
```

The list may be empty, if for example the cursor version is up-to-date.

All values are up to the resource to control. For example, `cursor` will be
the same format as one of the versions printed by `check`. The value of
`location` comes up upstream configuration; it is up to the resource to parse
it and perform any validations.

### `/tmp/resource/in`: Fetch a given version of the resource.

The `/in` script is passed a destination directory as `$1`, and is given
a version of the resource on stdin. The version passed in is entirely
determined by the output of `/check`.

For a `git` resource this will typically be something like:

```json
{
  "uri": "https://github.com/some/repo.git",
  "branch": "develop",
  "sha": "deadbeef"
}
```

The script should fetch the resource identified by the version and place it in
the given directory.
