FROM hashicorp/terraform:latest

RUN apk add --no-cache bash python3 py3-pip && pip install --break-system-packages terraform-local

WORKDIR /app
COPY . .

RUN chmod +x wait-for-it.sh
RUN chmod +x apply.sh

ENTRYPOINT [ "bash"]
CMD ["./wait-for-it.sh" , "-h", "localstack", "-p", "4566" , "--strict" , "--timeout=300" , "--" , "./apply.sh"]