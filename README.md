# update-ami
Amazon Machine Imageの更新作業を簡素化するCLI

# Usage

以下の作業を済ませておく
* AMIのバージョンを更新してlaunch config, Auto Scaliing Groupを更新する

# Arguments

## cluster (Required)

対象とするECSクラスタ名

## region (Required)

対象のAWSリージョン

## profile (Required)

AWS CLIを使用するユーザのprofile，MFA対応

## max-attempt (optional)

インスタンスのステータスチェックを行う最大試行回数

## delay (optional)

インスタンスのステータスチェックを行う間隔

## skip-abnormal-instance (optional)

処理に不具合のあったインスタンスが出た場合，一旦無視して他のインスタンスへ向けて処理を続行する

# 内部的に実行される手順

1. ECSクラスタごとに一つのインスタンスをドレイニングし、一定期間待つ
1. ドレイニングしたインスタンスをterminateする
1. 新しいインスタンスが立ち上がってECSクラスタインスタンスに登録されるのを待つ
1. ECSサービスをforce_deploymentで更新してタスクが偏ったインスタンスに存在しないようにする
1. 1.に戻り、これをクラスタ内のすべてのインスタンスが新しくなるまで続ける
