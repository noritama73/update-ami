# update-ami
更新されたAuto Scaling Groupの設定をECSコンテナインスタンスに反映させるため，インスタンスの入れ替え作業を行うCLI

# Usage

## Install

```
go install github.com/noritama73/update-ami/cmd/update-ami@latest
```

## Example

```
update-ami replace-instances --region ap-northeast-1 --profile <user>@<account> --cluster <cluster> --max-attempt 20 --delay 10 --skip-abnormal-instance
```

# Arguments

## cluster (Required)

対象とするECSクラスタ名．無ければAWS_ECS_CLUSTERを読む

## region (Required)

対象のAWSリージョン．無ければAWS_REGIONを読む

## profile (Required)

AWS CLIを使用するユーザのprofile，MFA対応．無ければAWS_PROFILEを読む

## max-attempt (optional)

インスタンスのステータスチェックを行う最大試行回数．デフォルト40回

## delay (optional)

インスタンスのステータスチェックを行う間隔．デフォルト20秒

## skip-abnormal-instance (optional)

処理に不具合のあったインスタンスが出た場合，一旦無視して他のインスタンスへ向けて処理を続行する

# 内部的に実行される手順

1. ECSクラスタごとに一つのインスタンスをドレイニングし、一定期間待つ
1. ドレイニングしたインスタンスをterminateする
1. 新しいインスタンスが立ち上がってECSクラスタインスタンスに登録されるのを待つ
1. ECSサービスをforce_deploymentで更新してタスクが偏ったインスタンスに存在しないようにする
1. 1.に戻り、これをクラスタ内のすべてのインスタンスが新しくなるまで続ける
