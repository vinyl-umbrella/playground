import json
import os
from datetime import datetime


time_out_of_handler = datetime.now()
env_user = os.environ.get("USER")


def lambda_handler(event, context):
    time_in_handler = datetime.now()
    print(f"Time out of handler: {time_out_of_handler}")
    print(f"Time in handler: {time_in_handler}")
    print(f"env value: {env_user}")
    return {
        "statusCode": 200,
        "body": json.dumps(
            {
                "time_out_of_handler": str(time_out_of_handler),
                "time_in_handler": str(time_in_handler),
                "env_user": env_user,
            }
        ),
    }
