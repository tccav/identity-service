FROM alpine

RUN adduser -D worker
USER worker
WORKDIR /home/worker

# Copy the Pre-built binary file from local.
COPY --chown=worker:worker /bin/app .

#Command to run the executable
CMD ["./app"]