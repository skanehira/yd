# yd
YAML Incremental Digger.

![](https://i.gyazo.com/521400d0740ed12c1606a8ab9b618632.gif)

## Installation
```sh
$ go install github.com/skanehira/yd@latest
```

## Usage
All of first, you can read yaml file as following.

```bash
$ yd file.yaml
$ yd https://sample.com/file.yaml
$ yd < file.yaml
$ yd -f file.yaml
```

Next, you can enter some query like `select(.a == .b)` to filter key or values.
`yd` is using [yq](github.com/mikefarah/yq) so the query can enter as well as `yq`.
Please refer [this document](https://mikefarah.gitbook.io/yq/operators) for learn how to use query of `yq`.

## Keybind

| Key      | Description      |
|----------|------------------|
| `Enter`  | focus view       |
| `Esc`    | focus input      |
| `Ctrl-n` | scroll down view |
| `Ctrl-p` | scroll up view   |

## Author
skanehira

## Thanks
- https://github.com/mikefarah/yq
