# Harbor go client

* [Harbor](http://github.com/vmware/harbor)

```
harborURL := "http://<your_harbor_url>"
auth := harborcli.LoginForm{
        Username: "<your_username>",
        Password:  "<your_password>",
}
client, _ := harborcli.NewHarborClient(harborURL, auth)
if err := client.Login(); err != nil {
        fmt.Println(err)
}
```

> __注：__ 支持 Harbor v1.8.x
