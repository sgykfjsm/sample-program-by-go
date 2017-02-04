# SQLX Sample Program

## Getting Start

Create your project directory under `$GOPATH`. In my case, `mkdir -pv ${GOPATH}/github.com/sgykfjsm/sample-program-by-go/sqlx`
Move to your project directory.
Install [glide](https://github.com/Masterminds/glide#install).
Initialize glide.yaml.

```
$ glide init
[INFO]  Generating a YAML configuration file and guessing the dependencies
[INFO]  Attempting to import from other package managers (use --skip-import to skip)
[INFO]  Scanning code to look for dependencies
[INFO]  Writing configuration file (glide.yaml)
[INFO]  Would you like Glide to help you find ways to improve your glide.yaml configuration?
[INFO]  If you want to revisit this step you can use the config-wizard command at any time.
[INFO]  Yes (Y) or No (N)?
N
[INFO]  You can now edit the glide.yaml file. Consider:
[INFO]  --> Using versions and ranges. See https://glide.sh/docs/versions/
[INFO]  --> Adding additional metadata. See https://glide.sh/docs/glide.yaml/
[INFO]  --> Running the config-wizard command to improve the versions in your configuration
```

Install dependencies.
```
$ glide get github.com/jmoiron/sqlx
$ glide get github.com/mattn/go-sqlite3
$ glide install
```



## Reference

- http://jmoiron.github.io/sqlx/
