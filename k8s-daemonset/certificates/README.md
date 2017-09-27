Certificates in this folder are for example purposes only,
and are not suitable for use in production.

To generate a new private key, and convert PKCS#1 to PKCS#8:

```bash
openssl req -x509 -nodes -newkey rsa:2048 -subj '/C=US/CN=My CA' -keyout tls.key -out tls.crt
openssl pkcs8 -topk8 -nocrypt -in tls.key -out tls.pk8
```
