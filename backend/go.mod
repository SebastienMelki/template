module github.com/{{ORG}}/{{PROJECT}}/backend

go 1.24

require (
	buf.build/go/protovalidate v0.14.0
	github.com/SebastienMelki/sebuf v0.0.6
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/testcontainers/testcontainers-go v0.35.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.35.0
	google.golang.org/protobuf v1.36.6
)
