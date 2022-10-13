# update-ami
更新されたAuto Scaling Groupの設定をECSコンテナインスタンスに反映させるため，インスタンスの入れ替え作業を行うCLI

# Usage

## Install

```
go install github.com/noritama73/update-ami/cmd/update-ami@latest
```

## Example

```
update-ami replace-instances --region ap-northeast-1 --profile <user>@<account> --cluster <cluster> --max-attempt 20 --delay 10
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

## asg-name (optional)

対象クラスタに紐づくAutoScaling Groupの名前．指定しないとクラスタと同じ文字列が入ります

# 内部的に実行される手順

1. 既存のコンテナインスタンスのIDを控える
2. ASGのdesired capacityを1増やす
3. 新しいインスタンスが追加されるのを待つ
4. 古いインスタンスを1つドレインする
5. ドレインされたらderegister→terminate
6. インスタンスが増えるのを待つ
7. サービスを強制更新
8. ちょっと待つ
9. 4.に戻る
10. 最後のインスタンスをterminateしたらdesired capacityを1減らす（元に戻す）
