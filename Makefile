build:
	docker compose up --build

run:
	docker compose up

stop:
	docker compose down

clear:
	docker compose down --volumes --rmi all --remove-orphans
