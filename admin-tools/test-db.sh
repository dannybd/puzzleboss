#!/bin/bash

export PERL5LIB='.:/canadia/puzzleboss/lib'

eval $(perl -MPB::Config -e 'PB::Config::export_to_bash();')

mysql -h $PB_DATA_DB_HOST -P $PB_DATA_DB_PORT -u $PB_DATA_DB_USER -p$PB_DATA_DB_PASS puzzleboss$PB_DEV_VERSION
