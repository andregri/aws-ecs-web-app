build-docker:
	docker build -t andregri/phonebook:1.0 backend/

run-docker:
	docker run --rm --network host andregri/phonebook:1.0

push-docker:
	docker push andregri/phonebook:1.0