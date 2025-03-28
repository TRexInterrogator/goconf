# GOCONF

Easy config loader package for .env files / os envs & typesafe struct mapping.

## Usage

Create a .env file inside your project to store your environment variables.

**Example .env:**

```.env
DBConnection=YourSqlConnectionString
DBUser=SomeDbUserName
DBPassword=SomeDbPassword
```

**Implementation:**

```.go
type MyConfig struct {
	DBConnection string
	DBUser       string
	DBPassword   string
}

...

config := &confmodels.MyConfig{}
if err := goconf.Load(config, nil); err != nil {
  log.Fatalf("error while loading app config: %v", err)
}
```

## Custom filenames

Exchange `Load(config, nil)` for `Load(config, ".myenv")` to change the name of your local .env file. The filename is represented a relative path to your main application.

## Loading from OS env

By default (if no local .env file can be found) the `Load()` function will try to load all environment variables from OS. OS envs must be the same as in your struct.

OS envs for previous example:

```bash
DBConnection=YourSqlConnectionString DBUser=SomeDbUserName DBPassword=SomeDbPassword go run main.go
```

## Other

- Struct nesting currently not supported
- .env file must be in KEY=VALUE format
- No support for numeric values (string only)
