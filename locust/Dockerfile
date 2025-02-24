FROM python:3.11-slim

WORKDIR /locust

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY locustfile.py .

EXPOSE 8089

CMD ["locust", "--host", "http://backend:8080"]
