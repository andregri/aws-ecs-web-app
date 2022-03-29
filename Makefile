build-docker:
	docker build -t andregri/phonebook:$(TAG) backend/

run-docker:
	docker run --rm --network host --env-file .env andregri/phonebook:$(TAG)

push-docker:
	docker push andregri/phonebook:$(TAG)