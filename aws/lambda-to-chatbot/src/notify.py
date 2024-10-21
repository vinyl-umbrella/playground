import json
import logging
import os

import boto3


def main(event, context):
    l = logging.getLogger(__name__)
    l.info(event)
    msg = json.loads(event["Records"][0]["Sns"]["Message"])

    sns = boto3.client("sns")

    msg = {
        "version": "1.0",
        "source": "custom",
        "content": {
            "title": f':warning: {msg["requestContext"]["functionArn"]} failed',
            "description": f'Error: {msg["responsePayload"]["errorMessage"]}\n```{msg["responsePayload"]["stackTrace"]}```',
        },
    }

    sns.publish(
        TopicArn=os.environ["TOPIC_ARN"],
        Message=json.dumps(msg),
    )
