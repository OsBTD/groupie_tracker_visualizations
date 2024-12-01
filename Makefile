.PHONY: run

run:
	reflex -r '\.go$$|\.html$$|\.css$$' -s -- sh -c "echo 'Running tasks...'; go run main.go"
