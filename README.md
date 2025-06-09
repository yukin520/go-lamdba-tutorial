

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


## DB talble

`TODO_ITEMS` テーブル

| フィールド名   | 型       | キー種別       | 説明                         |
|----------------|----------|----------------|------------------------------|
| id             | Number   | パーティションキー (PK) | TODOの一意なID              |
| updated_at     | String   | ソートキー (SK)        | 更新日時（ISO8601やUNIXタイム）|
| created_at     | String   | -              | 作成日時                     |
| name           | String   | -              | TODOのタイトル               |
| description    | String   | -              | TODOの詳細説明               |
| record_type    | String   | -              | データ種別（"todo" 固定）   |
| completed    | Bool   | -              | 完了状態  |

---

`record_type-index` グローバルセカンダリインデックス

| フィールド名   | キー種別           | 説明                         |
|----------------|--------------------|------------------------------|
| record_type    | パーティションキー | "todo"でフィルタリング       |
| updated_at     | ソートキー         | 更新日時順で並べ替え可能     |