import logging
from collections import Counter, defaultdict

import boto3

dynamodb = boto3.resource("dynamodb")
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger()

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

    for arch in sorted(results.keys()):
        print(f"cpu arch: {arch}")
        print("-" * 80)

        for memory_mb in sorted(results[arch].keys(), key=lambda x: int(x) if x.isdigit() else 0):
            print(f"\n  mem size: {memory_mb} MB")

            logical_stats = results[arch][memory_mb]["logical"]
            if logical_stats:
                print(f"    logical cores:")
                total = sum(logical_stats.values())
                for cores, count in sorted(logical_stats.items()):
                    percentage = (count / total * 100) if total > 0 else 0
                    print(f"      {cores} cores: {count:3d} times ({percentage:5.1f}%)")

            physical_stats = results[arch][memory_mb]["physical"]
            if physical_stats:
                print(f"    physical cores:")
                total = sum(physical_stats.values())
                for cores, count in sorted(physical_stats.items()):
                    percentage = (count / total * 100) if total > 0 else 0
                    print(f"      {cores} cores: {count:3d} times ({percentage:5.1f}%)")

            cpu_model_stats = results[arch][memory_mb]["cpu_model"]
            if cpu_model_stats:
                print(f"    CPU models:")
                total = sum(cpu_model_stats.values())
                # most_common() で頻度順にソート
                for model, count in cpu_model_stats.most_common():
                    percentage = (count / total * 100) if total > 0 else 0
                    print(f"      {model}: {count:3d} times ({percentage:5.1f}%)")

    print("\n" + "=" * 80)


def main():
    items = fetch_all_results()
    analyze_results(items)


if __name__ == "__main__":
    main()
