fs-automation
---

Application to execute a command on any event on file/directory. For example, on creation of file in custom directory 
we need to change their permissions:
```yaml
---
rules:
  - operation: CREATE
    path: /tmp
    cmd: chmod 755 {{fullpath}}
```

List of vars for `cmd` section to config:
  * path - will use Path variable as is;
  * fullpath - will use full path. Path option is as a prefix;
  * operation - will use a performed operation (ANY, CREATE, WRITE, REMOVE, RENAME, CHMOD);
  * action - real action which was done(in case of ANY it will show one of CREATE, WRITE, REMOVE, RENAME, CHMOD)
