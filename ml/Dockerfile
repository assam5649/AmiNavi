FROM python:3.10-slim

RUN apt-get update

WORKDIR /app

RUN apt-get update && apt-get install -y \
    libgl1 \
    libglib2.0-0 \  
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --timeout 300 -r requirements.txt

COPY app.py .
COPY model /app/model

CMD ["gunicorn", "--bind", "0.0.0.0:8501", "--workers", "1", "app:app", \
     "--log-level", "info", \
     "--access-logfile", "-", \
     "--error-logfile", "-"]
