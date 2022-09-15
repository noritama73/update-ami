# update-ami
Amazon Machine Imageの更新作業を簡素化するCLI

# Usage

以下の作業を済ませておく
* AMIのバージョンを更新してlaunch config, Auto Scaliing Groupを更新する

# Required Environment Variables

```
* AWS_REGION
* AWS_ACCESS_KEY
* AWS_SECRET_KEY
(* AWS_ECS_CLUSTER_ID)
```

## 内部的に実行される手順

1. ECSクラスタごとに一つのインスタンスをドレイニングし、一定期間待つ
1. ドレイニングしたインスタンスをterminateする
1. 新しいインスタンスが立ち上がってECSクラスタインスタンスに登録されるのを待つ
1. ECSサービスをforce_deploymentで更新してタスクが偏ったインスタンスに存在しないようにする
1. 1.に戻り、これをクラスタ内のすべてのインスタンスが新しくなるまで続ける
