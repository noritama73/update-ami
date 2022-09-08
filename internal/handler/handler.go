package handler

import "github.com/urfave/cli"

func ReplaceClusterInstnces(c *cli.Context) error {
	// クラスタのコンテナインスタンス一覧を取得

	// インスタンスをドレイン( update-container-instances-state )

	// インスタンスをterminate(termiane-instance)

	// 新しいインスタンスが登録されるのを待つ(ヘルスチェックの猶予は300秒)

	// ecsサービスを--force-new-deployment

	// 全てのインスタンスが更新されるまで繰り返す
	//（ひとまず最初に取得したインスタンスを全てterminateしたら正常終了？）
	return nil
}
