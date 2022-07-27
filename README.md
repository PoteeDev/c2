# Command & Control server for exploits
## Examples
config.yml
```
routes:
  ssrf:
    path: "ssrf"
    method: "POST"
```
Run 
```
$ docker run --rm -p 80:8080 -v $(pwd)/config.yml:/usr/app/config.yml explabs/potee-c2
```

Generate example action:
```
$ curl -X POST --cookie "hi=hi" "localhost/ssrf/b041306c930ce9?name=bob"
```
b041306c930ce9 - uniq token for action
Get result of action:
```
$ curl -X POST -H "X-Auth-Token: 123" localhost/admin?token=b041306c930ce9
{"cookies":{"hi":"hi"},"params":{"name":"bob"},"body":""}
```