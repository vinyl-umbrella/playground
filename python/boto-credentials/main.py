import pathlib

import boto3
import botocore.session
from botocore import credentials


def use_credentials_cache(profile: str) -> boto3.Session:
    cli_cache = pathlib.Path.home() / ".aws/cli/cache"

    botocore_session = botocore.session.Session(profile=profile)
    botocore_session.get_component("credential_provider").get_provider("assume-role").cache = (
        credentials.JSONFileCache(str(cli_cache))
    )

    return boto3.Session(botocore_session=botocore_session)


if __name__ == "__main__":
    session = use_credentials_cache("")
    STS = session.client("sts")
    res = STS.get_caller_identity()
    print(res)
