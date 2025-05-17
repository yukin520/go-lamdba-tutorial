

# Lambda × Golang Tutorial 

次のような構成で簡単なAPIを実装したい

![img](./assets/aws.drawio.svg)


## Usage

ローカル環境でデバックができます。

```shell
$ docker-compose up -d 
```

```shell
$ curl -XPSOT "http://localhost:9000/2015-03-31/functions/function/invocations"　-d @tests/getItem.json
```
