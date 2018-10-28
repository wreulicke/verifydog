Verfifydog
===========================

Git-based diff verify tool.

## Install 

You can get [here](https://github.com/wreulicke/verifydog/releases).

## Usage

### verify diff between commits

```bash
$ cat .verifydog.yml
verifiers:
  main.go: main.go
$ verifydog commit-hash1 commit-hash2
{"main.go", true}
```

### show history between commits

```bash
$ cat .verifydog.yml
verifiers:
  main.go: main.go
$ verifydog history commit-hash1 commit-hash2
main.go -->
commit 119e791a1ad5e220fa7096707e4c2ce4cbeb74b6
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 04:44:45 2018 +0900

    polish

commit ac89dfd5da61544eded8baf06854aafbfff15388
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 02:55:05 2018 +0900

    move action as function

commit 9cc8cf4e12481fd66d95d5dabaadbb01ebc40567
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 02:50:16 2018 +0900

    polish impl

commit 629ea6c5d558ce70314b319313ee383d2722d07c
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 02:32:46 2018 +0900

    polish impl

commit 48724784a040de73c7ac220d8028a38a34f33e1d
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 02:14:28 2018 +0900

    naive impl

commit 1eb76df8b8420dc7fe7403dd49d2914b8da3ce1f
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 01:57:28 2018 +0900

    polish

commit ff97ed11429c991df13592f6f8045e01f0c4324e
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 01:46:35 2018 +0900

    [WIP]

commit b02ed1e5657e87b1b94b98c73f5f92a9691547af
Author: wreulicke <saito.masaya@classmethod.jp>
Date:   Sun Oct 28 01:45:17 2018 +0900

    init

<-- main.go
```