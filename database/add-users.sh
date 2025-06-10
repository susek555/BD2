#!/bin/sh
set -e

BASE_URL="http://localhost:8080"
PASSWORD="PolskaGurom"
USER_COUNT=9

i=1
while [ "$i" -le "$USER_COUNT" ]; do
  username="herakles${i}"
  email="${username}@gmail.com"

  if [ $((i % 2)) -eq 1 ]; then
    # Odd → person
    selector="P"
    person_name="Herakles"
    person_surname="Wielki"

    cat <<EOF | curl --silent --show-error --location "${BASE_URL}/auth/register" \
      --header 'Content-Type: application/json' --data @-
{
  "username": "${email}",
  "email": "${email}",
  "password": "${PASSWORD}",
  "selector": "${selector}",
  "person_name": "${person_name}",
  "person_surname": "${person_surname}"
}
EOF

  else
    # Even → company
    selector="C"
    nip_index=$(( i / 2 - 1 ))
    nip="123456789${nip_index}"
    company_name="${username} Sp. z o.o."

    cat <<EOF | curl --silent --show-error --location "${BASE_URL}/auth/register" \
      --header 'Content-Type: application/json' --data @-
{
  "username": "${email}",
  "email": "${email}",
  "password": "${PASSWORD}",
  "selector": "${selector}",
  "company_name": "${company_name}",
  "company_nip": "${nip}"
}
EOF

  fi

  echo "→ Registered user ${email}"
  i=$((i + 1))
done