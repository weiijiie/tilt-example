#! /bin/sh

KEY="$(tr -dc A-Za-z0-9 < /dev/urandom | head -c 6)"
VALUE="$(tr -dc A-Za-z0-9 < /dev/urandom | head -c 6)"

echo "setting ${KEY}=${VALUE}:"
curl -s -X POST -H "Content-Type: application/json" --data "{ \"value\": \"${VALUE}\" }" "localhost:8080/${KEY}"

echo ""
echo "getting value"

RESPONSE="$(curl -s "localhost:8080/${KEY}")"
echo "${RESPONSE}"
if [ "${RESPONSE}" = "{\"value\":\"${VALUE}\"}" ]; then
  printf "value successfully set\n"
else
  printf "value not successfully set\n"
  exit 1
fi
