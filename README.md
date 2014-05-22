# Winston Resource Images

A resource type identifies abstractly how to pull down, push up, and check for
new source of some external resource.

For example, the `git` resource type can handle cloning repos, and checking
for new commits.

A source is defined to be a concrete manifestation of a logical resource: for
example, a `sha` identifying a commit of a `git` resource.

## Anatomy of a Resource

### `/tmp/resource/check`: Check for new sources.

A resource type's `/check` script is invoked to detect new sources. It is
given a current source as a point of reference on stdin. Note that the given
source may be missing the precise version, if for example this is the first
time the resource has been used.

For example, here's what the input for a `git` resource may look like:

```json
{
  "uri": "...",
  "branch": "develop",
  "ref": "61cebfdb274da579de4287347967b580d02d31e3"
}
```

In this case, `ref` may be omitted if there is no absolute point of reference.
The expected behavior of any resource is to return the most recent source.

The script should then output a list of all sources after `current`, in order:

```json
[
  {
    "uri": "...",
    "branch": "develop",
    "ref": "d74e0124818939e857f503734fdb0e7ea5f3b20c"
  },
  {
    "uri": "...",
    "branch": "develop",
    "ref": "7154febfa9b398361dcbd56566a161c35e7c5186"
  }
]
```

The list may be empty, if for example the given source is already the latest.

### `/tmp/resource/in`: Fetch a given source.

The `/in` script is passed a destination directory as `$1`, and is given
a source from of the resource on stdin. The source passed in may be the output
of `/check`, or a more open-ended source provided by user configuration (i.e.
containing a git branch but not a SHA).

Because the input may be open-ended, the `/in` script must print out the
source that it fetched. This allows the upstream to not have to perform
`/check` before `/in`, which can be slow (for git it implies two clones).

For a `git` resource the input source will typically be something like:

```json
{
  "uri": "https://github.com/some/repo.git",
  "branch": "develop",
  "sha": "deadbeef"
}
```

...but it could also just be:

```json
{
  "uri": "https://github.com/some/repo.git"
}
```

...if that's all the user specified in configuration, and they've directly
triggered a build.

The script must fetch the source and place it in the given directory.
