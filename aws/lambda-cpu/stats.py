from collections import Counter, defaultdict

import boto3

dynamodb = boto3.resource("dynamodb")
# Config
DYNAMODB_TABLE_NAME = "lambda-cpu-results"


def fetch_all_results():
    table = dynamodb.Table(DYNAMODB_TABLE_NAME)
    response = table.scan(
        ProjectionExpression="architecture, memory_mb, num_cpu, cpu_model_name, physical_cores, logical_cores"
    )
    items = response["Items"]

    # ページネーション処理
    while "LastEvaluatedKey" in response:
        response = table.scan(
            ExclusiveStartKey=response["LastEvaluatedKey"],
            ProjectionExpression="architecture, memory_mb, num_cpu, cpu_model_name, physical_cores, logical_cores",
        )
        items.extend(response["Items"])

    print(f"Retrieved {len(items)} records")
    return items


def analyze_results(items: list[dict]):

    # results[architecture][memory_mb] = {"logical": Counter, "physical": Counter, "cpu_model": Counter}
    results: dict[str, dict[str, dict[str, Counter]]] = defaultdict(
        lambda: defaultdict(
            lambda: {"logical": Counter(), "physical": Counter(), "cpu_model": Counter()}
        )
    )

    for item in items:
        arch: str = item["architecture"]
        memory_mb: str = item["memory_mb"]
        logical_cores = int(item["logical_cores"])
        physical_cores = int(item["physical_cores"])
        cpu_model_name = item.get("cpu_model_name", "Unknown CPU")

        # Counter で簡潔にカウント
        results[arch][memory_mb]["logical"][logical_cores] += 1
        results[arch][memory_mb]["physical"][physical_cores] += 1
        results[arch][memory_mb]["cpu_model"][cpu_model_name] += 1

    # Markdown テーブル出力
    print("| メモリサイズ (MB) | CPU アーキテクチャ | 論理コア | 物理コア | CPU モデル |")
    print("|------------------|-------------------|---------|---------|-----------|")

    for arch in sorted(results.keys()):
        for memory_mb in sorted(results[arch].keys(), key=lambda x: int(x) if x.isdigit() else 0):
            logical_stats = results[arch][memory_mb]["logical"]
            physical_stats = results[arch][memory_mb]["physical"]
            cpu_model_stats = results[arch][memory_mb]["cpu_model"]

            # 最も多いコア数を取得
            most_common_logical = logical_stats.most_common(1)[0] if logical_stats else (0, 0)
            most_common_physical = physical_stats.most_common(1)[0] if physical_stats else (0, 0)

            # CPUモデルを改行区切りで結合
            cpu_model_items = cpu_model_stats.most_common()
            cpu_models_str = "<br>".join(
                [f"{model} ({count}回)" for model, count in cpu_model_items]
            )

            logical_cores_str = f"{most_common_logical[0]} ({most_common_logical[1]}回)"
            physical_cores_str = f"{most_common_physical[0]} ({most_common_physical[1]}回)"

            print(
                f"| {memory_mb} | {arch} | {logical_cores_str} | {physical_cores_str} | {cpu_models_str} |"
            )

    print()


def main():
    items = fetch_all_results()
    analyze_results(items)


if __name__ == "__main__":
    main()
