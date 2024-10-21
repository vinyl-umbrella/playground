import json
import logging
import os

import boto3


def main(event, context):
    l = logging.getLogger(__name__)
    l.info(event)

    sns = boto3.client("sns")

    msg = {
        "version": "1.0",
        "source": "custom",
        "content": {
            "title": f':warning: {event["requestContext"]["functionArn"]} failed',
            "description": f'Error: {event["responsePayload"]["errorMessage"]}\n```{event["responsePayload"]["stackTrace"]}```',
        },
    }

    sns.publish(
        TopicArn=os.environ["TOPIC_ARN"],
        Message=json.dumps(msg),
    )
