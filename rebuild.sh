dep ensure
dropdb lambda
createdb lambda
sql-migrate up
rm -fr ./models
sqlboiler -b migrations postgres --no-hooks
