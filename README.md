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
a source from of the resource on stdin. The source passed in is entirely
determined by the output of `/check`.

For a `git` resource this will typically be something like:

```json
{
  "uri": "https://github.com/some/repo.git",
  "branch": "develop",
  "sha": "deadbeef"
}
```

The script should fetch the source and place it in the given directory.
