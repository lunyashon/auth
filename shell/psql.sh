#!/usr/bin/env bash
IFS=$'\n\t'
set -euo pipefail

. /etc/os-release 2>/dev/null || true
OS_FAMILY="${ID_LIKE:-} ${ID:-}"

case "$OS_FAMILY" in
	*debian*|*ubuntu*)
		PM=apt
		;;
	#*rhel*|*fedora*|*centos*|*rocky*|*almalinux*|*ol*|*amzn*)
	#	PM=dnf
	#	;;
	#*suse*|*sles*|*opensuse*)
	#	PM=zypper
	#	;;
	#*alpine*)
	#	PM=apk
	#	;;
	*)
		PM=unknown
		;;
esac

if [ "${PM}" == "unknown" ]; then
	echo "Error: OS is not supported"
	exit 1
fi

if command -v psql >> /dev/null; then
	echo "package postgresql installed"
else
	if sudo $PM update && sudo $PM install -y postgresql postgresql-client >> /dev/null; then
		if command -v psql >> /dev/null; then
			echo "package postgresql installed"
		else
			echo "ERROR: package postgresql installation failed"
			exit 1
		fi
    else
		echo "ERROR: package postgresql installation failed"
		exit 1
	fi
fi

POSTGRES_PASSWORD=""
DB_NAME="sso"

PARSED_FLAGS=$(getopt -o p:n: --long password:,db-name: -- "$@")
eval set -- "${PARSED_FLAGS}"

while true; do
	case "$1" in
		-p|--password) POSTGRES_PASSWORD="$2"; shift 2;;
		-n|--db-name) DB_NAME="$2"; shift 2;;
		--) shift; break;;
	esac
done

if [ -z "${POSTGRES_PASSWORD}" ]; then
	echo "Error: Password is not set"
	exit 1
fi

read -r -p "Are you sure you want to create a new database ${DB_NAME}? (y/n): " confirm
if [ "$confirm" != "y" ]; then
	echo "Database creation cancelled"
	exit 1
fi

psql "postgresql://postgres:${POSTGRES_PASSWORD}@127.0.0.1:5432/postgres" -v ON_ERROR_STOP=1<<SQL
    ALTER USER postgres WITH PASSWORD '${POSTGRES_PASSWORD}';
    DROP DATABASE IF EXISTS ${DB_NAME};
    CREATE DATABASE ${DB_NAME};
SQL

psql "postgresql://postgres:${POSTGRES_PASSWORD}@127.0.0.1:5432/${DB_NAME}" < ./migrations/sso.sql