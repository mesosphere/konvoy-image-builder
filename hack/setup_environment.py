import configparser
import os
import sys

AWS_PROFILE_ENV = "AWS_PROFILE"
CREDENTIALS_FILE = "~/.aws/credentials"
DEFAULT_REGION = "us-west-2"

keys = [
    "aws_access_key_id",
    "aws_secret_access_key",
    "aws_session_token",
    "aws_security_token",
]


def main():
    profile = os.environ.get(AWS_PROFILE_ENV)
    if not profile:
        print("{} is not set".format(AWS_PROFILE_ENV))
        sys.exit(1)

    config = configparser.ConfigParser()
    config.read([os.path.expanduser(CREDENTIALS_FILE)])

    if profile not in config.sections():
        print("profile {} does not exist".format(profile))
        sys.exit(1)

    creds = config[profile]

    for k in keys:
        print("export {}={}".format(k.upper(), creds[k]))

    print("export {}={}".format("AWS_DEFAULT_REGION", DEFAULT_REGION))


if __name__ == "__main__":
    main()
