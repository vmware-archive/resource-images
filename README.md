# Concourse Resource Interface

A resource type identifies abstractly how to pull down, push up, and check for
new versions of some external resource.

For example, the `git` resource type handles cloning repos, checking for new
commits, and pushing code to branches.

## Anatomy of a Resource

### `/opt/resource/check`: Check for new versions.

A resource type's `/check` script is invoked to detect new versions of the
resource. It is given the configured source and current version as a point of
reference on stdin. Note that the current version will be missing if this is
the first time the resource has been used. The expected behavior with no
version specified is to return the most recent version (*not* every version
since the resource's inception).

For example, here's what the input for a `git` resource may look like:

```json
{
  "source": {
    "uri": "...",
    "branch": "develop",
    "private_key": "..."
  },
  "version": { "ref": "61cebfdb274da579de4287347967b580d02d31e3" }
}
```

The script should then output a list of all versions after `version`, in order:

```json
[
  { "ref": "d74e0124818939e857f503734fdb0e7ea5f3b20c" },
  { "ref": "7154febfa9b398361dcbd56566a161c35e7c5186" }
]
```

The list may be empty, if the given version is already the latest.

### `/opt/resource/in`: Fetch a given resource.

The `/in` script is passed a destination directory as `$1`, and is given on
stdin the configured source and, optionally, a precise version of the resource
to fetch.

The script must fetch the resource and place it in the given directory.

Because the input may not specify a version, the `/in` script must print out
the version that it fetched. This allows the upstream to not have to perform
`/check` before `/in`, which can be slow (for git it implies two clones).

Additionally, the script may emit metadata as a list of key-value pairs. This
data is intended for public consumption and will make it upstream, intended to
be shown on the build's page.

Example input, in this case for the `git` resource:

```json
{
  "source": {
    "uri": "...",
    "branch": "develop",
    "private_key": "..."
  },
  "version": { "ref": "61cebfdb274da579de4287347967b580d02d31e3" }
}
```

Note that the `version` may be `null`.

Example output:

```json
{
  "version": { "ref": "61cebfdb274da579de4287347967b580d02d31e3" },
  "metadata": [
    { "name": "commit", "value": "61cebfdb274da579de4287347967b580d02d31e3" },
    { "name": "author", "value": "Hulk Hogan" }
  ]
}
```

### `/opt/resource/out`: Update a resource.

The `/out` script is called with a path to the directory containing the build's
full set of sources as the first argument, and is given on stdin the configured
params and the resource's source information. The source directory is as it was
at the end of the build.

The script must emit the resulting version of the resource. For example, the
`git` resource emits the sha of the commit that it just pushed.

Additionally, the script may emit metadata as a list of key-value pairs. This
data is intended for public consumption and will make it upstream, intended to
be shown on the build's page.

Example input, in this case for the `git` resource:

```json
{
  "params": {
    "branch": "develop",
  },
  "source": {
    "uri": "git@...",
    "private_key": "...",
  }
}
```

Example output:

```json
{
  "version": { "ref": "61cebfdb274da579de4287347967b580d02d31e3" },
  "metadata": [
    { "name": "commit", "value": "61cebfdb274da579de4287347967b580d02d31e3" },
    { "name": "author", "value": "Mick Foley" }
  ]
}
```
