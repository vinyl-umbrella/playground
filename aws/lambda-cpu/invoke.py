import json
import logging

import boto3

lambda_client = boto3.client("lambda")
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger()

# Config
FUNCTIONS = ["lambda-cpu-results-x86_64", "lambda-cpu-results-arm64"]
MEMORY_SIZES = [128, 256, 512, 1024, 2048, 3008, 4096, 8192, 10240]
ITERATIONS = 100


def update_memory_size(function_name: str, memory_size: int) -> None:
    logger.info(f"Updating {function_name} to {memory_size}MB")

    lambda_client.update_function_configuration(
        FunctionName=function_name,
        MemorySize=memory_size,
    )
    waiter = lambda_client.get_waiter("function_updated")
    waiter.wait(FunctionName=function_name, WaiterConfig={"Delay": 2, "MaxAttempts": 30})


def invoke_function(function_name: str):
    res = lambda_client.invoke(FunctionName=function_name, InvocationType="RequestResponse")
    payload = json.loads(res["Payload"].read().decode("utf-8"))
    logger.info(f"Result: {payload}")


def main():
    # arch, mem_size ごとに 10回ずつ．コールドスタートのために毎回メモリサイズを変更
    for function_name in FUNCTIONS:
        for i in range(ITERATIONS):
            for memory in MEMORY_SIZES:
                update_memory_size(function_name, memory)
                logger.info(
                    f"Invoking {function_name} with {memory}MB (Iteration {i + 1}/{ITERATIONS})"
                )
                invoke_function(function_name)


if __name__ == "__main__":
    main()
