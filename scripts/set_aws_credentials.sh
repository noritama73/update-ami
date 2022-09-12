#!/bin/bash

# You need to set ~/.aws/credential as below.

# [account_name]
# aws_access_key_id = AKIASAMPLECREDENTIAL
# aws_secret_access_key = AWsSecreTaccESSkEYSAmplEcreDENtiaLXXXXXX

PROFILE=$(cat .config | grep PROFILE | awk -F '=' '{print $2}')
SERIAL_NUMBER=$(cat .config | grep SERIAL_NUMBER | awk -F '=' '{print $2}')
TOKEN_CODE=$1

if [ ! -n "${TOKEN_CODE}" ]; then
	echo "You need to specify TOKEN_CODE."
	exit 1
fi

RES=($(aws --profile ${PROFILE} sts get-session-token \
--serial-number ${SERIAL_NUMBER} \
--token-code ${TOKEN_CODE} \
--query '[Credentials.AccessKeyId, Credentials.SecretAccessKey, Credentials.SessionToken]' \
--output text))

echo "[default]" > .credentials
echo "aws_access_key_id = ${RES[0]}" >> .credentials
echo "aws_secret_access_key = ${RES[1]}" >> .credentials
echo "aws_session_token = ${RES[2]}" >> .credentials
echo "SUCCESS!"
