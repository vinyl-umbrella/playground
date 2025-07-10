# API Gateway + Step Functions with API Key Authentication

このプロジェクトは、AWS SAMを使用してAPI Gateway + Step Functionsの統合を実装し、API Key認証を有効にしたサンプルです。

## アーキテクチャ

- **API Gateway**: API Key認証付きのREST API
- **Step Functions**: AWS API（STS GetCallerIdentity, IAM ListUsers）を呼び出すシンプルなワークフロー
- **Usage Plan**: APIの使用量制限とスロットリング

## デプロイ手順

1. SAM CLIがインストールされていることを確認
```bash
sam --version
```

2. アプリケーションのビルド
```bash
sam build
```

3. デプロイ
```bash
sam deploy --guided
```

初回デプロイ時は `--guided` フラグを使用して設定を保存します。

## 使用方法

### 1. API Keyの取得

デプロイ後、AWS コンソールまたは CLI を使用してAPI Keyの値を取得します：

```bash
aws apigateway get-api-key --api-key <API_KEY_ID> --include-value
```

### 2. API呼び出し

```bash
curl -X POST \
  https://<API_GATEWAY_ID>.execute-api.<REGION>.amazonaws.com/Prod/start-workflow \
  -H "x-api-key: <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "message": "Hello from API Gateway"
    }
  }'
```

### 3. レスポンス例

```json
{
  "executionArn": "arn:aws:states:us-east-1:123456789012:execution:SimpleStateMachine:...",
  "startDate": "2024-01-01T00:00:00.000Z"
}
```

## Step Functions ワークフロー

ワークフローは以下のステップを実行します：

1. **GetCallerIdentity**: 現在のAWSアカウント情報を取得
2. **ListUsers**: IAMユーザーの一覧を取得（最大5件）
3. **FormatResponse**: レスポンスを整形して返却

## 設定

### Usage Plan

- **日次制限**: 1,000リクエスト/日
- **レート制限**: 100リクエスト/秒
- **バースト制限**: 200リクエスト

### IAM権限

Step Functionsには以下の権限が付与されています：
- `sts:GetCallerIdentity`
- `iam:ListUsers`
- CloudWatch Logs への書き込み権限

## ローカル開発

### SAM Local で API を起動

```bash
sam local start-api
```

### Step Functions をローカルでテスト

```bash
sam local invoke SimpleStateMachine --event events/test-event.json
```

## クリーンアップ

```bash
sam delete
```

## トラブルシューティング

### API Key が見つからない場合

```bash
aws apigateway get-api-keys
```

### Step Functions の実行履歴を確認

```bash
aws stepfunctions list-executions --state-machine-arn <STATE_MACHINE_ARN>
```

### CloudWatch Logs を確認

Step Functions の実行ログは CloudWatch Logs で確認できます。
