for dir in ./cmd/native-plugins/*; do \
    name=$(basename "$dir"); \
    if [ -d "$dir" ]; then \
        echo "Building $name..."; \
        go build -o ./bin/"$name" "$dir"/. ; \
    fi \
done