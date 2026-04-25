# Akane
**Disclaimer**: It has been developed for educational purposes only. I am not responsible for your actions.

Akane is a small project that allows you to use misconfigured WordPress sites to send pings to another site.


### Usage

For **Linux**/**Mac**
```
git clone dépot
cd akane
go build -o akane main.go
./akane
```
For **Windows**
```
git clone dépot
cd akane
go build -o akane.exe main.go
akane.exe
```

The urls.json file is required for Akane to work. You must therefore ensure that it is present.
Its format must remain unchanged as the programme uses it as is; however, you can expand its contents.

To test it, I recommend setting up your own WordPress site locally and using a webhook address as the target, such as one from webhook.site, for example.

### Methodology

By default, WordPress has a file called xmlrpc.php in its root directory. This file enables various actions, which we can list using a simple Python query, for example.

```python
import requests

url = "http://localhost/xmlrpc.php"


headers = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:149.0) Gecko/20100101 Firefox/149.0",
    "Content-Type": "text/xml"
}

payload = """<?xml version="1.0"?>
<methodCall>
  <methodName>system.listMethods</methodName>
</methodCall>
"""

print(f"[*] Send request : {url}")
response = requests.post(url, data=payload, headers=headers)
print(response.text)
```

The response to this request is supposed to list all the methods available on the site in question. The method we are interested in is ‘pingback.ping’.
Essentially, this method exists to allow WordPress to verify that another site has indeed linked to it. However, it can be exploited to launch an attack.

The concept is simple: you send your XML packet containing the ‘pingback.ping’ method, the URL of your target, and the URL of an existing article.
The server will check whether the post exists; if it does, it will send the request. If not, it will report an error and do nothing.

What this offers an attacker:
    - Anonymity: The victim sees connections coming from legitimate WordPress servers, not from their own IP address.
    - Bandwidth: The attacker sends a small XML packet, and the server generates a full request to the target.
    - Bypass: Many firewalls allow requests from known WordPress sites to pass through, as they are considered "good".

### Licence

MIT (For authorized security research only)